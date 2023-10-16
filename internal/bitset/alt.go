package bitset

import (
	"fmt"
	"sudonters/zootler/internal/bag"
)

func IsFieldEmpty[F Field](b BitSet[F]) bool {
	for i := range b.field {
		if b.field[i] != 0 {
			return false
		}
	}
	return true
}

type Field interface {
	uint8 | uint16 | uint32 | uint64
}

func BucketsNeeded[F Field](n int) int {
	return (n / FieldSize[F]()) + 1
}

func FieldSize[F Field]() int {
	var f F

	switch any(f).(type) {
	case uint8:
		return 8
	case uint16:
		return 16
	case uint32:
		return 32
	case uint64:
		return 64
	default:
		panic("unreachable")
	}
}

func Copy[F Field](b BitSet[F]) (c BitSet[F]) {
	c.bits = b.bits
	c.k = b.k
	c.field = make([]F, len(b.field))
	copy(c.field, b.field)
	return
}

func WithCapacity[F Field](n int) (b BitSet[F]) {
	b.bits = FieldSize[F]()
	b.k = 1 + n/b.bits
	b.field = make([]F, b.k)
	return
}

func WithBuckets[F Field](k int) (b BitSet[F]) {
	b.bits = FieldSize[F]()
	b.k = k
	b.field = make([]F, k)
	return b
}

func Saturated[F Field](buckets int) BitSet[F] {
	return WithCapacity[F](buckets).Complement()
}

type BitSet[F Field] struct {
	k     int
	bits  int
	field []F
}

func (b BitSet[F]) String() string {
	return fmt.Sprintf("bitset[%d]{ k: %d }", b.bits, b.k)
}

func (b BitSet[F]) bucketFor(n int) int {
	return n / b.bits
}

func (b BitSet[F]) mask(n int) F {
	return F(1 << (n % b.bits))
}

func (b BitSet[F]) Set(n int) {
	b.field[b.bucketFor(n)] |= b.mask(n)
}

func (b BitSet[F]) Clear(n int) {
	b.field[b.bucketFor(n)] &= ^b.mask(n)
}

func (b BitSet[F]) Test(n int) bool {
	mask := b.mask(n)
	return mask == (mask & b.field[b.bucketFor(n)])
}

func (b BitSet[F]) Complement() BitSet[F] {
	c := Copy(b)
	for i, bucket := range c.field {
		c.field[i] = ^bucket
	}
	return c
}

func (b BitSet[F]) Intersect(bb BitSet[F]) (c BitSet[F]) {
	if bb.k > b.k {
		b, bb = bb, b
	}

	c.k = b.k
	c.bits = b.bits
	c.field = make([]F, c.k)

	for i := range bb.field {
		c.field[i] = b.field[i] & bb.field[i]
	}

	return
}

func (b BitSet[F]) Difference(bb BitSet[F]) BitSet[F] {
	return b.Intersect(bb.Complement())
}

func (b BitSet[F]) Eq(n BitSet[F]) bool {
	if b.k != n.k {
		return false
	}

	pairs := bag.ZipTwo(b.field, n.field)

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
