package bitpool

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"sudonters/zootler/internal/bag"
	"sudonters/zootler/internal/bitset"
	"sudonters/zootler/pkg/entity"
)

var ErrNoMoreIds = errors.New("no more ids available")
var ErrNoEntities = errors.New("no entities")
var _ entity.Pool = (*bitpool)(nil)
var _ entity.View = (*bitview)(nil)

func New(maxId int64) *bitpool {
	buckets := bitset.Buckets(maxId)
	b := bitset.New(buckets)
	return &bitpool{
		k:      buckets,
		maxId:  entity.Model(maxId),
		nextId: 0,
		all:    &b,
		tests:  make(map[reflect.Type]*bitset.Bitset64),
		table:  make(map[reflect.Type]map[entity.Model]entity.Component),
	}
}

func unknown(t reflect.Type, b *bitpool) error {
	return fmt.Errorf("unknown component: %q\navailable:\n%s", bag.NiceTypeName(t), bpSummary{*b})
}

// return a subset of the population that matches the provided selectors
func (b *bitpool) Query(qs []entity.Selector) ([]entity.View, error) {
	if b.nextId == 0 {
		return nil, ErrNoEntities
	}

	subset := bitset.NewFrom(*b.all)

	for _, q := range qs {
		test := q.Component()
		behavior := q.Behavior()
		members, ok := b.tests[test]
		if !ok {
			return nil, unknown(test, b)
		}

		switch behavior {
		case entity.ComponentInclude:
			subset = subset.Intersect(*members)
		case entity.ComponentExclude:
			subset = subset.Difference(*members)
		default:
			panic(fmt.Errorf("unknown behavior %q", &behavior))
		}

		if bitset.IsFieldEmpty(subset) {
			return nil, ErrNoEntities
		}
	}

	var entities []entity.View

	for i := int64(0); i < int64(b.nextId); i++ {
		if subset.Test(i) {
			e := bitview{entity.Model(i), b}
			entities = append(entities, &e)
		}
	}

	return entities, nil
}

func (b *bitpool) Get(m entity.Model, cs ...interface{}) {
	for i := range cs {
		_ = entity.AssignComponentTo(cs[i], tryFindCompFor(b, m))
	}
}

func (b *bitpool) Create() (entity.View, error) {
	if b.maxId == b.nextId {
		return nil, ErrNoMoreIds
	}
	v := bitview{id: b.nextId, p: b}
	b.nextId++
	b.all.Set(int64(v.id))
	_ = v.Add(v.id)
	return &v, nil
}

type bitpool struct {
	k      int64
	maxId  entity.Model
	nextId entity.Model
	all    *bitset.Bitset64
	tests  map[reflect.Type]*bitset.Bitset64
	table  map[reflect.Type]map[entity.Model]entity.Component
}

type bpSummary struct {
	bitpool
}

func (b bpSummary) String() string {
	return fmt.Sprintf("bitpool{ k: %d, maxId: %d, nextId: %d, all: %s, tests: %s }", b.k, b.maxId, b.nextId, b.all, b.summarizeTable())
}

func (b bpSummary) summarizeTable() string {
	repr := &strings.Builder{}
	repr.WriteRune('{')
	repr.WriteRune(' ')
	for typ, _ := range b.tests {
		fmt.Fprintf(repr, "%q", typ)
		repr.WriteRune(' ')
	}
	repr.WriteRune('}')
	return repr.String()
}

type bitview struct {
	id entity.Model
	p  *bitpool
}

func (b bitview) String() string {
	return fmt.Sprintf("bitview(%d)", b.id)
}

func (b *bitview) Model() entity.Model {
	return entity.Model(b.id)
}

func (b *bitview) Get(w interface{}) error {
	return entity.AssignComponentTo(w, tryFindCompFor(b.p, b.id))
}

func (b *bitview) AddMany(cs ...entity.Component) error {
	for _, c := range cs {
		if err := b.Add(c); err != nil {
			return err
		}
	}
	return nil
}

func (b *bitview) Add(c entity.Component) error {
	if many, ok := c.([]entity.Component); ok {
		return b.AddMany(many...)
	}

	typ := reflect.TypeOf(c)
	ensureTables(b.p, typ)
	b.p.tests[typ].Set(int64(b.id))
	b.p.table[typ][b.id] = c
	return nil
}

func (b *bitview) Remove(c entity.Component) error {
	typ := reflect.TypeOf(c)
	members, ok := b.p.tests[typ]
	if !ok {
		return nil
	}

	members.Clear(int64(b.id))
	delete(b.p.table[typ], b.id)
	return nil
}

func ensureTables(b *bitpool, t reflect.Type) {
	if t == nil {
		panic(errors.New("nil component type"))
	}

	if _, ok := b.tests[t]; !ok {
		set := bitset.New(b.k)
		b.tests[t] = &set
		b.table[t] = make(map[entity.Model]entity.Component)
	}
}

func existsInCompTable(b *bitpool, t reflect.Type, m entity.Model) bool {
	tab, ok := b.tests[t]
	if !ok {
		return false
	}
	return tab.Test(int64(m))
}

func tryFindCompFor(p *bitpool, m entity.Model) func(reflect.Type) (entity.Component, error) {
	return func(t reflect.Type) (entity.Component, error) {
		if !existsInCompTable(p, t, m) {
			return nil, entity.ErrNotAssigned
		}
		return p.table[t][m], nil
	}
}
