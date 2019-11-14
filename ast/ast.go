package ast

import (
	"fmt"
)

type Prints interface {
	Print(int)
}

type Equals interface {
	Equals(interface{}) bool
}

type AST interface {
	Prints
	Equals

	//FreeVars() FreeVarsPack
}

type FreeVarsPack struct {
	Nodes []AST
}

type Application struct {
	Body []AST
}

func NewApplication(body []AST) (Application, error) {
	return Application{
		body,
	}, nil
}

func (A Application) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Application:
		if len(A.Body) == len(B.Body) {
			for i := 0; i < len(A.Body); i++ {
				if !A.Body[i].Equals(B.Body[i]) {
					return false
				}
			}

			return true
		}
	}

	return false
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

type Pattern struct {
	Matches [][]AST
	Bodies  []AST
}

func NewPattern(m [][]AST, b []AST) (Pattern, error) {
	return Pattern{
		m,
		b,
	}, nil
}

func (A Pattern) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Pattern:
		if len(A.Matches) == len(B.Matches) && (len(A.Matches) == 0 || len(A.Matches[0]) == len(B.Matches[0])) && len(A.Bodies) == len(B.Bodies) {
			for i := 0; i < len(A.Matches); i++ {
				for j := 0; j < len(A.Matches[i]); j++ {
					if !A.Matches[i][j].Equals(B.Matches[i][j]) {
						return false
					}
				}
			}

			for i := 0; i < len(A.Bodies); i++ {
				if !A.Bodies[i].Equals(B.Bodies[i]) {
					return false
				}
			}

			return true
		}
	}

	return false
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

type Identifier struct {
	Value string
}

func NewIdentifier(v string) (Identifier, error) {
	return Identifier{
		v,
	}, nil
}

func (A Identifier) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Identifier:
		return A.Value == B.Value
	}

	return false
}

func (node Identifier) Print(tab int) {
	printTab(tab)
	fmt.Println(node.Value)
}

type Label struct {
	Value string
}

func NewLabel(v string) (Label, error) {
	return Label{
		v,
	}, nil
}

var True, _ = NewLabel("true")
var False, _ = NewLabel("false")

func (A Label) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Label:
		return A.Value == B.Value
	}

	return false
}

func (node Label) Print(tab int) {
	printTab(tab)
	fmt.Printf(".%s", node.Value)
}

type String struct {
	Value string
}

func NewString(v string) (String, error) {
	return String{
		v,
	}, nil
}

func (A String) Equals(b interface{}) bool {
	switch B := b.(type) {
	case String:
		return A.Value == B.Value
	}

	return false
}

func (node String) Print(tab int) {
	printTab(tab)
	fmt.Println(node.Value)
}

type List struct {
	Values []AST
}

func (A List) Equals(b interface{}) bool {
	switch B := b.(type) {
	case List:
		if len(A.Values) != len(B.Values) {
			return false
		}

		for i := 0; i < len(A.Values); i++ {
			if !A.Values[i].Equals(B.Values[i]) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

func (node List) Print(tab int) {
	printTab(tab)
	fmt.Println("[")

	for _, val := range node.Values {
		val.Print(tab + 1)
	}

	printTab(tab)
	fmt.Println("]")
}

type ListConstructor struct {
	Head AST
	Tail AST
}

func (A ListConstructor) Equals(b interface{}) bool {
	panic("TODO list constructor equals")
}

func (node ListConstructor) Print(tab int) {
	printTab(tab)
	fmt.Println("[")

	node.Head.Print(tab + 1)
	fmt.Println()
	node.Tail.Print(tab + 1)
	fmt.Println()

	printTab(tab)
	fmt.Println("]")
}

type Number struct {
	Value int
}

func NewNumber(v int) (Number, error) {
	return Number{
		v,
	}, nil
}

func (A Number) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Number:
		return A.Value == B.Value
	}

	return false
}

func (node Number) Print(tab int) {
	printTab(tab)
	fmt.Println(node.Value)
}

type Let struct {
	BoundIds    []Identifier
	BoundValues []AST
	Body        AST
}

func NewLet(ids []Identifier, vs []AST, b AST) (Let, error) {
	return Let{
		ids,
		vs,
		b,
	}, nil
}

func (A Let) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Let:
		if len(A.BoundIds) == len(B.BoundIds) {
			for i := 0; i < len(A.BoundIds); i++ {
				if !A.BoundIds[i].Equals(B.BoundIds[i]) || !A.BoundValues[i].Equals(B.BoundValues[i]) {
					return false
				}
			}

			return A.Body.Equals(B.Body)
		}
	}

	return false
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

type Where struct {
	Id        Identifier
	Condition AST
}

func NewWhere(i Identifier, c AST) (Where, error) {
	return Where{
		i,
		c,
	}, nil
}

func (A Where) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Where:
		return A.Id.Equals(B.Id) && A.Condition.Equals(B.Condition)
	}

	return false
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

func Print(ast AST) {
	ast.Print(0)
}
