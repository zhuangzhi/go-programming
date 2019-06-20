package partition_test

import (
	"fmt"
	"hash/crc32"
	"math"
	"testing"

	"github.com/montanaflynn/stats"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"

	"github.com/zhuangzhi/go-programming/partition"
)

var (
	crc32Table = crc32.MakeTable(crc32.Castagnoli)
	hasher     = crc32.New(crc32Table)
)

type CountTable map[string]int

func (t *CountTable) Add(node string, key int) {
	if t == nil {
		*t = make(map[string]int, 256)
	}
	m := *t
	m[node] = m[node] + 1
}

func (t CountTable) Samples() []float64 {
	samples := make([]float64, 0, len(t))
	for _, v := range t {
		samples = append(samples, float64(v))
	}
	return samples
}

func (t CountTable) Print() {
	for k := range t {
		fmt.Printf("node: %v,\tnumber:%v\n", k, t[k])
	}
}

type IdTable map[string][]int

func (t *IdTable) Add(node string, key int) {
	if t == nil {
		*t = make(map[string][]int, 256)
	}
	m := *t
	m[node] = append(m[node], key)
}

func (t IdTable) Print() {
	for k := range t {
		fmt.Printf("node: %v,\tnumber:%v,\tpartitions: %v\n", k, len(t[k]), t[k])
	}
}

func Contains(in, sub []int) bool {
	for _, v := range sub {
		find := false
		for _, t := range in {
			if t == v {
				find = true
				break
			}
		}
		if !find {
			return false
		}
	}
	return true
}

func TestConsistent(t *testing.T) {
	nodes := [1024]string{}
	nodeCount := 5
	MockNodes(nodes[:], nodeCount)
	count := 10000
	table := IdTable{}
	hash := partition.NewRendezvous(nodes[:nodeCount], 16384)
	for i := 0; i < count; i++ {
		node := hash.Get(fmt.Sprintf("%v", i))
		table.Add(node, i)
	}

	hash2 := partition.NewRendezvous(nodes[:nodeCount-1], 16384)
	table2 := IdTable{}
	for i := 0; i < count; i++ {
		node := hash2.Get(fmt.Sprintf("%v", i))
		table2.Add(node, i)
	}

	for k := range table2 {
		assert.True(t, Contains(table2[k], table[k]))
	}
}

func MockNodes(nodes []string, n int) {
	for i := 0; i < n; i++ {
		nodes[i] = fmt.Sprint("10.13.3.", i)
	}
}

func MockNodesUUID(nodes []string, n int) {
	for i := 0; i < n; i++ {
		nodes[i] = uuid.New().String()
	}
}

func TestGetN(t *testing.T) {
	nodes := [1024]string{}
	nodeCount := 15
	MockNodes(nodes[:], nodeCount)
	hash := partition.NewRendezvous(nodes[:nodeCount], 16384)
	for i := 0; i < 1024; i++ {
		key := fmt.Sprint(i)
		os := hash.GetN(3, key)
		node := hash.Get(key)
		assert.Equal(t, os[0], node)
	}
}

func TestHash(t *testing.T) {
	nodes := [1024]string{}
	nodeCount := 15
	MockNodes(nodes[:], nodeCount)
	table := CountTable{}
	hash := partition.NewRendezvous(nodes[:nodeCount], 16384)
	for i := 0; i < 1000000; i++ {
		node := hash.Get(uuid.New().String())
		table.Add(node, i)
	}
	table.Print()
	samples := table.Samples()
	PrintStatictis(samples)
}

const formatSample = `
Mean:                           %v
Standard deviation:             %v
Standard deviation population:  %v
Sample Standard Deviation:      %v
Coefficient of Variation:       %v%%
Square Variance:                %v
`

func PrintStatictis(samples []float64) {
	mean, _ := stats.Mean(samples)
	sd, _ := stats.StandardDeviation(samples)
	sdp, _ := stats.StandardDeviationPopulation(samples)
	variance, _ := stats.Variance(samples)
	sds, _ := stats.StandardDeviationSample(samples)
	//Coefficient of variation
	coefficient := sd / mean
	fmt.Printf(
		formatSample,
		mean,
		sd,
		sdp,
		sds,
		coefficient*100,
		math.Sqrt(variance),
	)
}
