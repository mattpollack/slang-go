package ast

import (
	"errors"
	"fmt"
	"reflect"
)

/*

APPROACH:

Increase depth while tracking defined variables

Decrease depth while replacing patterns with applications, and populating an environment


*/

func Defun(ast AST, env *Environment) (AST, error) {
	return (&DefunVisitor{
		hasDef: env.HasDef(),
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
	switch A := a.(type) {
	case Application:
		return v.VisitApplication(A)
	case Pattern:
		return v.VisitPattern(A)
	case Identifier:
		return a, nil
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

func (v *DefunVisitor) VisitApplication(application Application) (AST, error) {
	for i := range application.Body {
		var err error
		application.Body[i], err = v.Copy().Visit(application.Body[i])

		if err != nil {
			return nil, err
		}
	}

	return application, nil
}

func (v *DefunVisitor) VisitLet(let Let) (AST, error) {
	// Set which identifiers are assigned to values, so we know which patterns need extra arguments
	for i, id := range let.BoundIds {
		if id.Value != "_" && !IsBuiltin(id.Value) {
			v.hasDef[id.Value] = true
		}

		var err error
		let.BoundValues[i], err = v.Copy().Visit(let.BoundValues[i])

		if err != nil {
			return nil, err
		}
	}

	var err error
	let.Body, err = v.Visit(let.Body)

	if err != nil {
		return nil, err
	}

	return let, nil
}

func (v *DefunVisitor) VisitPattern(pattern Pattern) (AST, error) {
	// Search for free identifiers, if any are hasDef[k] then this pattern needs extra parameters
	freeVars, err := FreeVars(pattern, v.hasDef)

	if err != nil {
		return nil, err
	}

	// TODO [x]: Create new pattern with free vars on each matchgroup
	// TODO [ ]: Replace pattern with application passing free vars
	newPattern := pattern.Copy().(Pattern)
	newVars := []AST{}

	// NOTE: Issue here, we need this to be a set order, currently ranging over a map
	for k := range freeVars {
		id, _ := NewIdentifier(k)

		newVars = append(newVars, id)
	}

	for i := range newPattern.Matches {
		newPattern.Matches[i] = append(append([]AST{}, newVars...), newPattern.Matches[i]...)
	}

	for i := range newPattern.Bodies {
		var err error
		newPattern.Bodies[i], err = v.Visit(newPattern.Bodies[i])

		if err != nil {
			return nil, err
		}
	}

	return newPattern, nil // nil, errors.New("TODO visit pattern")
}

type FreeVarsVisitor struct {
	freeVars     map[string]bool
	parentHasDef map[string]bool
	localHasDef  map[string]bool
}

func FreeVars(a AST, parentHasDef map[string]bool) (map[string]bool, error) {
	v := &FreeVarsVisitor{
		freeVars:     map[string]bool{},
		parentHasDef: parentHasDef,
		localHasDef:  map[string]bool{},
	}

	if err := v.Visit(a); err != nil {
		return nil, err
	}

	return v.freeVars, nil
}

func (v *FreeVarsVisitor) Visit(a AST) error {
	switch A := a.(type) {
	case Application:
		return v.VisitApplication(A)
	case Pattern:
		return v.VisitPattern(A)
	case Identifier:
		if A.Value == "_" || IsBuiltin(A.Value) {
			return nil
		}

		return v.VisitIdentifier(A)
	case Label:
		return nil
	case String:
		return nil
	case List:
	case ListConstructor:
	case Number:
		return nil
	case Let:
		return v.VisitLet(A)
	case Where:
	}

	return errors.New(fmt.Sprintf("Unhandled ast kind '%s' passed to FreeVarsVisitor", reflect.TypeOf(a).Name()))
}

func (v *FreeVarsVisitor) defScope(id Identifier) (bool, bool) {
	definedLocal, okLocal := v.localHasDef[id.Value]
	definedParent, okParent := v.parentHasDef[id.Value]

	return definedLocal && okLocal, definedParent && okParent
}

func (v *FreeVarsVisitor) VisitPattern(pattern Pattern) error {
	// Any free variables defined by a parent must be passed into all branches?
	for _, matchGroup := range pattern.Matches {
		for _, match := range matchGroup {
			if id, ok := match.(Identifier); ok {
				if err := v.VisitIdentifierPattern(id); err != nil {
					return nil
				}
			}

			if err := v.Visit(match); err != nil {
				return err
			}
		}
	}

	for _, body := range pattern.Bodies {
		if err := v.Visit(body); err != nil {
			return err
		}
	}

	return nil
}

// Identifier possibly defined in a pattern
func (v *FreeVarsVisitor) VisitIdentifierPattern(id Identifier) error {
	local, parent := v.defScope(id)

	if !local {
		if !parent {
			v.localHasDef[id.Value] = true
		} else {
			v.freeVars[id.Value] = true
		}
	}

	return nil
}

// Identifier defined in a let
func (v *FreeVarsVisitor) VisitIdentifierLet(id Identifier) error {
	v.localHasDef[id.Value] = true

	return nil
}

// All cases of access
func (v *FreeVarsVisitor) VisitIdentifier(id Identifier) error {
	local, parent := v.defScope(id)

	if local {
		return nil
	} else if parent {
		v.freeVars[id.Value] = true
		return nil
	}

	return errors.New(fmt.Sprintf("Identifier '%s' has no definition", id.Value))
}

func (v *FreeVarsVisitor) VisitApplication(app Application) error {
	for _, ast := range app.Body {
		if err := v.Visit(ast); err != nil {
			return err
		}
	}

	return nil
}

func (v *FreeVarsVisitor) VisitLet(let Let) error {
	for i, id := range let.BoundIds {
		if err := v.VisitIdentifierLet(id); err != nil {
			return err
		}

		if err := v.Visit(let.BoundValues[i]); err != nil {
			return err
		}
	}

	if err := v.Visit(let.Body); err != nil {
		return err
	}

	return nil
}
