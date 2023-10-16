package bitset

import (
	"fmt"

	"sudonters/zootler/internal/bag"
)

func IsEmpty(b Bitset64) bool {
	for i := range b.bs {
		if b.bs[i] != 0 {
			return false
		}
	}
	return true
}

func Buckets(max int64) int64 {
	return max/bs64Size + 1
}

func New(k int64) Bitset64 {
	bs := make([]int64, k)
	return Bitset64{k, bs}
}

func NewFrom(b Bitset64) Bitset64 {
	n := New(b.k)
	copy(n.bs, b.bs)
	return n
}

func SetMany(b *Bitset64, toSet ...int64) {
	for _, i := range toSet {
		b.Set(i)
	}
}

func ClearMany(b *Bitset64, toClear ...int64) {
	for _, i := range toClear {
		b.Clear(i)
	}
}

func TestMany(b Bitset64, toTest ...int64) []bool {
	res := make([]bool, len(toTest))

	for i := range toTest {
		res[i] = b.Test(toTest[i])
	}

	return res
}

func (b *Bitset64) Set(i int64) {
	idx := bs64idx(i)
	bit := bs64bit(i)
	b.bs[idx] |= bit
}

func (b *Bitset64) Clear(i int64) {
	b.bs[bs64idx(i)] &= ^bs64bit(i)
}

func (b *Bitset64) Test(i int64) bool {
	bit := bs64bit(i)
	return bit == (bit & b.bs[bs64idx(i)])
}

func (b Bitset64) Complement() Bitset64 {
	n := NewFrom(b)
	for i, bits := range n.bs {
		n.bs[i] = ^bits
	}
	return n
}

func (b Bitset64) Intersect(n Bitset64) Bitset64 {
	if n.k > b.k {
		b, n = n, b
	}

	ret := NewFrom(b)
	for i, d := range n.bs {
		ret.bs[i] &= d
	}

	return ret
}

func (b Bitset64) Union(n Bitset64) Bitset64 {
	if n.k > b.k {
		b, n = n, b
	}

	ret := NewFrom(b)
	for i, bits := range n.bs {
		ret.bs[i] |= bits
	}

	return ret
}

func (b Bitset64) Difference(n Bitset64) Bitset64 {
	return b.Intersect(n.Complement())
}

func (b Bitset64) Eq(n Bitset64) bool {
	if b.k != n.k {
		return false
	}

	pairs := bag.ZipTwo(b.bs, n.bs)

	for {
		p := pairs.Next()
		if p == nil {
			break
		}
		if p.A != p.B {
			return false
		}
	}

	return true
}

type Bitset64 struct {
	k  int64
	bs []int64
}

func (b Bitset64) String() string {
	return fmt.Sprintf("Bitset64 { k: %d }", b.k)
}

const bs64Size = 64

func bs64idx(i int64) int64 {
	return i / bs64Size
}

func bs64bit(i int64) int64 {
	return 1 << (i % bs64Size)
}
