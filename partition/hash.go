package partition

// Integer hash function
// http://web.archive.org/web/20071223173210/http://www.concentric.net/~Ttwang/tech/inthash.htm

// Mix96 Robert Jenkins' 96 bit Mix Function
func Mix96(a, b, c uint32) uint32 {
	a = a - b
	a = a - c
	a = a ^ (c >> 13)
	b = b - c
	b = b - a
	b = b ^ (a << 8)
	c = c - a
	c = c - b
	c = c ^ (b >> 13)
	a = a - b
	a = a - c
	a = a ^ (c >> 12)
	b = b - c
	b = b - a
	b = b ^ (a << 16)
	c = c - a
	c = c - b
	c = c ^ (b >> 5)
	a = a - b
	a = a - c
	a = a ^ (c >> 3)
	b = b - c
	b = b - a
	b = b ^ (a << 10)
	c = c - a
	c = c - b
	c = c ^ (b >> 15)
	return c
}

// Hash2int hash to uint32 to a uint64
func Hash2int(key0, key1 uint32) uint64 {
	key := uint64(key0) | (uint64(key1) << 32)
	return Hash64(key)
}

// Hash64 64 bit Mix Function
func Hash64(key uint64) uint64 {
	key = (^key) + (key << 21) // key = (key << 21) - key - 1;
	key ^= (key >> 24)
	key += (key << 3) + (key << 8) // key * 265
	key ^= (key >> 14)
	key += (key << 2) + (key << 4) // key * 21
	key ^= (key >> 28)
	key += (key << 31)

	return key
}

func Hash32shift(key uint32) uint32 {
	key = ^key + (key << 15) // key = (key << 15) - key - 1;
	key = key ^ (key >> 12)
	key = key + (key << 2)
	key = key ^ (key >> 4)
	key = key * 2057 // key = (key + (key << 3)) + (key << 11);
	key = key ^ (key >> 16)
	return key
}

func Hash32(a uint32) uint32 {
	a = (a + 0x7ed55d16) + (a << 12)
	a = (a ^ 0xc761c23c) ^ (a >> 19)
	a = (a + 0x165667b1) + (a << 5)
	a = (a + 0xd3a2646c) ^ (a << 9)
	a = (a + 0xfd7046c5) + (a << 3)
	a = (a ^ 0xb55a4f09) ^ (a >> 16)
	return a
}

func MurmurHash3Mixer(key uint64) uint64 {
	key ^= (key >> 33)
	key *= 0xff51afd7ed558ccd
	key ^= (key >> 33)
	key *= 0xc4ceb9fe1a85ec53
	key ^= (key >> 33)
	return key
}
