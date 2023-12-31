package componenttable

import (
	"errors"
	"fmt"
	"reflect"
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/skelly/bitset"
	"github.com/etc-sudonters/substrate/stageleft"
)

var strType reflect.Type = mirrors.TypeOf[string]()
var ErrCorruptedTable = errors.New("table became corrupted")

type componentGetter struct {
	*Table
}

func (c componentGetter) GetComponent(e entity.Model, t reflect.Type) (entity.Component, error) {
	return c.Get(e, t)
}

func New(maxEntities int) *Table {
	t := new(Table)
	t.rows = make([]*Row, 1, 32)
	t.typemap = make(mirrors.TypeMap, 32)
	t.rows[entity.INVALID_ENTITY] = nil
	t.typemap[nil] = mirrors.TypeId(entity.INVALID_COMPONENT)
	t.getter = componentGetter{t}
	t.entityBuckets = bitset.Buckets(maxEntities)
	return t
}

type Table struct {
	entityBuckets int
	rows          []*Row
	typemap       mirrors.TypeMap
	getter        entity.ComponentGetter
}

func (t *Table) Set(e entity.Model, c entity.Component) entity.ComponentId {
	typ := entity.PierceComponentType(c)
	if typ == strType {
		panic(fmt.Errorf("string component added to %d: %q", e, c))
	}
	row := t.RowOf(typ)
	row.Set(e, c)
	return row.id
}

func (t *Table) Unset(e entity.Model, typ reflect.Type) entity.ComponentId {
	if r := t.rowFor(typ); r != nil {
		r.Unset(e)
		return r.id
	}
	return entity.INVALID_COMPONENT
}

func (t *Table) IdOf(typ reflect.Type) (entity.ComponentId, error) {
	id, err := t.typemap.IdOf(typ)
	if err != nil {
		return 0, entity.ErrUnknownComponent{T: typ}
	}
	return entity.ComponentId(id), nil
}

func (t *Table) Get(e entity.Model, typ reflect.Type) (entity.Component, error) {
	r := t.rowFor(typ)
	if r == nil {
		return nil, entity.ErrNotAssigned
	}

	c := r.Get(e)
	if c == nil {
		return nil, entity.ErrNotAssigned
	}

	return c, nil
}

func (t Table) Getter() entity.ComponentGetter {
	return t.getter
}

func (t Table) Len() int {
	return len(t.rows)
}

func (t Table) rowFor(typ reflect.Type) *Row {
	id, err := t.typemap.IdOf(typ)
	if err != nil {
		return nil
	}

	if len(t.rows) > int(id) {
		return t.rows[int(id)]
	}

	return nil
}

func (t *Table) RowOf(typ reflect.Type) *Row {
	if r := t.rowFor(typ); r != nil {
		return r
	}

	id := t.typemap.Add(typ)

	if len(t.rows) != int(id) {
		panic(stageleft.AttachExitCode(ErrCorruptedTable, stageleft.ExitCode(117)))
	}

	r := new(Row)
	r.Init(entity.ComponentId(id), t.entityBuckets)
	t.rows = append(t.rows, r)
	r.typ = typ
	return r
}
