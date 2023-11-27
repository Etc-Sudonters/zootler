package logic

import "context"

type Placer interface {
	Place(ctx context.Context, eg EntityGraph, tokens, locations []ComponentTupleIterator)
}
