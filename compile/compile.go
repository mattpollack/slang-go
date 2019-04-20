package compile

import (
	"../ast"
	"../vm"

	"fmt"
)

var builtins map[string]vm.Instruction = map[string]vm.Instruction{
	"+": &vm.OP{vm.INS_ADD},
	"-": &vm.OP{vm.INS_SUB},
}

type Meta struct {
	name string
}

func Compile(AST ast.Expression) *vm.Program {
	AST = pass0(AST)

	main := &vm.BasicBlock{Label: "__main_entry"}
	prog := &vm.Program{}
	prog.Push(main)

	visitExpr(prog, main, AST)

	// Remove the return from main
	switch ins := main.Body[len(main.Body)-1].(type) {
	case *vm.OP:
		if ins.Kind == vm.INS_RETURN {
			main.Body = main.Body[:len(main.Body)-1]
		}
	}

	main.Push(&vm.OP{vm.INS_EXIT})

	return prog
}

func visitExpr(prog *vm.Program, block *vm.BasicBlock, _e ast.Expression) {
	if _e == nil {
		return
	}

	switch node := _e.(type) {
	case ast.Application:
		for i := len(node.Body) - 1; i >= 0; i-- {
			visitExpr(prog, block, node.Body[i])
		}

	case ast.Pattern:
		block.Push(&vm.OP{vm.INS_CALL})
		block.Push(&vm.Address{fmt.Sprintf("%s_0", block.Label)})
		block.Push(&vm.OP{vm.INS_RETURN})

		// Push a basic block for each body
		for i, body := range node.Bodies {
			next := &vm.BasicBlock{Label: fmt.Sprintf("%s_%d", block.Label, i)}
			prog.Push(next)

			for _, m := range node.Matches[i] {
				next.Push(&vm.OP{vm.INS_DUP})
				visitMatch(prog, next, m)
				next.Push(&vm.OP{vm.INS_JNE})
				next.Push(&vm.Address{fmt.Sprintf("%s_%d", block.Label, i+1)})
				next.Push(&vm.OP{vm.INS_POP})
			}

			visitExpr(prog, next, body)
			next.Push(&vm.OP{vm.INS_RETURN})
		}

		// Push the error handling block
		errorHandler := &vm.BasicBlock{Label: fmt.Sprintf("%s_%d", block.Label, len(node.Bodies))}
		prog.Push(errorHandler)
		errorHandler.Push(&vm.TODO{"error handling"})
		errorHandler.Push(&vm.OP{vm.INS_EXIT})

	case ast.Identifier:
		if v, ok := builtins[node.Value]; ok {
			block.Push(v)
		} else {
			block.Push(&vm.OP{vm.INS_CALL})
			block.Push(&vm.Address{fmt.Sprintf("%s", node.Value)})
		}

	case ast.Label:
		block.Push(&vm.TODO{"Label"})

	case ast.String:
		block.Push(&vm.TODO{"String"})

	case ast.Number:
		block.Push(&vm.OP{vm.INS_PUSH})
		block.Push(vm.NewInt32(int32(node.Value)))

	case ast.Let:
		for i, id := range node.BoundIds {
			next := &vm.BasicBlock{Label: id.Value}
			visitExpr(prog, next, node.BoundValues[i])
			prog.Push(next)

			// If the last operation in a let binding isn't a return, add a return
			switch ins := next.Body[len(next.Body)-1].(type) {
			case *vm.OP:
				if ins.Kind != vm.INS_RETURN {
					next.Push(&vm.OP{vm.INS_RETURN})
				}
			default:
				next.Push(&vm.OP{vm.INS_RETURN})
			}
		}

		visitExpr(prog, block, node.Body)

	case ast.If:
		block.Push(&vm.TODO{"If"})

	case BoundValue:
		block.Push(&vm.TODO{"BoundValue"})

	default:
		panic("Unexpected type")
	}
}

func visitMatch(prog *vm.Program, block *vm.BasicBlock, m ast.Match) {
	if m == nil {
		return
	}

	switch node := m.(type) {
	case ast.Identifier:
		block.Push(&vm.TODO{"Match Identifier"})

	case ast.Label:
		block.Push(&vm.TODO{"Match Label"})

	case ast.String:
		block.Push(&vm.TODO{"Match String"})

	case ast.Number:
		block.Push(&vm.OP{vm.INS_PUSH})
		block.Push(vm.NewInt32(int32(node.Value)))

	case ast.Where:
		block.Push(&vm.TODO{"Match Where"})

	default:
		panic(node)
	}
}
