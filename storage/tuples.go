package storage

type Tuple struct {
	Id    uint64
	Items []interface{}
}

type TupleIterator interface {
	Length() int
	Current() *Tuple
	Advance() bool
}
