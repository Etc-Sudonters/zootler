package interpreter

import (
	"fmt"
	"sudonters/zootler/pkg/rules/ast"

	"github.com/etc-sudonters/substrate/stageleft"
)

type Evaluation[T any] interface {
	EvalBinOp(op *ast.BinOp, env Environment) T
	EvalBoolOp(op *ast.BoolOp, env Environment) T
	EvalCall(call *ast.Call, env Environment) T
	EvalIdentifier(ident *ast.Identifier, env Environment) T
	EvalLiteral(str *ast.Literal, env Environment) T
	EvalSubscript(subscript *ast.Subscript, env Environment) T
	EvalTuple(tup *ast.Tuple, env Environment) T
	EvalUnary(unary *ast.UnaryOp, env Environment) T
}

func Evaluate[T any](v Evaluation[T], node ast.Expression, env Environment) T {
	switch node := node.(type) {
	case *ast.BinOp:
		return v.EvalBinOp(node, env)
	case *ast.BoolOp:
		return v.EvalBoolOp(node, env)
	case *ast.Call:
		return v.EvalCall(node, env)
	case *ast.Identifier:
		return v.EvalIdentifier(node, env)
	case *ast.Literal:
		return v.EvalLiteral(node, env)
	case *ast.Subscript:
		return v.EvalSubscript(node, env)
	case *ast.Tuple:
		return v.EvalTuple(node, env)
	case *ast.UnaryOp:
		return v.EvalUnary(node, env)
	default:
		panic(stageleft.AttachExitCode(
			fmt.Errorf("unknown node type %T", node),
			stageleft.ExitCode(86),
		))
	}
}
