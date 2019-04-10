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

	prog, err := ast.Parse(
		`
let rec = {
  0 -> 1
  n -> (+ n (rec (- n 1)))
}

(fact 10)

`)

	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Print(prog)

	fmt.Println("\n---------------------------------------\n")

	compile.Compile(prog)
}
