package logic

import "context"

type Solver interface {
	Solve(ctx context.Context, eg EntityGraph)
}
