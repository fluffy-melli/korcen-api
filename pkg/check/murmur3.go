// pkg/check/murmur3.go

package check

import (
	"encoding/binary"
	"math/bits"
)

func Murmur3Hash(key string) uint32 {
	data := []byte(key)
	var seed uint32 = 0
	const (
		c1 uint32 = 0xcc9e2d51
		c2 uint32 = 0x1b873593
	)

	var h1 uint32 = seed
	var roundedEnd int = (len(data) & ^0x3)

	for i := 0; i < roundedEnd; i += 4 {
		k1 := binary.LittleEndian.Uint32(data[i : i+4])
		k1 *= c1
		k1 = bits.RotateLeft32(k1, 15)
		k1 *= c2

		h1 ^= k1
		h1 = bits.RotateLeft32(h1, 13)
		h1 = h1*5 + 0xe6546b64
	}

	var k1 uint32 = 0
	val := len(data) & 0x3
	switch val {
	case 3:
		k1 ^= uint32(data[roundedEnd+2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(data[roundedEnd+1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(data[roundedEnd])
		k1 *= c1
		k1 = bits.RotateLeft32(k1, 15)
		k1 *= c2
		h1 ^= k1
	}

	h1 ^= uint32(len(data))

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}
