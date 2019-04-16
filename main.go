package main

import (
	"./ast"
	"./compile"

	"fmt"
	"time"
)

func main() {
	// Timer
	startTime := time.Now()
	defer func() {
		fmt.Println("\nExecution time:", time.Now().Sub(startTime))
	}()

	AST, err := ast.Parse([]byte(
		`
(1
 2
 3
 4)
`))

	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Print(AST)

	fmt.Println("\n---------------------------------------\n")

	prog := compile.Compile(AST)
	prog.Print()

	/*
		fmt.Println("\n---------------------------------------\n")
		block := vm.BasicBlock{Label: "test"}
		block.Push(&vm.OP{vm.INS_ADD})
		block.Push(vm.NewInt32(10))
		block.Push(vm.NewInt32(10))
		block.Push(&vm.OP{vm.INS_SUB})
		block.Push(&vm.OP{vm.INS_MUL})
		block.Push(&vm.OP{vm.INS_DIV})
		block.Print()
	*/
}
