package bitarrpool

import (
	"fmt"
	"reflect"
	"sudonters/zootler/pkg/entity"

	"github.com/etc-sudonters/substrate/skelly/set/bits"
)

type implSet = bits.Bitset64

func newSet(k int) implSet {
	return bits.New(k)
}

func implBitSize() int {
	return 64
}

func getComponenter(b bitarrview) func(reflect.Type) (entity.Component, error) {
	return func(t reflect.Type) (entity.Component, error) {
		id, ok := b.p.table.lookup[t]
		if !ok {
			return nil, entity.ErrUnknownComponent{Type: t}
		}

		if !b.comps.Test(int(id)) {
			return nil, entity.ErrNotAssigned
		}

		row := b.p.table.row(id)
		comp := row.get(b.id)
		if comp == nil {
			return nil, entity.ErrNilComponent{
				Entity:    b.Model(),
				Component: t,
			}
		}

		return comp, nil
	}
}

type CompressedRepr struct {
	bitarrpool
}

type compressedTableRepr struct {
	componentTable
}

func (c compressedTableRepr) String() string {
	return fmt.Sprintf("componentTable{ rows: %d }", len(c.rows))
}

func (c CompressedRepr) String() string {
	return fmt.Sprintf(
		"bitarrpool{ k: %d, entities: %d, table: %s }",
		c.componentBucketCount,
		len(c.entities),
		compressedTableRepr{c.table},
	)
}
