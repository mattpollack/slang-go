package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"./ast"
)

func main() {
	// Timer
	startTime := time.Now()
	defer func() {
		fmt.Println("\nExecution time:", time.Now().Sub(startTime))
	}()

	switch len(os.Args) {
	case 2:
		/*
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Runtime Error (%s)\n", err)
				}
			}()*/

		src, err := ioutil.ReadFile(os.Args[1])

		if err != nil {
			fmt.Println(err)
			return
		}

		prog, err := ast.Parse(src)

		if err != nil {
			fmt.Println(err)
			return
		}

		ast.Print(prog)

		fmt.Println("--- RESULT ---------------------------------")

		prog, err = prog.Eval(ast.StdLib)

		fmt.Println("--------------------------------------------")

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	default:
		fmt.Println("Unexpected number of args")
	}
}
