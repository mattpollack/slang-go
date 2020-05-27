package ast

import (
	"errors"
	"fmt"
	"reflect"
)

func Defun(in AST, env *Environment) (AST, error) {
	if _, ok := in.(Let); !ok {
		in, _ = NewLet([]Identifier{}, []AST{}, in)
	}

	// Make sure the input is a let, for non-branching logic
	inLet := in.(Let)
	resLet, _ := NewLet([]Identifier{}, []AST{}, nil)
	nextUniqueId := 0
	hasDef := env.HasDef()

	// Visit inLet bound
	for i := range inLet.BoundIds {
		//hasDef[inLet.BoundIds[i].Value] = true

		v := &DefunVisitor{
			hasDef:       hasDef,
			nextUniqueId: &nextUniqueId,
			bound:        []DefunBound{},
		}

		body, err := v.Visit(inLet.BoundValues[i])

		if err != nil {
			return nil, err
		}

		for _, bound := range v.bound {
			resLet.Bind(
				bound.id,
				bound.body,
			)
		}

		resLet.Bind(inLet.BoundIds[i], body)
	}

	v := &DefunVisitor{
		hasDef:       hasDef,
		nextUniqueId: &nextUniqueId,
		bound:        []DefunBound{},
	}

	var err error
	resLet.Body, err = v.Visit(inLet.Body)

	if err != nil {
		return nil, err
	}

	for _, bound := range v.bound {
		resLet.Bind(
			bound.id,
			bound.body,
		)
	}

	return resLet, nil
}

type DefunVisitor struct {
	hasDef       map[string]bool
	nextUniqueId *int
	bound        []DefunBound
}

type DefunBound struct {
	id   Identifier
	body AST
}

func (v *DefunVisitor) NextUniqueId() Identifier {
	*v.nextUniqueId = *v.nextUniqueId + 1

	return Identifier{
		Value: fmt.Sprintf("'%d", *v.nextUniqueId),
	}
}

func (v *DefunVisitor) Copy() *DefunVisitor {
	hasDef := map[string]bool{}

	for k, v := range v.hasDef {
		hasDef[k] = v
	}

	return &DefunVisitor{
		hasDef:       hasDef,
		nextUniqueId: v.nextUniqueId,
		bound:        v.bound,
	}
}

func (v *DefunVisitor) Bind(id Identifier, body AST) {
	v.bound = append(v.bound, DefunBound{id, body})
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
		return v.VisitList(A)
	case ListConstructor:
	case Number:
		return a, nil
	case Let:
		return v.VisitLet(A)
	case Where:
	}

	return nil, errors.New(fmt.Sprintf("Unhandled ast kind '%s' passed to DefunVisitor", reflect.TypeOf(a).Name()))
}

func (v *DefunVisitor) VisitList(list List) (AST, error) {
	for i := range list.Values {
		var err error
		list.Values[i], err = v.Copy().Visit(list.Values[i])

		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (v *DefunVisitor) VisitApplication(application Application) (AST, error) {
	for i := range application.Body {
		var err error
		application.Body[i], err = v.Visit(application.Body[i])

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
		let.BoundValues[i], err = v.Visit(let.BoundValues[i])

		if err != nil {
			return nil, err
		}

		// TODO: [ ]
		// Recursive calls are generating bad applications like:
		//   g = ( '1 str g )
		// If an application is recursive, it needs to pass a reference of itself, including all parent bound values
		//   g = ( '1 str )
		//   h = g g
	}

	var err error
	let.Body, err = v.Visit(let.Body)

	if err != nil {
		return nil, err
	}

	return let, nil
}

func (v *DefunVisitor) VisitPattern(pattern Pattern) (AST, error) {
	freeVars, err := FreeVars(pattern, v.hasDef)

	if err != nil {
		return nil, err
	}

	boundIds, err := MatchesThatBind(pattern, v.hasDef)

	for k := range boundIds {
		v.hasDef[k] = true
	}

	newPattern := pattern.Copy().(Pattern)

	for i := range newPattern.Matches {
		newPattern.Matches[i] = append(append([]AST{}, freeVars...), newPattern.Matches[i]...)
	}

	for i := range newPattern.Bodies {
		var err error
		newPattern.Bodies[i], err = v.Visit(newPattern.Bodies[i])

		if err != nil {
			return nil, err
		}
	}

	newId := v.NextUniqueId()

	v.Bind(newId, newPattern)

	// Build application to replace pattern
	app, _ := NewApplication([]AST{newId})

	for _, id := range freeVars {
		app.Body = append(app.Body, id)
	}

	return app, nil
}

type FreeVarsVisitor struct {
	freeVarsMap  map[string]bool
	freeVarsList []AST
	parentHasDef map[string]bool
	localHasDef  map[string]bool
}

func FreeVars(a AST, parentHasDef map[string]bool) ([]AST, error) {
	v := &FreeVarsVisitor{
		freeVarsMap:  map[string]bool{},
		freeVarsList: []AST{},
		parentHasDef: parentHasDef,
		localHasDef:  map[string]bool{},
	}

	if err := v.Visit(a); err != nil {
		return nil, err
	}

	return v.freeVarsList, nil
}

func MatchesThatBind(p Pattern, parentHasDef map[string]bool) (map[string]bool, error) {
	v := &FreeVarsVisitor{
		freeVarsMap:  map[string]bool{},
		freeVarsList: []AST{},
		parentHasDef: parentHasDef,
		localHasDef:  map[string]bool{},
	}

	if err := v.VisitPattern(p, true); err != nil {
		return nil, err
	}

	return v.localHasDef, nil
}

func (v *FreeVarsVisitor) markFree(id string) {
	if _, ok := v.freeVarsMap[id]; !ok {
		v.freeVarsMap[id] = true
		v.freeVarsList = append(v.freeVarsList, Identifier{id})
	}
}

func (v *FreeVarsVisitor) Visit(a AST) error {
	switch A := a.(type) {
	case Application:
		return v.VisitApplication(A)
	case Pattern:
		return v.VisitPattern(A, false)
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
		return v.VisitList(A)
	case ListConstructor:
	case Number:
		return nil
	case Let:
		return v.VisitLet(A)
	case Where:
	}

	return errors.New(fmt.Sprintf("Unhandled ast kind '%s' passed to FreeVarsVisitor", reflect.TypeOf(a).Name()))
}

func (v *FreeVarsVisitor) VisitMatch(a AST, onlyMatches bool) error {
	switch A := a.(type) {
	case Application:
	case Pattern:
	case Identifier:
		if A.Value == "_" || IsBuiltin(A.Value) {
			return nil
		}

		return v.VisitIdentifierMatch(A)
	case Label:
		return nil
	case String:
		return nil
	case List:
		return v.VisitListMatch(A, onlyMatches)
	case ListConstructor:
		return v.VisitListConstructorMatch(A, onlyMatches)
	case Number:
		return nil
	case Let:
	case Where:
		return v.VisitWhereMatch(A, onlyMatches)
	}

	return errors.New(fmt.Sprintf("Unhandled ast kind '%s' passed to FreeVarsVisitor Match", reflect.TypeOf(a).Name()))
}

func (v *FreeVarsVisitor) defScope(id Identifier) (bool, bool) {
	definedLocal, okLocal := v.localHasDef[id.Value]
	definedParent, okParent := v.parentHasDef[id.Value]

	return definedLocal && okLocal, definedParent && okParent
}

func (v *FreeVarsVisitor) VisitIdentifierMatch(id Identifier) error {
	local, parent := v.defScope(id)

	if !local {
		if !parent {
			v.localHasDef[id.Value] = true
		} else {
			v.markFree(id.Value)
		}
	}

	return nil
}

func (v *FreeVarsVisitor) VisitWhereMatch(where Where, onlyMatches bool) error {
	if err := v.VisitMatch(where.Match, onlyMatches); err != nil {
		return err
	}

	if !onlyMatches {
		if err := v.Visit(where.Condition); err != nil {
			return err
		}
	}

	return nil
}

func (v *FreeVarsVisitor) VisitListMatch(list List, onlyMatches bool) error {
	for _, ast := range list.Values {
		if err := v.VisitMatch(ast, onlyMatches); err != nil {
			return err
		}
	}

	return nil
}

func (v *FreeVarsVisitor) VisitListConstructorMatch(listCons ListConstructor, onlyMatches bool) error {
	if err := v.VisitMatch(listCons.Head, onlyMatches); err != nil {
		return err
	}

	if err := v.VisitMatch(listCons.Tail, onlyMatches); err != nil {
		return err
	}

	return nil
}

func (v *FreeVarsVisitor) VisitList(list List) error {
	for _, ast := range list.Values {
		if err := v.Visit(ast); err != nil {
			return err
		}
	}

	return nil
}

func (v *FreeVarsVisitor) VisitPattern(pattern Pattern, onlyMatches bool) error {
	// Any free variables defined by a parent must be passed into all branches?
	for _, matchGroup := range pattern.Matches {
		for _, match := range matchGroup {
			if err := v.VisitMatch(match, onlyMatches); err != nil {
				return err
			}
		}
	}

	if !onlyMatches {
		for _, body := range pattern.Bodies {
			if err := v.Visit(body); err != nil {
				return err
			}
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
		v.markFree(id.Value)
		return nil
	}

	// NOTE: there may be an issue with this check?
	// return nil
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
