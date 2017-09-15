package cityhash

import (
	"encoding/binary"
)

type Uint128 struct {
	First  uint64
	Second uint64
}

var k0 uint64 = 0xc3a5c85c97cb3127
var k1 uint64 = 0xb492b66fbe98f273
var k2 uint64 = 0x9ae16a3b2f90404f
var k3 uint64 = 0xc949d7c7509e6557

func rotate(val uint64, shift uint) uint64 {
	if shift == 0 {
		return val
	} else {
		return (val >> shift) | (val << (64 - shift))
	}
}

func rotateByAtLeast1(val uint64, shift uint) uint64 {
	return (val >> shift) | (val << (64 - shift))
}

func weakHashLen32WithSeeds(w uint64, x uint64, y uint64, z uint64, a uint64, b uint64) Uint128 {
	a += w
	b = rotate(b+a+z, 21)
	c := a
	a += x + y
	b += rotate(a, 44)
	return Uint128{a + z, b + c}
}

func shiftMix(val uint64) uint64 {
	return val ^ (val >> 47)
}

func hashLen0to16(s []byte) uint64 {
	if len(s)>8 {
		a:=fetch64(s)
		b:=fetch64(s[len(s)-8:])
		return hashLen16(Uint128{a,rotateByAtLeast1(b+uint64(len(s)), uint(len(s)))})^b;
	}
	if len(s)>=4{
		a:=fetch32(s)
		return hashLen16(Uint128{uint64(len(s))+(uint64(a)<<3),uint64(fetch32(s[len(s)-4:]))})
	}
	if len(s)>0 {
		a:=s[0]
		b:=s[len(s)>>1]
		c:=s[len(s)-1]
		y:=uint64(a)+(uint64(b)<<8)
		z:=uint64(len(s))+(uint64(c)<<8)
		return shiftMix(y*k2^z*k3)*k2;
	}
	return k2;
}

func cityMurmur(s []byte, seed Uint128) Uint128 {
	a := seed.First
	b := seed.Second
	var c uint64
	var d uint64
	if len(s) <= 16 {
		a = shiftMix(a*k1) * k1
		c = b*k1 + hashLen0to16(s)
		var m uint64
		if len(s) >= 8 {
			m = fetch64(s)
		} else {
			m = c
		}
		d = shiftMix(a + m)
	} else {
		c = hashLen16(Uint128{fetch64(s[len(s)-8:]) + k1, a})
		d = hashLen16(Uint128{b + uint64(len(s)), c + fetch64(s[len(s)-16:])})
		a += d
		for len(s) > 16 {
			a ^= shiftMix(fetch64(s)*k1) * k1
			a *= k1
			b ^= a
			c ^= shiftMix(fetch64(s[8:])*k1) * k1
			c *= k1
			d ^= c
			s = s[16:]
		}
	}
	a = hashLen16(Uint128{a, c})
	b = hashLen16(Uint128{d, b})
	return Uint128{a ^ b, hashLen16(Uint128{b, a})}
}

func cityHash128WithSeed(s []byte, seed Uint128) Uint128 {
	if len(s) < 128 {
		return cityMurmur(s, seed)
	}
	l:=len(s)
	offset:=0
	x := seed.First
	y := seed.Second
	z := uint64(len(s)) * k1
	var v, w Uint128
	v.First = rotate(y^k1, 49)*k1 + fetch64(s)
	v.Second = rotate(v.First, 42)*k1 + fetch64(s[8:])
	w.First = rotate(y+z, 35)*k1 + x
	w.Second = rotate(x+fetch64(s[88:]), 53) * k1
	for l > 128 {
		for i := 0; i < 2; i++ {
			x = rotate(x+y+v.First+fetch64(s[offset+16:]), 37) * k1
			y = rotate(y+v.Second+fetch64(s[offset+48:]), 42) * k1
			x ^= w.Second
			y ^= v.First
			z = rotate(z^w.First, 33)
			v = weakHashLen32WithSeeds(fetch64(s[offset:]), fetch64(s[8+offset:]), fetch64(s[16+offset:]), fetch64(s[24+offset:]), v.Second*k1, x+w.First)
			w = weakHashLen32WithSeeds(fetch64(s[32+offset:]), fetch64(s[40+offset:]), fetch64(s[48+offset:]), fetch64(s[56+offset:]), z+w.Second, y)
			z, x = x, z
			offset+=64
		}
		l-=128
	}
	y += rotate(w.First, 37)*k0 + z
	x += rotate(v.First+z, 49) * k0
	for tailDone := 0; tailDone < l; {
		tailDone += 32
		y = rotate(y-x, 42)*k0 + v.Second
		w.First += fetch64(s[offset+l-tailDone+16:])
		x = rotate(x, 49)*k0 + w.First
		w.First += v.First
		v = weakHashLen32WithSeeds(fetch64(s[offset+l-tailDone:]), fetch64(s[offset+l-tailDone+8:]), fetch64(s[offset+l-tailDone+16:]), fetch64(s[offset+l-tailDone+24:]), v.First, v.Second)
	}
	x = hashLen16(Uint128{x, v.First})
	y = hashLen16(Uint128{y, w.First})
	return Uint128{hashLen16(Uint128{x + v.Second, w.Second}) + y,
		hashLen16(Uint128{x + w.Second, y + v.Second})}
}

func hashLen16(x Uint128) uint64 {
	var mul uint64 = 0x9ddfea08eb382d69
	a := (x.First ^ x.Second) * mul
	a ^= a >> 47
	b := (x.Second ^ a) * mul
	b ^= b >> 47
	b *= mul
	return b
}

func fetch64(d []byte) uint64 {
	return binary.LittleEndian.Uint64(d)
}

func fetch32(d []byte) uint32 {
	return binary.LittleEndian.Uint32(d)
}

func CityHash128(s []byte) Uint128 {
	if len(s) >= 16 {
		return cityHash128WithSeed(s[16:], Uint128{fetch64(s) ^ k3, fetch64(s[8:])})
	} else if len(s) >= 8 {
		return cityHash128WithSeed(nil, Uint128{fetch64(s) ^ (uint64(len(s)) * k0), fetch64(s[uint64(len(s))-8:]) ^ k1})
	} else {
		return cityHash128WithSeed(s, Uint128{k0, k1})
	}
}
