package storage

import "context"

type Database struct {
	vm        *VirtualMachine
	allocator Allocator
	documents DocumentEngine
	graph     GraphEngine
	columns   ColumnEngine
}

func (s Database) Run(context.Context, *Program) (TupleIterator, error) {
	return nil, nil
}

func New() (*Database, error) {
	return nil, nil
}
