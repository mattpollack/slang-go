package ast

import (
	"fmt"
	"errors"
)

/*

APPROACH:

Increase depth while tracking defined variables

Decrease depth while replacing patterns with applications, and populating an environment


*/

func Defun(ast AST) (AST, error) {
	return (&DefunVisitor{}).Visit(ast)
}

type DefunVisitor struct {
	hasDef map[Identifier]bool
}

func (v *DefunVisitor) Copy() *DefunVisitor {
	hasDef := map[Identifier]bool{}

	for k, v := range v.hasDef {
		hasDef[k] = v
	}

	return &DefunVisitor{
		hasDef: hasDef,
	}
}

func (v *DefunVisitor) Visit(a AST) (AST, error) {
	Print(a)
	fmt.Println("-----------------")
	switch A := a.(type) {
	case Let:
		return v.VisitLet(A)
	}

	return nil, errors.New("Unhandled ast kind passed to DefunVisitor")
}

func (v *DefunVisitor) VisitLet(let Let) (AST, error) {
	fmt.Println(let)

	// Set which identifiers are assigned to values, so we know which patterns need extra arguments
	for i, id := range let.BoundIds {
		v.hasDef[id] = true

		// recur
	}

	return let, nil
}
