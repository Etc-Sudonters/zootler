package componenttable

import (
	"reflect"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/reitertools"

	"github.com/etc-sudonters/substrate/skelly/set/bits"
)

type Row struct {
	id         entity.ComponentId
	typ        reflect.Type
	components []entity.Component
	members    bits.Bitset64
}

func (r *Row) Components() reitertools.Iterator[RowEntry] {
	i := reitertools.SubsliceIter(r.components, 1)
	f := reitertools.Filter(i, func(c entity.Component, _ int) bool {
		return c != nil
	})
	m := reitertools.Map(f, func(c entity.Component, idx int) RowEntry {
		return RowEntry{
			Entity:    entity.Model(idx),
			Component: c,
		}
	})
	return m
}

func (r *Row) Init(id entity.ComponentId, entityBuckets int) {
	r.id = id
	r.components = make([]entity.Component, 0)
	r.members = bits.New(entityBuckets)
}

func (row *Row) Set(e entity.Model, c entity.Component) {
	row.EnsureSize(int(e))
	row.components[e] = c
	row.members.Set(int(e))
}

func (row *Row) Unset(e entity.Model) {
	if len(row.components) < int(e) {
		return
	}

	row.components[e] = nil
	row.members.Clear(int(e))
}

func (row Row) Get(e entity.Model) entity.Component {
	if !row.members.Test(int(e)) {
		return nil
	}

	return row.components[e]
}

func (row *Row) EnsureSize(n int) {
	if len(row.components) > n {
		return
	}

	expaded := make([]entity.Component, n+1, n*2)
	copy(expaded, row.components)
	row.components = expaded
}