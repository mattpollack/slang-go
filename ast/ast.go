package ast

import (
	"fmt"
)

type Meta interface {
	MetaSet(interface{}) interface{}
	MetaGet() interface{}
}

type Expression interface {
	Meta
	EqualsExpr(Expression) bool

	IsExpression()
}

type Match interface {
	Meta
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

// --------------------------------------------------------

func printTab(tab int) {
	for ; tab > 0; tab-- {
		fmt.Print("  ")
	}
}

func PrintMatch(m Match, tab int) {
	switch node := m.(type) {
	case Identifier:
		fmt.Print(node.Value)
	case Label:
		fmt.Printf(".%s", node.Value)
	case String:
		fmt.Print("\"%s\"", node.Value)
	case Number:
		fmt.Print(node.Value)
	case Where:
		fmt.Printf("(")
		printHelp(node.Id, 0)
		fmt.Printf(" : \n")
		printHelp(node.Condition, tab+1)
		fmt.Println()
		printTab(tab)
		fmt.Print(")")
	default:
		fmt.Print("<Match>")
	}
}

func Print(ast Expression) {
	printHelp(ast, 0)
	fmt.Println()
}

func printHelp(ast Expression, tab int) {
	if ast == nil {
		printTab(tab)
		fmt.Println("<nil>")
		return
	}

	switch node := ast.(type) {
	case Application:
		printTab(tab)
		fmt.Println("(")

		for _, e := range node.Body {
			printHelp(e, tab+1)
			fmt.Println()
		}

		printTab(tab)
		fmt.Print(")")
	case Pattern:
		printTab(tab)
		fmt.Println("{")

		for i := 0; i < len(node.Bodies); i++ {
			printTab(tab)

			for _, m := range node.Matches[i] {
				PrintMatch(m, tab)
				fmt.Print(" ")
			}

			fmt.Println("-> ")
			printHelp(node.Bodies[i], tab+1)
			fmt.Println()
		}
		printTab(tab)
		fmt.Print("}")
	case Identifier:
		printTab(tab)
		fmt.Print(node.Value)
	case Label:
		printTab(tab)
		fmt.Printf(".%s", node.Value)
	case String:
		printTab(tab)
		fmt.Printf("\"%s\"", node.Value)
	case Number:
		printTab(tab)
		fmt.Print(node.Value)
	case Let:
		for i, id := range node.BoundIds {
			printTab(tab)
			fmt.Print("let ")
			printHelp(id, 0)
			fmt.Println(" =")
			printHelp(node.BoundValues[i], tab+1)
			fmt.Println()
		}

		printHelp(node.Body, tab)
		fmt.Println()
	case If:
		printTab(tab)
		fmt.Println("if ")
		printHelp(node.Condition, tab+1)
		fmt.Println()
		printHelp(node.Tbody, tab+1)
		fmt.Println()
		printTab(tab)
		fmt.Println("else ")
		printHelp(node.Fbody, tab+1)
	default:
		printTab(tab)
		fmt.Printf("<Unexpected>")
	}
}
