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

		AST, err := ast.Parse(src)
		AST.Print(0)
	default:
		fmt.Println("Unexpected number of args")
	}
}
