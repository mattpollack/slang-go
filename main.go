package main

import (
	"./ast"
	"./compile"
	"./vm"

	"fmt"
	"time"
)

func main() {
	// Timer
	startTime := time.Now()
	defer func() {
		fmt.Println("\nExecution time:", time.Now().Sub(startTime))
	}()

	src := []byte(
		`
let test = {
  0 -> 5
  1 -> 6
  2 -> 7
  3 -> 8
  4 -> 9
  5 -> 15
}

(test 5)
`)
	AST, err := ast.Parse(src)

	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Print(AST)

	fmt.Println("\n---------------------------------------\n")

	prog := compile.Compile(AST)
	prog.Print()

	fmt.Printf("\nSummary\n- Bytecode:\t %d bytes\n- Source:\t %d bytes\n", prog.Size(), len(src))

	fmt.Println("\n---------------------------------------\n")

	run := vm.NewVM(prog.Render())

	for run.Status == vm.VM_RUNNING {
		//run.PrintStack()
		//fmt.Println()

		run.Step()
	}

	if run.Err != nil {
		fmt.Println(run.Err)
	} else {
		fmt.Println("# RESULT:")
		run.PrintStack()
	}
}
