package bitarrpool

import (
	"reflect"
	"sudonters/zootler/pkg/entity"
)

type componentId int

const (
	INVALID_COMPONENT componentId  = 0
	INVALID_ENTITY    entity.Model = 0
)

type componentTable struct {
	rows   []*componentRow
	lookup map[reflect.Type]componentId
}

func (tbl *componentTable) rowFor(c entity.Component) *componentRow {
	id, ok := tbl.idValue(c)
	if !ok {
		return tbl.addrow(reflect.TypeOf(c))
	}

	return tbl.rows[id]
}

func (tbl componentTable) idType(typ reflect.Type) (componentId, bool) {
	id, ok := tbl.lookup[typ]
	if !ok {
		return INVALID_COMPONENT, false
	}
	return id, true
}

func (tbl componentTable) idValue(c entity.Component) (componentId, bool) {
	val := reflect.Indirect(reflect.ValueOf(c))
	return tbl.idType(val.Type())
}

func (tbl componentTable) row(r componentId) *componentRow {
	return tbl.rows[r]
}

func (tbl *componentTable) addrow(t reflect.Type) *componentRow {
	row := make(componentRow, 1, 128)
	tbl.lookup[t] = componentId(len(tbl.rows))
	tbl.rows = append(tbl.rows, &row)
	return &row
}

func (tbl *componentTable) init() {
	tbl.rows = make([]*componentRow, 1, 128)
	tbl.lookup = make(map[reflect.Type]componentId, 128)
	tbl.lookup[nil] = 0
}

type componentRow []entity.Component // idx'd by entId

func (row *componentRow) set(e entity.Model, c entity.Component) {
	row.ensureSize(int(e))
	(*row)[e] = c
}

func (row *componentRow) unset(e entity.Model) {
	if len(*row) < int(e) {
		return
	}

	(*row)[e] = nil
}

func (row componentRow) get(e entity.Model) entity.Component {
	if len(row) < int(e) {
		return nil
	}

	return row[e]
}

func (row *componentRow) ensureSize(n int) {
	if len(*row) > n {
		return
	}

	expaded := make(componentRow, n+1, n*2)
	copy(expaded, *row)
	*row = expaded
}
