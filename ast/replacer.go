package ast

import (
	"fmt"
	"errors"
	"reflect"
)

// Replaces AST with another AST

type Replacer struct {
	replaced AST
	replacement AST
}

func Replace(ast AST, replaced AST, replacement AST) (AST, error) {
	return (&Replacer{replaced, replacement}).Visit(ast)
}

func (v *Replacer) Visit(a AST) (AST, error) {
	if a.Equals(v.replaced) {
		return v.replacement.Copy(), nil
	}

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
		return v.VisitListConstructor(A)
	case Number:
		return a, nil
	case Let:
		return v.VisitLet(A)
	case Where:
	}

	return nil, errors.New(fmt.Sprintf("Unhandled ast kind '%s' passed to Replacer", reflect.TypeOf(a).Name()))
}

func (v *Replacer) VisitListConstructor(listCons ListConstructor) (AST, error) {
	var err error
	listCons.Head, err = v.Visit(listCons.Head)

	if err != nil {
		return nil, err
	}

	listCons.Tail, err = v.Visit(listCons.Tail)

	if err != nil {
		return nil, err
	}

	return listCons, nil
}

func (v *Replacer) VisitList(list List) (AST, error) {
	var err error

	for i := range list.Values {
		list.Values[i], err = v.Visit(list.Values[i])

		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (v *Replacer) VisitApplication(app Application) (AST, error) {
	var err error

	for i := range app.Body {
		app.Body[i], err = v.Visit(app.Body[i])

		if err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (v *Replacer) VisitLet(let Let) (AST, error) {
	var err error
	for i := range let.BoundValues {
		let.BoundValues[i], err = v.Visit(let.BoundValues[i])

		if err != nil {
			return nil, err
		}
	}

	let.Body, err = v.Visit(let.Body)

	if err != nil {
		return nil, err
	}

	return let, nil
}

func (v *Replacer) VisitPattern(pattern Pattern) (AST, error) {
	var err error

	for _, matchGroup := range pattern.Matches {
		for i := range matchGroup {
			matchGroup[i], err = v.Visit(matchGroup[i])

			if err != nil {
				return nil, err
			}
		}
	}

	for i := range pattern.Bodies {
		pattern.Bodies[i], err = v.Visit(pattern.Bodies[i])

		if err != nil {
			return nil, err
		}
	}

	return pattern, nil
}


