package compile

import (
	"../ast"
	"../vm"
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

	prog := &vm.Program{}
	block := &vm.BasicBlock{Label: "__main_entry"}

	prog.Push(block)

	visitExpr(prog, block, AST)

	return prog
}

func visitExpr(prog *vm.Program, block *vm.BasicBlock, e ast.Expression) {
	if e == nil {
		return
	}

	switch e.(type) {
	case ast.Application:
		block.Push(&vm.TODO{"Application"})

	case ast.Pattern:
		block.Push(&vm.TODO{"Pattern"})

	case ast.Identifier:
		block.Push(&vm.TODO{"Identifier"})

	case ast.Label:
		block.Push(&vm.TODO{"Label"})

	case ast.String:
		block.Push(&vm.TODO{"String"})

	case ast.Number:
		block.Push(&vm.TODO{"Number"})

	case ast.Let:
		block.Push(&vm.TODO{"Let"})

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
		block.Push(&vm.TODO{"Match Number"})

	case ast.Where:
		block.Push(&vm.TODO{"Match Where"})

	default:
		panic(node)
	}
}
