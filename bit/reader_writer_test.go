package bit

import (
	"fmt"
	"testing"
)

func TestLongReaderLoop(t *testing.T) {
	size := 250
	w := NewWriter(make([]byte, size))
	for i := 0; i < size; i++ {
		v := i % 255
		w.WriteBits(uint64(byte(v)), bitWidthTable[v])
	}
	to := w.Bytes()
	fmt.Printf("len:%v\n", len(to))
	r := NewReader(to)
	for i := 0; i < size; i++ {
		if i%16 == 0 {
			fmt.Printf("\n")
		}
		v, _ := r.ReadBits(bitWidthTable[i%255])
		fmt.Printf("%v:%v,", v, i)
	}

}

func TestLongReader1(t *testing.T) {
	w := NewWriter(make([]byte, 256))
	w.WriteBits(uint64(19), bitWidthTable[19])
	to := w.Bytes()
	r := NewReader(to)
	v, _ := r.ReadBits(bitWidthTable[19])
	fmt.Printf("%v:%v,", v, 19)

}
