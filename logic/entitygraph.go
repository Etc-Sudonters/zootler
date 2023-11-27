package logic

import (
	"context"
)

type EntityGraph interface {
	Successors(ctx context.Context, start EntityId, strategy WalkStrategy, load Loader) (ComponentTupleCoroutine, error)
	Predecessors(ctx context.Context, start EntityId, strategy WalkStrategy, load Loader) (ComponentTupleCoroutine, error)
	Relate(ctx context.Context, src, dest EntityId) (EntityId, error)
	Select(ctx context.Context, load LoadSelector) (ComponentTupleIterator, error)
	Count(ctx context.Context, selector Selector) (int, error)
	Set(ctx context.Context, selector Selector, component Component) (ComponentTupleIterator, error)
	SetOne(ctx context.Context, entity EntityId, component Component) error
	New(ctx context.Context, components ...Component) (EntityId, error)
	Unique(ctx context.Context, uniqueComponent Component) (EntityId, error)
}

type WalkStrategy uint8

const (
	_ WalkStrategy = iota
	BreadthFirst
	DepthFirst
)

type Selector interface {
	With() []ComponentId
	Without() []ComponentId
	On() []EntityId
}

type Loader interface {
	Load() []ComponentId
}

type LoadSelector interface {
	Loader
	Selector
}

type ComponentTuple struct {
	Entity     EntityId
	Components []Component
}

type ComponentTupleIterator interface {
	Error() error
	Advance() bool
	Current() *ComponentTuple
	Length() int
}

type ComponentTupleCoroutine interface {
	ComponentTupleIterator
	Accept(EntityId)
}

type ComponentTupleIteratorOf[T any] interface {
	ComponentTupleIterator
	CastCurrent() *T
}

type ComponentView[T any] struct {
	Entity EntityId
	View   *T
}

func CastTuple[T any](*ComponentTuple) (ComponentView[T], error) {
	var t ComponentView[T]
	return t, nil
}

/*
	```go
	entities, _ := SelectAs[struct {
		Name string
		Qty  float64
	}](context.Background(), nil, nil)

	_, this := entities.CastCurrent()
	fmt.Print(this.Name)
	```
*/
func SelectAs[T any](context.Context, EntityGraph, Selector) (ComponentTupleIteratorOf[T], error) {
	return nil, nil
}
