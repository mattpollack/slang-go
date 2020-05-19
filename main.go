package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"./ast"
)

func loadFile(path string) (*ast.SourceFile, error) {
	// TODO: fix relative pathing
	src, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	srcFile, err := ast.Parse(src)

	if err != nil {
		return nil, err
	}

	env := ast.NewEnv(ast.StdLib)

	for _, imp := range srcFile.Imports {
		impSrcFile, err := loadFile(imp.Path)

		// ERRORING HERE???????????
		if err != nil {
			return nil, err
		}

		name := imp.Name

		if name == "" {
			name = impSrcFile.PackageName
		}

		env.Set(name, impSrcFile.Definition)
	}

	srcFile.Definition, err = srcFile.Eval(env)

	if err != nil {
		return nil, err
	}

	return srcFile, nil
}

func main() {
	// Timer
	startTime := time.Now()
	defer func() {
		fmt.Println("\nExecution time:", time.Now().Sub(startTime))
	}()

	switch len(os.Args) {
	case 2:
		// Read queue, loaded with first file
		_, err := loadFile(os.Args[1])

		if err != nil {
			fmt.Println(err)
			return
		}

	default:
		fmt.Println("Unexpected number of args")
	}
}
