package world

import (
	"context"
	"errors"

	"sudonters/zootler/internal/bag"
	"sudonters/zootler/internal/ioutil"
	"sudonters/zootler/internal/queue"
	"sudonters/zootler/pkg/entity"
	"sudonters/zootler/pkg/logic"
)

type ConstGoal bool

func (c ConstGoal) Reachable(context.Context, World) (bool, error) {
	return bool(c), nil
}

type Goal interface {
	Reachable(context.Context, World) (bool, error)
}

type Filler interface {
	Fill(context.Context, World, Goal) error
}

type AssumedFill struct {
	Locations []entity.Selector
	Items     []entity.Selector
}

func (a *AssumedFill) Fill(ctx context.Context, w World, g Goal) error {
	var err error = nil
	var locs []entity.View
	var items []entity.View

	var filt []entity.Selector

	filt = make([]entity.Selector, len(a.Locations)+1)
	filt[0] = entity.With[logic.Location]{}
	copy(filt[1:], a.Locations)

	locs, err = w.Entities.Pool.Query(filt)
	if err != nil {
		return ioutil.AttachExitCode(err, ioutil.ExitQueryFail)
	}

	filt = make([]entity.Selector, len(a.Items)+1)
	filt[0] = entity.With[logic.Token]{}
	copy(filt[1:], a.Items)
	items, err = w.Entities.Pool.Query(filt)
	if err != nil {
		return ioutil.AttachExitCode(err, ioutil.ExitQueryFail)
	}

	L := queue.From(locs)
	I := queue.From(items)

	var solved bool
	maxTries := len(L) * len(I)

	for i := 0; i <= maxTries; i++ {
		if len(L) == 0 || len(I) == 0 {
			break
		}

		if err != nil {
			return err
		}

		var loc entity.View
		var item entity.View

		bag.Shuffle(L)
		bag.Shuffle(I)
		loc, L, err = L.Pop()
		if err != nil {
			return err
		}
		item, I, err = I.Pop()
		if err != nil {
			return err
		}

		loc.Add(logic.Inhabited(item.Model()))
		item.Add(logic.Inhabits(loc.Model()))

		solved, err = g.Reachable(ctx, w)
		if err != nil {
			return err
		}

		if !solved {
			loc.Remove(logic.Inhabited(item.Model()))
			item.Remove(logic.Inhabits(loc.Model()))
			L = L.Push(loc)
			I = I.Push(item)
		}
	}

	if !solved {
		err = errors.New("could not solve placement")
	}

	return err
}
