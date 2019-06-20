package partition

import (
	"hash"
	"hash/crc32"
	"sort"
)

var (
	crc32Table       = crc32.MakeTable(crc32.Castagnoli)
	EnablePartitions = false
)

type NodeSum struct {
	Name  string
	Sum   uint32
	Score uint64
	Index int
}

type Hash struct {
	//
	nodes          []NodeSum
	partitions     []uint8
	hasher         hash.Hash32
	partitionCount uint32
}

// NewRendezvous new a hasher
func NewRendezvous(nodes []string, partitionCount int) *Hash {
	// Queue.
	hash := &Hash{}
	hash.hasher = crc32.New(crc32Table)
	hash.partitionCount = uint32(partitionCount)
	hash.Add(nodes...)

	return hash
}
func (h *Hash) cachePartitions() {
	h.partitions = make([]uint8, h.partitionCount)
	for i := 0; i < int(h.partitionCount); i++ {
		var maxScore uint64
		var maxNode int
		var score uint64

		for _, node := range h.nodes {
			// G
			score = Hash2int(node.Sum, uint32(i))
			if score > maxScore {
				maxScore = score
				maxNode = node.Index
			}
		}
		h.partitions[i] = uint8(maxNode)
	}
}

func (h *Hash) Add(nodes ...string) {
	for i, node := range nodes {
		h.nodes = append(h.nodes, NodeSum{node, h.hash([]byte(node)), 0, i})
	}
}

func (h *Hash) GetI(key int) string {
	if len(h.partitions) == int(h.partitionCount) {
		p := key % int(h.partitionCount)
		return h.nodes[h.partitions[p]].Name
	}

	keySum := uint32(key)

	var maxScore uint64
	var maxNode string
	var score uint64

	for _, node := range h.nodes {
		// G
		score = Hash2int(keySum, node.Sum)
		if score > maxScore {
			maxScore = score
			maxNode = node.Name
		}
	}

	return maxNode
}

func (h *Hash) GetNI(n, key int) []string {
	if len(h.nodes) == 0 || n == 0 {
		return []string{}
	}
	if n > len(h.nodes) {
		n = len(h.nodes)
	}

	tNodes := [256]NodeSum{}
	nodes := append(tNodes[:0], h.nodes...)

	keySum := uint32(key) % h.partitionCount

	for i := 0; i < len(nodes); i++ {
		nodes[i].Score = Hash2int(keySum, nodes[i].Sum)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Score > nodes[j].Score
	})

	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = string(nodes[i].Name)
	}

	return names

}

func (h *Hash) Partititon(key string) int {
	return int(h.hash([]byte(key)) % h.partitionCount)
}

func (h *Hash) Get(key string) string {
	return h.GetI(h.Partititon(key))
}

// GetN get n nodes
func (h *Hash) GetN(
	n int,
	key string,
) []string {
	return h.GetNI(n, h.Partititon(key))
}

func (h *Hash) hash(node []byte) uint32 {
	h.hasher.Reset()
	h.hasher.Write(node)
	return h.hasher.Sum32()
}
