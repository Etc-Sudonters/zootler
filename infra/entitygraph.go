package infra

import (
	"context"
	"sudonters/zootler/logic"
	"sudonters/zootler/storage"
)

func NewEntityGraph(db *storage.Database) EntityGraph {
	return EntityGraph{db}
}

type RunnableEntityGraph interface {
	logic.EntityGraph
	Run(ctx context.Context, program *storage.Program) (logic.ComponentTupleCoroutine, error)
}

type EntityGraph struct {
	db *storage.Database
}

func (eg EntityGraph) Run(ctx context.Context, program *storage.Program) (logic.ComponentTupleCoroutine, error) {
	panic("not implemented") // TODO: Implement
}

func (eg EntityGraph) Successors(ctx context.Context, start logic.EntityId, strategy logic.WalkStrategy, load logic.Loader) (logic.ComponentTupleCoroutine, error) {
	panic("not implemented") // TODO: Implement
}

func (eg EntityGraph) Predecessors(ctx context.Context, start logic.EntityId, strategy logic.WalkStrategy, load logic.Loader) (logic.ComponentTupleCoroutine, error) {
	panic("not implemented") // TODO: Implement
}

func (eg EntityGraph) Relate(ctx context.Context, src logic.EntityId, dest logic.EntityId) (logic.EntityId, error) {
	panic("not implemented") // TODO: Implement
}

func (eg EntityGraph) Select(ctx context.Context, load logic.LoadSelector) (logic.ComponentTupleIterator, error) {
	panic("not implemented") // TODO: Implement
}

func (eg EntityGraph) Count(ctx context.Context, selector logic.Selector) (int, error) {
	panic("not implemented") // TODO: Implement
}

func (eg EntityGraph) Set(ctx context.Context, selector logic.Selector, component logic.Component) (logic.ComponentTupleIterator, error) {
	panic("not implemented") // TODO: Implement
}

func (eg EntityGraph) SetOne(ctx context.Context, entity logic.EntityId, component logic.Component) error {
	panic("not implemented") // TODO: Implement
}

func (eg EntityGraph) New(ctx context.Context, components ...logic.Component) (logic.EntityId, error) {
	panic("not implemented") // TODO: Implement
}

func (eg EntityGraph) Unique(ctx context.Context, uniqueComponent logic.Component) (logic.EntityId, error) {
	panic("not implemented") // TODO: Implement
}

func NewSelectorBuilder() *LoadSelectorBuilder {
	return new(LoadSelectorBuilder)
}

type LoadSelectorBuilder struct {
	load    []logic.ComponentId
	on      []logic.EntityId
	with    []logic.ComponentId
	without []logic.ComponentId
}

func (ls *LoadSelectorBuilder) Load(c logic.ComponentId) *LoadSelectorBuilder {
	ls.load = append(ls.load, c)
	return ls
}

func (ls *LoadSelectorBuilder) With(c logic.ComponentId) *LoadSelectorBuilder {
	ls.with = append(ls.with, c)
	return ls
}

func (ls *LoadSelectorBuilder) Without(c logic.ComponentId) *LoadSelectorBuilder {
	ls.without = append(ls.without, c)
	return ls
}

func (ls *LoadSelectorBuilder) On(e logic.EntityId) *LoadSelectorBuilder {
	ls.on = append(ls.on, e)
	return ls
}

func (ls LoadSelectorBuilder) Build() LoadSelector {
	return LoadSelector{
		load:    ls.load,
		on:      ls.on,
		with:    ls.with,
		without: ls.without,
	}
}

type LoadSelector struct {
	load    []logic.ComponentId
	on      []logic.EntityId
	with    []logic.ComponentId
	without []logic.ComponentId
}

func (ls LoadSelector) Load() []logic.ComponentId {
	return ls.load
}

func (ls LoadSelector) With() []logic.ComponentId {
	return ls.with
}

func (ls LoadSelector) Without() []logic.ComponentId {
	return ls.without
}

func (ls LoadSelector) On() []logic.EntityId {
	return ls.on
}
