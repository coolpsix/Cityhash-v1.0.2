package cityhash

import (
	"strconv"
	"testing"
)

// static const uint64 k0 = 0xc3a5c85c97cb3127ULL;
// static const int kDataSize = 1 << 20;
// static const int kTestSize = 300;
const dataSize = 1 << 20
const testSize = 300

var data [dataSize]byte

// void setup() {
//   uint64 a = 9;
//   uint64 b = 777;
//   for (int i = 0; i < kDataSize; i++) {
//     a += b;
//     b += a;
//     a = (a ^ (a >> 41)) * k0;
//     b = (b ^ (b >> 41)) * k0 + i;
//     uint8 u = b >> 37;
//     memcpy(data + i, &u, 1);  // uint8 -> char
//   }
// }

func setup() {
	var a uint64 = 9
	var b uint64 = 777
	var i uint64
	for i = 0; i < dataSize; i++ {
		a += b
		b += a
		a = (a ^ (a >> 41)) * k0
		b = (b^(b>>41))*k0 + i
		data[i] = byte(b >> 37)
	}
}

// TestGoogleHashSet from https://github.com/google/cityhash/blob/master/src/city-test.cc
func TestGoogleHashSet(t *testing.T) {
	setup()
	// for i:=0; i < testSize-1; i++
	result := CityHash128(data[:])
	t.Log(strconv.FormatUint(result.First, 16), strconv.FormatUint(result.Second, 16))
	t.Log(dataSize, testSize)
	t.Error("Expected some impl")
}
