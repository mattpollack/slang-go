package compile

import (
	"../ast"
	"../vm"

	"fmt"
)

func Compile(AST ast.Expression) *vm.Program {
	return nil
}

func visitExpr(prog *vm.Program, block *vm.BasicBlock, _e ast.Expression) {
	panic("todo visit match")
}

func visitMatch(prog *vm.Program, block *vm.BasicBlock, m ast.Match, addr string) {
	panic("todo visit match")
}
