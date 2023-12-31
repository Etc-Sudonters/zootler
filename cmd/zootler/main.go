package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"sudonters/zootler/cmd/zootler/tui"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/filler"
	"sudonters/zootler/pkg/world"
	"sudonters/zootler/pkg/world/components"
	"sudonters/zootler/pkg/world/filter"

	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/mirrors"
	"github.com/etc-sudonters/substrate/stageleft"
)

type cliOptions struct {
	logicDir   string `short:"-l" description:"Path to logic files" required:"t"`
	dataDir    string `short:"-d" description:"Path to data files" required:"t"`
	visualizer bool   `short:"-v" description:"Open visualizer" required:"f"`
}

func (opts *cliOptions) init() {
	flag.StringVar(&opts.logicDir, "l", "", "Directory where logic files are located")
	flag.StringVar(&opts.dataDir, "d", "", "Directory where data files are stored")
	flag.BoolVar(&opts.visualizer, "v", false, "Open visualizer")
	flag.Parse()
}

func (c cliOptions) validate() error {
	if c.logicDir == "" {
		return missingRequired("-l")
	}

	if c.dataDir == "" {
		return missingRequired("-d")
	}

	return nil
}

func main() {
	var opts cliOptions
	var exit stageleft.ExitCode = stageleft.ExitSuccess
	stdio := dontio.Std{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stdout,
	}
	defer func() {
		os.Exit(int(exit))
	}()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				fmt.Fprintf(stdio.Err, "%s\n", err)
			}
			_, _ = stdio.Err.Write(debug.Stack())
			if exit != stageleft.ExitSuccess {
				exit = stageleft.AsExitCode(r, stageleft.ExitCode(126))
			}
		}
	}()

	ctx := context.Background()
	ctx = dontio.AddStdToContext(ctx, &stdio)

	(&opts).init()

	if cliErr := opts.validate(); cliErr != nil {
		fmt.Fprintf(stdio.Err, "%s\n", cliErr.Error())
		exit = stageleft.ExitCode(2)
		return
	}

	b := world.DefaultBuilder()

	stampTokens(b)

	w := b.Build()

	if opts.visualizer {
		v := tui.Tui(w)
		if err := v.Run(ctx); err != nil {
			panic(err)
		}
		return
	}

	song := entity.FilterBuilder{}.With(mirrors.TypeOf[components.Song]())

	assumed := &filler.AssumedFill{
		Locations: entity.BuildFilter(filter.Song),
		Items:     entity.BuildFilter(filter.Song),
	}
	if err := assumed.Fill(ctx, w, filler.ConstGoal(true)); err != nil {
		exit = stageleft.ExitCodeFromErr(err, stageleft.ExitCode(2))
		fmt.Fprintf(stdio.Err, "Error during placement: %s\n", err.Error())
		return
	}

	if err := showTokenPlacements(ctx, w, song.Clone()); err != nil {
		exit = stageleft.ExitCodeFromErr(err, stageleft.ExitCode(2))
		fmt.Fprintf(stdio.Err, "Error during placement review: %s\n", err.Error())
		return
	}
}

type missingRequired string // option name

func (arg missingRequired) Error() string {
	return fmt.Sprintf("%s is required", string(arg))
}

func showTokenPlacements(ctx context.Context, w world.World, fb entity.FilterBuilder) error {
	placed, err := w.Entities.Query(entity.BuildFilter(filter.Inhabits).Combine(fb).Build())
	if err != nil {
		return fmt.Errorf("while querying placements: %w", err)
	}
	stdio, _ := dontio.StdFromContext(ctx)

	for _, tok := range placed {
		var itemName components.Name
		var placementName components.Name
		var placement components.Inhabits

		err = tok.Get(&itemName)
		if err != nil {
			return err
		}
		err = tok.Get(&placement)
		if err != nil {
			return err
		}
		w.Entities.Get(entity.Model(placement), []interface{}{&placementName})
		if placementName == "" {
			return fmt.Errorf("%v did not have an attached name", placement)
		}

		fmt.Fprintf(stdio.Out, "%s placed at %s", itemName, placementName)
	}

	return nil
}

func stampTokens(b *world.Builder) {
	tokens, err := b.Pool.Query(entity.FilterBuilder{}.With(mirrors.TypeOf[components.Token]()).Build())
	if err != nil {
		panic(err)
	}

	var name components.Name

	for _, token := range tokens {
		token.Get(&name)
		stamp := b.TypedStrs.Typed(string(name))
		token.Add(stamp)
	}
}
