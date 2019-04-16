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

	src := []byte(
		`
let test = {
  0 -> 1
}

(test 0)
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

	fmt.Printf("\nSummary\n- Bytecode:\t %d bytes\n- Source:\t %d bytes\n", prog.Size()/8, len(src))
}
