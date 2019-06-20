package partition_test

import (
	"fmt"
	"testing"

	"github.com/zhuangzhi/go-programming/partition"
)

func BenchmarkRendezvousBench(b *testing.B) {
	nodes := [1024]string{}
	nodeCount := 10
	for i := 0; i < nodeCount; i++ {
		// id := uuid.New()
		nodes[i] = fmt.Sprint("10.13.3.", i)
	}
	hash := partition.NewRendezvous(nodes[:nodeCount], 16384)
	for i := 0; i < b.N; i++ {
		// uuid.New().String()
		hash.Get("helkaspodiurafspodfuaposdiuf")
		// if _, ok := table[node]; !ok {
		// 	table[node] = 1
		// } else {
		// 	table[node] = table[node] + 1
		// }
	}
}
