package bitarrpool

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"sudonters/zootler/internal/errs"
	"sudonters/zootler/internal/set"
	"sudonters/zootler/pkg/entity"
)

func dump(t *testing.T, v interface{}) {
	t.Logf("%+v", v)
}

func dumpView(w io.Writer, v interface{}) {
	fmt.Fprintf(w, "Dump:\n%+v\n", v)
}

func expectedComponent(w io.Writer, v bitarrview, expected entity.Component) {
	fmt.Fprintf(w, "expected component %s to be loaded on %d\n", expected, v.id)
}

func didNotExpectError(t *testing.T, err error) {
	w := &strings.Builder{}
	fmt.Fprintf(w, "did not expected error:\n%s", err)
	t.Fatal(w.String(), errs.ShowPanicTrace())
}

func expectedEqualComponents[T entity.Component](w io.Writer, v bitarrview, expect T, actual T) {
	fmt.Fprintf(w, "expected equal values for component %s on %d\n", entity.ComponentName(expectedComponent), v.id)
	fmt.Fprint(w, "Actual ")
	dumpView(w, actual)
	fmt.Fprint(w, "Expected ")
	dumpView(w, expect)
}

func TestCanStoreAndRetrievePointerToComp(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatal(p)
		}
	}()
	initialValue := 10
	changedValue := 9999

	p := New(100)
	v, _ := p.Create()
	v.Add(myTestComponent{-1})
	v.Add(&myTestComponent{initialValue})

	var c *myTestComponent
	// ptrB -> ptrA, assigns ptrA
	if err := v.Get(&c); err != nil {
		didNotExpectError(t, err)
	}

	c.V = changedValue

	var d *myTestComponent
	if err := v.Get(&d); err != nil {
		didNotExpectError(t, err)
	}

	if d.V != c.V {
		msg := &strings.Builder{}
		expectedEqualComponents(msg, v.(bitarrview), c, d)
		t.Fatal(msg.String())
	}
}

func TestCanRemoveCustomComponent(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatal(p)
		}
	}()
	p := New(1000)
	v, _ := p.Create()
	ent := v.(bitarrview)
	ent.Add(myTestComponent{})

	var c myTestComponent
	if err := ent.Get(&c); err != nil {
		didNotExpectError(t, err)
	}

	// doesn't need to be loaded instance
	ent.Remove(myTestComponent{})

	var d myTestComponent
	if err := ent.Get(&d); err != nil && !errors.Is(err, entity.ErrNotAssigned) {
		didNotExpectError(t, err)
	}
}

func TestCanQueryForEntitiesByComponentExistence(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatal(p)
		}
	}()
	entitiesToMake := 10000
	tagRatio := 7
	p := New(10000)

	totalEnts := set.New[entity.Model]()
	taggedEnts := set.New[entity.Model]()

	for i := 0; i < entitiesToMake; i++ {
		v, _ := p.Create()
		ent := v.(bitarrview)
		totalEnts.Add(ent.Model())

		if (i % tagRatio) == 0 {
			ent.Add(myTestComponent{i})
			taggedEnts.Add(ent.Model())
		}
	}

	if len(totalEnts) != entitiesToMake {
		t.Fatalf("mismatched entity count\nexpected:\t%d\nactual:\t%d", len(totalEnts), entitiesToMake)
	}

	filter := []entity.Selector{entity.With[myTestComponent]{}}
	queryedFor, err := p.Query(filter)
	if err != nil {
		didNotExpectError(t, err)
	}

	if len(queryedFor) != len(taggedEnts) {
		t.Logf("mismatched entity count\nexpected:\t%d\nactual:\t%d", len(taggedEnts), len(queryedFor))
		t.FailNow()
	}

	if !reflect.DeepEqual(
		taggedEnts,
		set.MapFromSlice(
			queryedFor,
			func(u entity.View) entity.Model {
				return u.Model()
			},
		),
	) {
		t.Logf("tagged set and returned set are different!")
		t.Fail()
	}
}

func TestCanUseMultipleComponents(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatal(p)
		}
	}()
	componentsToMake := 35
	goodTagRatio := 7
	badTagRation := 5
	p := New(1000)

	totalEnts := set.New[entity.Model]()
	goodTaggedEnts := set.New[entity.Model]()
	badTaggedEnts := set.New[entity.Model]()

	for i := 0; i <= componentsToMake; i++ {
		v, _ := p.Create()
		ent := v.(bitarrview)
		totalEnts.Add(ent.Model())

		if (i % goodTagRatio) == 0 {
			ent.Add(myTestComponent{i})
			goodTaggedEnts.Add(ent.Model())
		}

		if (i % badTagRation) == 0 {
			ent.Add(anotherComponent{float64(i)})
			badTaggedEnts.Add(ent.Model())
		}
	}

	comboTagSet := set.Intersection(goodTaggedEnts, badTaggedEnts)
	comboQueries, err := p.Query([]entity.Selector{
		entity.DebugSelector{
			F:        func(string, ...any) {}, //t.Logf,
			Selector: entity.With[myTestComponent]{},
		},
		entity.DebugSelector{
			F:        t.Logf,
			Selector: entity.With[anotherComponent]{},
		}},
	)

	if err != nil {
		didNotExpectError(t, err)
	}

	actualEntities := set.MapFromSlice(comboQueries, entity.View.Model)

	if !reflect.DeepEqual(
		comboTagSet, actualEntities,
	) {
		t.Log("selected a different set of entities than expected")
		t.Logf("expected:\t%+v", comboTagSet)
		t.Logf("actual:\t\t%+v", actualEntities)
		t.Fail()
	}
}

func TestCanExcludeEntitiesBasedOnComponent(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatal(p)
		}
	}()
	componentsToMake := 1001
	firstTagRatio := 7
	secondTagRatio := 5
	p := New(int(componentsToMake))

	totalEnts := set.New[entity.Model]()
	firstTagEnts := set.New[entity.Model]()
	secondTagEnt := set.New[entity.Model]()

	for i := 1; i <= componentsToMake; i++ {
		ent, _ := p.Create()

		totalEnts.Add(ent.Model())

		if (i % firstTagRatio) == 0 {
			ent.Add(myTestComponent{i})
			firstTagEnts.Add(ent.Model())
		}

		if (i % secondTagRatio) == 0 {
			ent.Add(anotherComponent{float64(i)})
			secondTagEnt.Add(ent.Model())
		}
	}

	allEntitiesWithoutTags := set.Difference(
		totalEnts,
		set.Union(firstTagEnts, secondTagEnt),
	)

	debuggify := func(q entity.Selector) entity.DebugSelector {
		return entity.DebugSelector{
			Selector: q,
			F:        t.Logf,
		}
	}

	queriedAllUntagged, err := p.Query([]entity.Selector{
		debuggify(entity.Without[myTestComponent]{}),
		debuggify(entity.Without[anotherComponent]{}),
	})

	if err != nil {
		didNotExpectError(t, err)
	}

	actualEntities := set.MapFromSlice(queriedAllUntagged, entity.View.Model)

	if actualEntities.Exists(entity.Model(firstTagRatio*secondTagRatio)) ||
		actualEntities.Exists(entity.Model(firstTagRatio)) ||
		actualEntities.Exists(entity.Model(secondTagRatio)) {
		t.Log("selected incorrect entity group")
		t.FailNow()
	}

	if !reflect.DeepEqual(allEntitiesWithoutTags, actualEntities) {
		t.Log("selected a different set of entities than expected")
		t.Logf("expected:\t%+v", allEntitiesWithoutTags)
		t.Logf("actual:\t\t%+v", actualEntities)
		t.Fail()
	}

}

func TestCanFilterWithoutLoading(t *testing.T) {
	var lastI = 0
	defer func() {
		if p := recover(); p != nil {
			t.Fatal(p, errs.ShowPanicTrace(), lastI)
		}
	}()
	entitiesToMake := 10000
	tagRatio := 7
	p := New(10)

	totalEnts := set.New[entity.Model]()
	taggedEnts := set.New[entity.Model]()
	taggedCount := 0

	for i := 1; i <= entitiesToMake; i++ {
		lastI = i
		v, _ := p.Create()
		ent := v.(bitarrview)
		totalEnts.Add(ent.Model())

		if (i % tagRatio) == 0 {
			ent.Add(myTestComponent{i})
			taggedEnts.Add(ent.Model())
			taggedCount += 1
		}
	}

	queriedEnts, err := p.Query([]entity.Selector{entity.Without[myTestComponent]{}})

	if err != nil {
		didNotExpectError(t, err)
	}

	var comp myTestComponent
	for _, ent := range queriedEnts {
		err := ent.Get(&comp)
		if errors.Is(err, entity.ErrNotAssigned) {
			continue
		}

		if err != nil {
			didNotExpectError(t, err)
		}

		t.Logf("did not expect %T to be loaded on %d", comp, ent.Model())
		t.Fail()
	}

	if len(queriedEnts) != (entitiesToMake - taggedCount) {
		t.Logf("expected %d elements", entitiesToMake-taggedCount)
		t.Logf("got %d elements", len(queriedEnts))
		t.Fail()
	}
}

func TestCanRetrieveArbitraryEntityWithComps(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatal(p, errs.ShowPanicTrace())
		}
	}()
	p := New(1000)
	v, _ := p.Create()
	ent := v.(bitarrview)
	ent.Add(myTestComponent{99})

	var c1 myTestComponent
	var c2 *anotherComponent

	p.Get(ent.id, &c1, &c2)

	if c1.V != 99 {
		t.Logf("did not retrieve expected instance of %[1]T: %[1]v", c1)
		t.Fail()
	}

	if c2 != nil {
		t.Logf("did not expect to retrieve %T from %s", c2, ent)
		t.Fail()
	}

}

type myTestComponent struct {
	V int
}

type anotherComponent struct {
	K float64
}
