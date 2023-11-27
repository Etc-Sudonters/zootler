package main

import (
	"context"
	"fmt"
	"strings"
	"sudonters/zootler/infra"
	"sudonters/zootler/logic"
	"sudonters/zootler/storage"
)

func main() {
	ctx := context.TODO()
	store, _ := storage.New()
	var eg infra.RunnableEntityGraph = infra.NewEntityGraph(store)

	iter, _ := eg.Run(
		ctx,
		panicWhenErr(func() (*storage.Program, error) {
			return storage.CompileStr("SELECT (Name, Price) WITH (Scrub, Grotto)")
		}))

	for elems := iter.Current(); iter.Advance() && iter.Error() == nil; {
		var trace strings.Builder
		grottoScrub, _ := logic.CastTuple[struct {
			Name  NameComponent
			Price PriceComponent
		}](elems)

		pathHere, _ := eg.Predecessors(
			ctx, grottoScrub.Entity, logic.DepthFirst,
			infra.NewSelectorBuilder().
				Load(NameComponentId).
				Load(SpawnComponentId).
				With(LocationComponentId).
				Build(),
		)

		trace.WriteString(string(grottoScrub.View.Name))

		for parent := pathHere.Current(); pathHere.Advance() && pathHere.Error() == nil; {
			segment, _ := logic.CastTuple[struct {
				Name  NameComponent
				Spawn *SpawnComponent
			}](parent)

			fmt.Fprintf(&trace, " <- %s", segment.View.Name)

			if segment.View.Spawn == nil {
				pathHere.Accept(segment.Entity)
			}
		}
	}
}

const (
	NameComponentId logic.ComponentId = iota
	PriceComponentId
	SpawnComponentId
	LocationComponentId
)

type NameComponent string
type PriceComponent float64
type SpawnComponent struct{}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func panicWhenErr[T any](f func() (T, error)) T {
	t, err := f()
	panicIfErr(err)
	return t
}
