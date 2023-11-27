package storage

import "context"

type Program struct {
	ops []Op
}

type VirtualMachine struct {
}

func (vm *VirtualMachine) Run(context.Context, Program) (TupleIterator, error) {
	return nil, nil
}
