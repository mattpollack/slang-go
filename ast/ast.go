package ast

import (
	"fmt"
)

// Provides access to a meta object
type Meta interface {
	MetaSet(interface{}) interface{}
	MetaGet() interface{}
}

type Prints interface {
	Print(int)
}

type Expression interface {
	Meta
	Prints
	EqualsExpr(Expression) bool

	IsExpression()
}

type Match interface {
	Meta
	Prints
	EqualsMatch(Match) bool

	IsMatch()
}

// --------------------------------------------------------

type Application struct {
	Body []Expression

	meta interface{}
}

func (e Application) IsExpression() {}

func NewApplication(body []Expression) (Application, error) {
	return Application{
		body,
		nil,
	}, nil
}

func (A Application) EqualsExpr(b Expression) bool {
	switch B := b.(type) {
	case Application:
		if len(A.Body) == len(B.Body) {
			for i := 0; i < len(A.Body); i++ {
				if !A.Body[i].EqualsExpr(B.Body[i]) {
					return false
				}
			}

			return true
		}
	}

	return false
}

func (e Application) MetaGet() interface{} {
	return e.meta
}

func (e Application) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (E Application) Print(tab int) {
	printTab(tab)
	fmt.Println("(")

	for _, e := range E.Body {
		e.Print(tab + 1)
		fmt.Println()
	}

	printTab(tab)
	fmt.Print(")")
}

// --------------------------------------------------------

type If struct {
	Condition Expression
	Tbody     Expression
	Fbody     Expression

	meta interface{}
}

func (e If) IsExpression() {}

func NewIf(c, t, f Expression) (If, error) {
	return If{
		c,
		t,
		f,
		nil,
	}, nil
}

func (A If) EqualsExpr(b Expression) bool {
	switch B := b.(type) {
	case If:
		return A.Condition.EqualsExpr(B.Condition) && A.Tbody.EqualsExpr(B.Tbody) && A.Fbody.EqualsExpr(B.Fbody)
	}

	return false
}

func (e If) MetaGet() interface{} {
	return e.meta
}

func (e If) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (node If) Print(tab int) {
	printTab(tab)
	fmt.Println("if ")
	node.Condition.Print(tab + 1)
	fmt.Println()
	node.Tbody.Print(tab + 1)
	fmt.Println()
	printTab(tab)
	fmt.Println("else ")
	node.Fbody.Print(tab + 1)
}

// --------------------------------------------------------

type Pattern struct {
	Matches [][]Match
	Bodies  []Expression

	meta interface{}
}

func (e Pattern) IsExpression() {}

func NewPattern(m [][]Match, b []Expression) (Pattern, error) {
	return Pattern{
		m,
		b,
		nil,
	}, nil
}

func (A Pattern) EqualsExpr(b Expression) bool {
	switch B := b.(type) {
	case Pattern:
		if len(A.Matches) == len(B.Matches) && (len(A.Matches) == 0 || len(A.Matches[0]) == len(B.Matches[0])) && len(A.Bodies) == len(B.Bodies) {
			for i := 0; i < len(A.Matches); i++ {
				for j := 0; j < len(A.Matches[i]); j++ {
					if !A.Matches[i][j].EqualsMatch(B.Matches[i][j]) {
						return false
					}
				}
			}

			for i := 0; i < len(A.Bodies); i++ {
				if !A.Bodies[i].EqualsExpr(B.Bodies[i]) {
					return false
				}
			}

			return true
		}
	}

	return false
}

func (e Pattern) MetaGet() interface{} {
	return e.meta
}

func (e Pattern) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (node Pattern) Print(tab int) {
	printTab(tab)
	fmt.Println("{")

	for i := 0; i < len(node.Bodies); i++ {
		for _, m := range node.Matches[i] {
			m.Print(tab)
		}

		fmt.Println("-> ")
		node.Bodies[i].Print(tab + 1)
		fmt.Println()
	}
	printTab(tab)
	fmt.Print("}")
}

// --------------------------------------------------------

type Identifier struct {
	Value string

	meta interface{}
}

func (e Identifier) IsExpression() {}
func (e Identifier) IsMatch()      {}

func NewIdentifier(v string) (Identifier, error) {
	return Identifier{
		v,
		nil,
	}, nil
}

func (A Identifier) EqualsExpr(b Expression) bool {
	switch B := b.(type) {
	case Identifier:
		return A.Value == B.Value
	}

	return false
}

func (A Identifier) EqualsMatch(b Match) bool {
	switch B := b.(type) {
	case Identifier:
		return A.Value == B.Value
	}

	return false
}

func (e Identifier) MetaGet() interface{} {
	return e.meta
}

func (e Identifier) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (node Identifier) Print(tab int) {
	printTab(tab)
	fmt.Println(node.Value)
}

// --------------------------------------------------------

type Label struct {
	Value string

	meta interface{}
}

func (e Label) IsExpression() {}
func (e Label) IsMatch()      {}

func NewLabel(v string) (Label, error) {
	return Label{
		v,
		nil,
	}, nil
}

func (A Label) EqualsExpr(b Expression) bool {
	switch B := b.(type) {
	case Label:
		return A.Value == B.Value
	}

	return false
}

func (A Label) EqualsMatch(b Match) bool {
	switch B := b.(type) {
	case Label:
		return A.Value == B.Value
	}

	return false
}

func (e Label) MetaGet() interface{} {
	return e.meta
}

func (e Label) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (node Label) Print(tab int) {
	printTab(tab)
	fmt.Print(node.Value)
}

// --------------------------------------------------------

type String struct {
	Value string

	meta interface{}
}

func (e String) IsExpression() {}
func (e String) IsMatch()      {}

func NewString(v string) (String, error) {
	return String{
		v,
		nil,
	}, nil
}

func (A String) EqualsExpr(b Expression) bool {
	switch B := b.(type) {
	case String:
		return A.Value == B.Value
	}

	return false
}

func (A String) EqualsMatch(b Match) bool {
	switch B := b.(type) {
	case String:
		return A.Value == B.Value
	}

	return false
}

func (e String) MetaGet() interface{} {
	return e.meta
}

func (e String) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (node String) Print(tab int) {
	printTab(tab)
	fmt.Println(node.Value)
}

// --------------------------------------------------------

type Number struct {
	Value int

	meta interface{}
}

func (e Number) IsExpression() {}
func (e Number) IsMatch()      {}

func NewNumber(v int) (Number, error) {
	return Number{
		v,
		nil,
	}, nil
}

func (A Number) EqualsExpr(b Expression) bool {
	switch B := b.(type) {
	case Number:
		return A.Value == B.Value
	}

	return false
}

func (A Number) EqualsMatch(b Match) bool {
	switch B := b.(type) {
	case Number:
		return A.Value == B.Value
	}

	return false
}

func (e Number) MetaGet() interface{} {
	return e.meta
}

func (e Number) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (node Number) Print(tab int) {
	printTab(tab)
	fmt.Println(node.Value)
}

// --------------------------------------------------------

type Let struct {
	BoundIds    []Identifier
	BoundValues []Expression
	Body        Expression

	meta interface{}
}

func (e Let) IsExpression() {}

func NewLet(ids []Identifier, vs []Expression, b Expression) (Let, error) {
	return Let{
		ids,
		vs,
		b,
		nil,
	}, nil
}

func (A Let) EqualsExpr(b Expression) bool {
	switch B := b.(type) {
	case Let:
		if len(A.BoundIds) == len(B.BoundIds) {
			for i := 0; i < len(A.BoundIds); i++ {
				if !A.BoundIds[i].EqualsExpr(B.BoundIds[i]) || !A.BoundValues[i].EqualsExpr(B.BoundValues[i]) {
					return false
				}
			}

			return A.Body.EqualsExpr(B.Body)
		}
	}

	return false
}

func (e Let) MetaGet() interface{} {
	return e.meta
}

func (e Let) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (node Let) Print(tab int) {
	for i, id := range node.BoundIds {
		printTab(tab)
		fmt.Print("let ")
		id.Print(0)
		fmt.Println(" =")
		node.BoundValues[i].Print(tab + 1)
		fmt.Println()
	}

	node.Body.Print(tab)
	fmt.Println()
}

// --------------------------------------------------------

type Where struct {
	Id        Expression
	Condition Expression

	meta interface{}
}

func (m Where) IsMatch() {}

func NewWhere(i Identifier, c Expression) (Where, error) {
	return Where{
		i,
		c,
		nil,
	}, nil
}

func (A Where) EqualsMatch(b Match) bool {
	switch B := b.(type) {
	case Where:
		return A.Id.EqualsExpr(B.Id) && A.Condition.EqualsExpr(B.Condition)
	}

	return false
}

func (e Where) MetaGet() interface{} {
	return e.meta
}

func (e Where) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (node Where) Print(tab int) {
	fmt.Printf("(")
	node.Id.Print(0)
	fmt.Printf(" : \n")
	node.Condition.Print(tab + 1)
	fmt.Println()
	printTab(tab)
	fmt.Println(")")
}

// --------------------------------------------------------

func printTab(tab int) {
	for ; tab > 0; tab-- {
		fmt.Print("  ")
	}
}

func PrintTab(tab int) {
	for ; tab > 0; tab-- {
		fmt.Print("  ")
	}
}

func Print(ast Expression) {
	ast.Print(0)
}
