package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"./ast"
	//"./compile"
	//"./vm"
)

var DEBUG = false

func main() {
	// Timer
	startTime := time.Now()
	defer func() {
		fmt.Println("\nExecution time:", time.Now().Sub(startTime))
	}()

	switch len(os.Args) {
	case 2:
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Runtime Error (%s)\n", err)
			}
		}()

		src, err := ioutil.ReadFile(os.Args[1])

		if err != nil {
			fmt.Println(err)
			return
		}

		AST, err := ast.Parse(src)

		if err != nil {
			fmt.Println(err)
			return
		}

		AST = ast.Interpret(AST)

		if DEBUG {
			fmt.Println("\n---------------------------------------\nAST:\n")
			AST.Print(0)
		}
	default:
		fmt.Println("Unexpected number of args")
	}
}

/*
func main() {
	// Timer
	startTime := time.Now()
	defer func() {
		fmt.Println("\nExecution time:", time.Now().Sub(startTime))
	}()

	src := []byte(
		`

let range = {
  min max : (> max min) v : (&& (>= v min) (< v max))
  -> .true
  => .false
}

let do = {
  s : (s.end) _  -> s
  s           fn -> (do (fn s) fn)
}

let _ = (do
  {
    .end -> .false
    .i   -> 10
  }
  {
    s : (> (s.i) 0)
    -> let _ = (print (s.i))
       { .end -> .false .i -> (- (s.i) 1) }
    => { .end -> .true }
  }
)

(print "test")
`)
	AST, err := ast.Parse(src)

	if err != nil {
		fmt.Println(err)
		return
	}

	AST = ast.Interpret(AST)
	fmt.Println("\n---------------------------------------\nAST:\n")
	AST.Print(0)

	if false {
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
}
*/
