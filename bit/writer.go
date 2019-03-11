package bit

import "encoding/binary"

// Writer Bit writer
type Writer interface {
	WriteBits(val uint64, width byte)
	Bytes() []byte
	BitLength() int
}

func NewWriter(to []byte) Writer {
	w := new(WriterLong)
	w.writer = NewBytesLongWriter(to)
	w.cache = 0
	w.bits = 0
	return w
}

// LongWriter is
type LongWriter interface {
	Write(uint64)
}

// BytesLongWriter ..
type BytesLongWriter struct {
	data []byte
	pos  int
	tmp  []byte // 8 byte array
}

// NewBytesLongWriter ..
func NewBytesLongWriter(to []byte) *BytesLongWriter {
	w := new(BytesLongWriter)
	w.data = to
	w.pos = 0
	w.tmp = make([]byte, 8)
	return w
}

func (r *BytesLongWriter) Write(v uint64) {
	if len(r.data)-r.pos < 8 {
		// Copy
		binary.BigEndian.PutUint64(r.tmp, v)
		r.data = append(r.data[:r.pos], r.tmp...)
	} else {
		binary.BigEndian.PutUint64(r.data[r.pos:], v)
	}
	r.pos += 8
}

type WriterLong struct {
	writer *BytesLongWriter
	cache  uint64 // unwritten bits are stored here
	bits   byte   // number of unwritten bits in cache
}

func (w *WriterLong) Bytes() []byte {
	if w.bits > 0 {
		w.writer.Write(w.cache)
		w.cache = 0
		w.bits = 0
	}
	return w.writer.data
}
func (w *WriterLong) BitLength() int {
	return w.writer.pos*64 + int(w.bits)
}

func (w *WriterLong) WriteBits(val uint64, width byte) {
	newbits := w.bits + width
	if newbits < 64 {
		w.cache |= val << (64 - newbits)
		w.bits = newbits
		return
	} else if newbits == 64 {
		w.writer.Write(w.cache | val)
		w.cache = 0
		w.bits = 0
		return
	} else {
		// cache will be filled, and there will be more bits to write
		// "Fill cache" and write it out
		free := 64 - w.bits
		w.writer.Write(w.cache | (val >> (width - free)))
		width -= free
		if width > 0 {
			// Note: n < 8 (in case of n=8, 1<<n would overflow byte)
			w.cache = (val & ((1 << width) - 1)) << (64 - width)
			w.bits = width
		} else {
			w.cache = 0
			w.bits = 0
		}
		return
	}
}
