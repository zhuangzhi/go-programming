package partition

import (
	"container/list"
	"errors"
	"math/rand"
	"time"

	"github.com/Workiva/go-datastructures/queue"
)

const (
	MaxReplicaCount          = 7
	DefaultRetryMultiplier   = 10
	RangeCheckRatio          = 1.1
	MaxRetryCount            = 3
	AggressiveRetryThreshold = 1
	AggressiveIndexThreshold = 3
	MinAvgOwnerDiff          = 3
)

type Replica struct {
	Address string
	UUID    string
}

type PartitionTable []Partitions

func NewPartitionTable(size int) PartitionTable {
	return make([]Partitions, size)
}

var (
	ErrIndexOutoffRange = errors.New("error: index outoff range")
)

type Partitions map[int]bool

func (p Partitions) Clear() {
	if p != nil {
		for k := range p {
			delete(p, k)
		}
	}
}
func (p Partitions) Ids() []int {
	ids := make([]int, 0, len(p))
	for k := range p {
		ids = append(ids, k)
	}
	return ids
}
func (p Partitions) Add(partitionID int) bool {
	if _, ok := p[partitionID]; ok {
		return false
	}
	p[partitionID] = true
	return true
}

func (t PartitionTable) GetPartitions(index int) Partitions {
	if index >= len(t) {
		panic(ErrIndexOutoffRange)
	}

	set := t[index]
	if set == nil {
		set = make(map[int]bool, 1<<14)
		t[index] = set
	}
	return set
}

func (t PartitionTable) Add(index, partitionID int) bool {
	return t.GetPartitions(index).Add(partitionID)
}

func (t PartitionTable) Contains(index, partitionID int) bool {
	if p := t.GetPartitions(index); p != nil {
		_, ok := p[partitionID]
		return ok
	}
	return false
}

func (t PartitionTable) Delete(index, partitionID int) bool {
	if p := t.GetPartitions(index); p != nil {
		_, ok := p[partitionID]
		if ok {
			delete(p, partitionID)
		}
		return ok
	}
	return false
}

func (t PartitionTable) Size(index int) int {
	return len(t.GetPartitions(index))
}

func (t PartitionTable) Reset() {
	for _, ps := range t {
		ps.Clear()
	}
}

type Node struct {
	Replica Replica
	Table   PartitionTable
}

type Generator struct {
}

func (g *Generator) Arrange() [][]Replica {
	return nil
}

func Shuffle(partitions []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(partitions), func(i, j int) { partitions[i], partitions[j] = partitions[j], partitions[i] })
}

func DistributeUnownedPartitions(
	nodes, freePartitions *queue.RingBuffer,
	index int,
) {
	groupSize := nodes.Len()
	maxTries := freePartitions.Len() * groupSize * 10 //DEFAULT_RETRY_MULTIPLIER;
	tries := uint64(0)
	partitionID, err := freePartitions.Get()
	for err == nil && tries < maxTries {
		tries++
		n, _ := nodes.Get()
		node := n.(*Node)
		if node.Table.Add(index, partitionID.(int)) {
			partitionID, err = freePartitions.Poll(time.Millisecond)
		}
		nodes.Offer(node)
	}
}

func assignFreePartitionsToNodeGroup(
	freePartitions *queue.RingBuffer,
	index int,
	node *Node,
) {
	size := freePartitions.Len()
	for i := uint64(0); i < size; i++ {
		partitionID, err := freePartitions.Poll(time.Microsecond)
		if err != nil {
			break
		}
		if !node.Table.Add(index, partitionID.(int)) {
			freePartitions.Offer(partitionID)
		} else {
			break
		}
	}
}

func tryToDistributeUnownedPartitions(
	underLoadedGroups, freePartitions *queue.RingBuffer,
	avgPartitionPerGroup, index, plusOneGroupCount int,
) int {

	// distribute free partitions among under-loaded groups
	maxPartitionPerGroup := avgPartitionPerGroup + 1
	maxTries := freePartitions.Len() * underLoadedGroups.Len()
	tries := uint64(0)
	for tries < maxTries && freePartitions.Len() > 0 && underLoadedGroups.Len() > 0 {
		tries++
		group, err := underLoadedGroups.Poll(time.Millisecond)
		if err != nil {
			break
		}
		node := group.(*Node)

		assignFreePartitionsToNodeGroup(freePartitions, index, node)
		count := node.Table.Size(index)
		if plusOneGroupCount > 0 && count == maxPartitionPerGroup {
			plusOneGroupCount--
			if plusOneGroupCount == 0 {
				// all (avg + 1) partitions owned groups are found
				// if there is any group has avg number of partitions in under-loaded queue
				// remove it.
				old := underLoadedGroups
				underLoadedGroups = queue.NewRingBuffer(old.Cap())
				for {
					g, err := old.Poll(time.Millisecond)
					if err != nil {
						break
					}
					node := g.(*Node)
					if node.Table.Size(index) < avgPartitionPerGroup {
						underLoadedGroups.Put(node)
					}
				}
			}
		} else if (plusOneGroupCount > 0 && count < maxPartitionPerGroup) || (count < avgPartitionPerGroup) {
			underLoadedGroups.Offer(group)
		}
	}
	return plusOneGroupCount
}

func getUnownedPartitions(state [][]*Replica, replicaIndex int) []int {
	freePartitions := make([]int, 0, len(state))
	// if owner of a partition can not be found then add partition to free partitions queue.
	for partitionID := 0; partitionID < len(state); partitionID++ {
		replicas := state[partitionID]
		if replicas[replicaIndex] == nil {
			freePartitions = append(freePartitions, partitionID)
		}
	}
	Shuffle(freePartitions)
	return freePartitions
}

func partitionOwnerAvailable(nodes []*Node, partitionID, replicaIndex int, owner *Replica) bool {
	for _, node := range nodes {
		if node.Replica.Address == owner.Address {
			return node.Table.Contains(replicaIndex, partitionID)
		}
	}
	return false
}

func ContainsInt(set []int, value int) bool {
	for _, v := range set {
		if v == value {
			return true
		}
	}
	return false
}

func initializeGroupPartitions(state [][]*Replica, nodes []*Node, replicaCount int,
	aggressive bool, toBeArrangedPartitions []int) {
	// reset partition before reuse
	for _, node := range nodes {
		node.Table.Reset()
	}
	for partitionID := 0; partitionID < len(state); partitionID++ {
		replicas := state[partitionID]

		for replicaIndex := 0; replicaIndex < MaxReplicaCount; replicaIndex++ {
			if replicaIndex >= replicaCount {
				replicas[replicaIndex] = nil
				continue
			}

			owner := replicas[replicaIndex]
			valid := false
			if owner != nil {
				valid = partitionOwnerAvailable(nodes, partitionID, replicaIndex, owner)
			}
			if !valid {
				replicas[replicaIndex] = nil
			} else if aggressive && replicaIndex < AggressiveIndexThreshold && (toBeArrangedPartitions == nil || ContainsInt(toBeArrangedPartitions, partitionID)) {
				for i := AggressiveIndexThreshold; i < replicaCount; i++ {
					replicas[i] = nil
				}
			}
		}
	}
}

func selectToGroupPartitions(
	index, expectedPartitionCount int,
	toNode, fromNode *Node,
) {
	from := fromNode.Table.GetPartitions(index)
	to := toNode.Table.GetPartitions(index)
	ps := from.Ids()
	for _, p := range ps {
		if len(from) <= expectedPartitionCount ||
			len(to) >= expectedPartitionCount {
			break
		}
		if to.Add(p) {
			delete(from, p)
		}
	}
}

func transferPartitionsBetweenGroups(
	underLoadedGroups *list.List,
	overLoadedGroups *list.List,
	index, avgPartitionPerGroup, plusOneGroupCount int,
) {

	maxPartitionPerGroup := avgPartitionPerGroup + 1
	maxTries := underLoadedGroups.Len() * overLoadedGroups.Len() * DefaultRetryMultiplier
	tries := 0
	expectedPartitionCount := avgPartitionPerGroup
	if plusOneGroupCount > 0 {
		expectedPartitionCount = maxPartitionPerGroup
	}

	for tries < maxTries && underLoadedGroups.Len() > 0 {
		tries++
		v := underLoadedGroups.Front()
		underLoadedGroups.Remove(v)
		toNode := v.Value.(*Node)
		overNode := overLoadedGroups.Front()
		for overNode != nil {
			tmp := overLoadedGroups.Front()
			overLoadedGroups.Remove(tmp)
			fromNode := tmp.Value.(*Node)
			selectToGroupPartitions(index, expectedPartitionCount, toNode, fromNode)
			fromCount := len(fromNode.Table.GetPartitions(index))
			if plusOneGroupCount > 0 && fromCount == maxPartitionPerGroup {
				plusOneGroupCount--
				if plusOneGroupCount == 0 {
					expectedPartitionCount = avgPartitionPerGroup
				}
			}
			if fromCount <= expectedPartitionCount {
				// overLoadedGroupsIterator.remove()
			}
			toCount := len(toNode.Table.GetPartitions(index))
			if plusOneGroupCount > 0 && toCount == maxPartitionPerGroup {
				plusOneGroupCount--
				if plusOneGroupCount == 0 {
					expectedPartitionCount = avgPartitionPerGroup
				}
			}

			if toCount >= expectedPartitionCount {
				break
			}
		}
		if len(toNode.Table.GetPartitions(index)) < avgPartitionPerGroup {
			underLoadedGroups.PushBack(toNode)
		}
	}
}
