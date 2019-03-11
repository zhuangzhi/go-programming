package bit

import (
	"encoding/binary"
	"io"
)

// LongReader is
type LongReader interface {
	Read() (uint64, error)
}

type BytesLongReader struct {
	data []byte
}

func (r *BytesLongReader) Read() (uint64, error) {
	data := r.data
	size := len(data)
	if size == 0 {
		return 0, io.EOF
	}

	if size >= 8 {
		t := data[:8]
		data = data[8:]
		return binary.BigEndian.Uint64(t), nil
	}

	tmp := uint64(0)
	for i := 0; i < size; i++ {
		shift := uint(8 * (7 - i))
		tmp |= uint64(data[i]) << shift
	}
	data = nil
	return tmp, nil
}

// Reader bit reader
type Reader interface {
	ReadBits(bits byte) (value uint64, err error)
}

func NewReader(data []byte) Reader {
	return &ReaderLong{
		reader: &BytesLongReader{
			data: data,
		},
		cache: 0,
		bits:  0,
	}
}

// ReaderLong Long integer aligned bit reader
type ReaderLong struct {
	reader LongReader
	cache  uint64
	bits   byte
}

// ReadBits read (bits) number bits from array.
func (lr *ReaderLong) ReadBits(bits byte) (value uint64, err error) {
	if bits > lr.bits {
		value = lr.cache
		bits -= lr.bits
		lr.cache, err = lr.reader.Read()
		if err != nil {
			return 0, err
		}
		shift := 64 - bits
		value = value<<bits + uint64(lr.cache>>shift)
		lr.cache &= bitMask[shift]
		lr.bits = shift
		return
	}
	shift := lr.bits - bits
	value = lr.cache >> shift
	//
	// It should be :lr.cache &= 1<<shift - 1, we optimize to a array.
	lr.cache &= bitMask[shift]
	lr.bits = shift
	return
}
