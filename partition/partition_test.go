package partition_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Workiva/go-datastructures/queue"

	"github.com/stretchr/testify/assert"

	"github.com/zhuangzhi/go-programming/partition"
)

func TestPartitionTable(t *testing.T) {
	table := partition.NewPartitionTable(1024)
	for i := 0; i < 1024; i++ {
		for j := 0; j < 1024; j++ {
			assert.True(t, table.Add(i, j))
		}
	}
	assert.False(t, table.Add(1, 1))
	table.Reset()
}

func TestDistributeUnownedPartitions(t *testing.T) {
	partitionNumber := 271
	nodes := queue.NewRingBuffer(1024)
	ids := make([]int, partitionNumber)
	for i := 0; i < partitionNumber; i++ {
		ids[i] = i
	}
	partition.Shuffle(ids)
	partitionIDs := queue.NewRingBuffer(1024)
	for i := 0; i < 30; i++ {
		nodes.Put(&partition.Node{
			Replica: partition.Replica{
				Address: fmt.Sprint(i),
				UUID:    fmt.Sprint(i),
			},
			Table: partition.NewPartitionTable(300),
		})
	}
	for i := 0; i < partitionNumber; i++ {
		partitionIDs.Put(ids[i])
	}
	partition.DistributeUnownedPartitions(nodes, partitionIDs, 0)
	for {
		n, err := nodes.Poll(time.Millisecond)
		if err != nil {
			break
		}
		node := n.(*partition.Node)
		ids := node.Table.GetPartitions(0).Ids()
		fmt.Printf("Node: %v, Count: %v, partitions:%v\n", node.Replica.Address, len(ids), ids)
	}
}
