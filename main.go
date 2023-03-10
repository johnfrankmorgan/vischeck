package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(&analysis.Analyzer{
		Name: "vischeck",
		Doc:  "flurb",
		Run:  run,
	})
}

type visibility string

func (vis visibility) valid() bool {
	switch vis {
	case visReadonly:
		return true
	}

	return false
}

const (
	visTag      = "visibility"
	visReadonly = "readonly"
)

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.Field:
				if node.Tag == nil {
					break
				}

				tag, _ := strconv.Unquote(node.Tag.Value)
				vis := reflect.StructTag(tag).Get(visTag)

				if !visibility(vis).valid() {
					pass.Reportf(node.Tag.Pos(), "invalid %s tag: %q", visTag, vis)
				}

			case *ast.AssignStmt:
				for _, lhs := range node.Lhs {
					lhs, ok := lhs.(*ast.SelectorExpr)
					if !ok {
						continue
					}

					check(pass, lhs, "cannot assign")
				}

			case *ast.IncDecStmt:
				typ := "decrement"

				if node.Tok == token.INC {
					typ = "increment"
				}

				check(pass, node.X, "cannot "+typ)

			case *ast.UnaryExpr:
				if _, ok := pass.TypesInfo.TypeOf(node).(*types.Pointer); ok {
					// can't take address of readonly fields, as
					// pointers imply mutation
					check(pass, node.X, "cannot take address")
				}
			}

			return true
		})
	}

	return nil, nil
}

func check(pass *analysis.Pass, node ast.Expr, message string) {
	expr, ok := node.(*ast.SelectorExpr)
	if !ok {
		return
	}

	sel := pass.TypesInfo.Selections[expr]

	var typ *types.Named

	switch val := sel.Recv().(type) {
	case *types.Named:
		typ = val

	case *types.Pointer:
		// dereference type until it's not a pointer
		for typ == nil {
			switch t := val.Elem().(type) {
			case *types.Pointer:
				val = t

			case *types.Named:
				typ = t

			default:
				return
			}
		}

	default:
		return
	}

	for i := 0; i < typ.NumMethods(); i++ {
		if typ.Method(i).Scope().Contains(node.Pos()) {
			// in receiver method, assume mutation is safe
			return
		}
	}

	str, ok := typ.Underlying().(*types.Struct)
	if !ok {
		return
	}

	for i := 0; i < str.NumFields(); i++ {
		if str.Field(i).Name() != sel.Obj().Name() {
			continue
		}

		tag, ok := reflect.StructTag(str.Tag(i)).Lookup(visTag)
		if !ok {
			continue
		}

		switch tag {
		case visReadonly:
			pass.Reportf(expr.Pos(), "misuse of %s field: %s", tag, message)

		default:
			pass.Reportf(str.Field(i).Pos(), "invalid %s tag: %q", visTag, tag)
		}
	}

	return
}
