package bytecode

import (
	"../ast"

	"fmt"
)

var builtins map[string]func() = map[string]func(){
	"+": func() {
		fmt.Println("  add")
	},
	"-": func() {
		fmt.Println("  sub")
	},
	"eq": func() {
		fmt.Println("  eq")
	},
}

type Meta struct {
	name string
}

func Compile(prog ast.Expression) {
	prog = pass0(prog)

	fmt.Println("__main_entry:")
	visitExpr(prog)
}

func visitMatch(m ast.Match) {
	if m == nil {
		return
	}

	switch node := m.(type) {
	case ast.Identifier:
		fmt.Println("  TODO identifier")
	case ast.Label:
		fmt.Println("  TODO label")
	case ast.String:
		fmt.Println("  TODO string")
	case ast.Number:
		fmt.Printf("  push %d\n", node.Value)
	case ast.Where:
		fmt.Printf("  arg_set %d\n", node.Id.(BoundValue).Index)
		//fmt.Println("  TODO: if true set arg")
		visitExpr(node.Condition)
	default:
		panic(node)
	}
}

func visitExpr(e ast.Expression) {
	if e == nil {
		return
	}

	switch node := e.(type) {
	case ast.Application:
		for i := len(node.Body) - 1; i >= 0; i-- {
			visitExpr(node.Body[i])
		}

	case ast.Pattern:
		for i, group := range node.Matches {
			for _, match := range group {
				if i != 0 {
					fmt.Printf("%s_%d:\n", node.MetaGet().(Meta).name, i)
				}

				visitMatch(match)

				if i < len(node.Matches)-1 {
					fmt.Printf("  jne %s_%d\n", node.MetaGet().(Meta).name, i+1)
				} else {
					fmt.Printf("  jne __RUNTIME_ERROR\n")
				}

				visitExpr(node.Bodies[i])
				fmt.Printf("  return\n")

				if i+1 != len(node.Matches) {
					fmt.Println()
				}
			}
		}

	case ast.Identifier:
		if fn, ok := builtins[node.Value]; ok {
			fn()
		} else {
			fmt.Printf("  call %s\n", node.Value)
		}

	case ast.Label:
		fmt.Println("  TODO label")

	case ast.String:
		fmt.Println("  TODO string")

	case ast.Number:
		fmt.Printf("  push %d\n", node.Value)

	case ast.Let:
		visitExpr(node.Body)
		fmt.Println()

		for i, e := range node.BoundValues {
			meta := e.MetaGet().(Meta)
			meta.name = node.BoundIds[i].Value
			e = e.MetaSet(meta).(ast.Expression) // Value doesn't need to be set???

			fmt.Printf("%s:\n", meta.name)
			visitExpr(e)
		}

	case ast.If:
		fmt.Println("  TODO if")

	case BoundValue:
		fmt.Printf("  arg_get %d\n", node.Index)

	default:
		panic("TODO compile node if this type")
	}
}
