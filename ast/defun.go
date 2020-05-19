package ast

import (
	"fmt"
	"errors"
	"reflect"
)

/*

APPROACH:

Increase depth while tracking defined variables

Decrease depth while replacing patterns with applications, and populating an environment


*/

func Defun(ast AST) (AST, error) {
	return (&DefunVisitor{
		hasDef: map[string]bool{},
	}).Visit(ast)
}

type DefunVisitor struct {
	hasDef map[string]bool
}

func (v *DefunVisitor) Copy() *DefunVisitor {
	hasDef := map[string]bool{}

	for k, v := range v.hasDef {
		hasDef[k] = v
	}

	return &DefunVisitor{
		hasDef: hasDef,
	}
}

func (v *DefunVisitor) Visit(a AST) (AST, error) {
	Print(a)
	fmt.Println("~~~")
	fmt.Println(v.hasDef)
	fmt.Println("-----------------")
	switch A := a.(type) {
	case Application:
	case Pattern:
		return v.VisitPattern(A)
	case Identifier:
	case Label:
		return a, nil
	case String:
		return a, nil
	case List:
	case ListConstructor:
	case Number:
		return a, nil
	case Let:
		return v.VisitLet(A)
	case Where:

	}

	return nil, errors.New(fmt.Sprintf("Unhandled ast kind '%s' passed to DefunVisitor", reflect.TypeOf(a).Name()))
}

func (v *DefunVisitor) VisitLet(let Let) (AST, error) {
	// Set which identifiers are assigned to values, so we know which patterns need extra arguments
	for i, id := range let.BoundIds {
		if id.Value != "_" {
			v.hasDef[id.Value] = true
		}

		var err error
		let.BoundValues[i], err = v.Copy().Visit(let.BoundValues[i])

		if err != nil {
			return nil, err
		}
	}

	return let, nil
}

func (v *DefunVisitor) VisitPattern(pattern Pattern) (AST, error) {
	// TODO: Search for free identifiers, if any have are hasDef[k] then this pattern needs extra parameters
	freeVars, err := GetFreeVars(pattern, v.hasDef)

	if err != nil {
		return nil, err
	}

	fmt.Println("*****")
	fmt.Println(freeVars)

	return nil, errors.New("TODO visit pattern")
}

type FreeVarsVisitor struct {
	freeVars map[string]bool
	parentHasDef map[string]bool
	localHasDef map[string]bool
}

func GetFreeVars(a AST, parentHasDef map[string]bool) (map[string]bool, error) {
	v := &FreeVarsVisitor{
		freeVars: map[string]bool{},
		parentHasDef: parentHasDef,
		localHasDef: map[string]bool{},
	}

	if err := v.Visit(a); err != nil {
		return nil, err
	}

	return v.freeVars, nil
}

func (v *FreeVarsVisitor) Visit(a AST) error {
	switch A := a.(type) {
	case Application:
	case Pattern:
		// If a pattern is trying to define an identifier that already has a parent definition, then its a free variable
		// If an identifier is neither a parent or local definition, then its an error
		return v.VisitPattern(A)
	case Identifier:
		return v.VisitIdentifier(A)
	case Label:
	case String:
	case List:
	case ListConstructor:
	case Number:
	case Let:
	case Where:
	}

	return errors.New(fmt.Sprintf("Unhandled ast kind '%s' passed to FreeVarsVisitor", reflect.TypeOf(a).Name()))
}

func (v *FreeVarsVisitor) VisitPattern(pattern Pattern) error {
	// Any free variables defined by a parent must be passed into all branches?
	for _, matchGroup := range pattern.Matches {
		for _, match := range matchGroup {
			v.Visit(match)
		}
	}

	fmt.Println(v)

	return errors.New("TODO free vars visit pattern")
}

// Identifier possibly defined in a pattern
func (v *FreeVarsVisitor) VisitIdentifierPattern(id Identifier) error {
	return errors.New("TODO free vars visit id pattern")
}

// Identifier defined in a let
func (v *FreeVarsVisitor) VisitIdentifierLet(id Identifier) error {
	return errors.New("TODO free vars visit id let")
}

// All cases of access
func (v *FreeVarsVisitor) VisitIdentifier(id Identifier) error {
	definedLocal, okLocal := v.localHasDef[id.Value]
	definedLocal = definedLocal && okLocal
	definedParent, okParent := v.parentHasDef[id.Value]
	definedParent = definedParent && okParent

	if !definedLocal && definedParent {
		v.freeVars[id.Value] = true
		return nil
	}

	return errors.New(fmt.Sprintf("Identifier '%s' has no definition", id.Value))
}
