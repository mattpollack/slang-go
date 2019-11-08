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

type Equals interface {
	Equals(interface{}) bool
}

// A bit of a misnomer
type Expression interface {
	Meta
	Prints
	Equals

	Eval(Environment, bool) Expression
	Apply(Expression, bool) Expression
	Copy() Expression

	IsExpression()
}

type Match interface {
	Meta
	Prints
	Equals

	IsMatch()
}

// --------------------------------------------------------
// See interpret.go

type Environment struct {
	Bound map[string]Expression
}

func NewEnvironment() Environment {
	return Environment{map[string]Expression{}}
}

func (e Environment) Set(id string, val Expression) Environment {
	if id != "_" {
		e.Bound[id] = val
	}

	return e
}

func (e Environment) Get(id string) (Expression, bool) {
	exp, ok := e.Bound[id]

	if ok {
		return exp.Copy(), true
	} else {
		return nil, false
	}
}

func (e Environment) Unset(id string) Environment {
	delete(e.Bound, id)

	return e
}

func (e Environment) Copy() Environment {
	bound := map[string]Expression{}

	for k, v := range e.Bound {
		bound[k] = v
	}

	return Environment{bound}
}

// --------------------------------------------------------

type Application struct {
	Body []Expression

	meta interface{}
}

func (e Application) IsExpression() {}
func (e Application) Copy() Expression {
	body := []Expression{}

	for _, b := range e.Body {
		body = append(body, b.Copy())
	}

	return Application{body, e.meta}
}

func NewApplication(body []Expression) (Application, error) {
	return Application{
		body,
		nil,
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

func (A If) Equals(b interface{}) bool {
	switch B := b.(type) {
	case If:
		return A.Condition.Equals(B.Condition) && A.Tbody.Equals(B.Tbody) && A.Fbody.Equals(B.Fbody)
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

	env Environment

	meta interface{}
}

func (e Pattern) IsExpression() {}

func (e Pattern) Copy() Expression {
	bodies := []Expression{}

	for _, b := range e.Bodies {
		bodies = append(bodies, b.Copy())
	}

	return Pattern{
		e.Matches,
		bodies,
		e.env.Copy(),
		e.meta,
	}
}

func NewPattern(m [][]Match, b []Expression) (Pattern, error) {
	return Pattern{
		m,
		b,
		NewEnvironment(),
		nil,
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
func (e Identifier) Copy() Expression {
	return e
}

func NewIdentifier(v string) (Identifier, error) {
	return Identifier{
		v,
		nil,
	}, nil
}

func (A Identifier) Equals(b interface{}) bool {
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
func (e Label) Copy() Expression {
	return e
}

func NewLabel(v string) (Label, error) {
	return Label{
		v,
		nil,
	}, nil
}

func (A Label) Equals(b interface{}) bool {
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
	fmt.Printf(".%s", node.Value)
}

// --------------------------------------------------------

type String struct {
	Value string

	meta interface{}
}

func (e String) IsExpression() {}
func (e String) IsMatch()      {}
func (e String) Copy() Expression {
	return e
}

func NewString(v string) (String, error) {
	return String{
		v,
		nil,
	}, nil
}

func (A String) Equals(b interface{}) bool {
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

type Slice struct {
	Low  Expression
	High Expression

	meta interface{}
}

func (e Slice) IsExpression() {}
func (e Slice) Copy() Expression {
	return Slice{
		e.Low.Copy(),
		e.High.Copy(),
		e.meta,
	}
}

func (e Slice) Equals(b interface{}) bool {
	panic("TODO slice equals")
}

func (e Slice) MetaGet() interface{} {
	return e.meta
}

func (e Slice) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (node Slice) Print(tab int) {
	panic("TODO print slice")
}

// --------------------------------------------------------

type List struct {
	Values []Expression

	meta interface{}
}

func (e List) IsExpression() {}
func (e List) IsMatch()      {}
func (e List) Copy() Expression {
	values := []Expression{}

	for _, v := range e.Values {
		values = append(values, v.Copy())
	}

	return List{values, e.meta}
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

func (e List) MetaGet() interface{} {
	return e.meta
}

func (e List) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
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

// --------------------------------------------------------

type ListConstructor struct {
	Head Expression
	Tail Expression

	meta interface{}
}

func (e ListConstructor) IsExpression() {}
func (e ListConstructor) IsMatch()      {}

func (A ListConstructor) Equals(b interface{}) bool {
	panic("TODO list constructor equals")
}

func (e ListConstructor) MetaGet() interface{} {
	return e.meta
}

func (e ListConstructor) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
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

// --------------------------------------------------------

type Number struct {
	Value int

	meta interface{}
}

func (e Number) IsExpression() {}
func (e Number) IsMatch()      {}
func (e Number) Copy() Expression {
	return e
}

func NewNumber(v int) (Number, error) {
	return Number{
		v,
		nil,
	}, nil
}

func (A Number) Equals(b interface{}) bool {
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
func (e Let) Copy() Expression {
	bvs := []Expression{}

	for _, bv := range e.BoundValues {
		bvs = append(bvs, bv.Copy())
	}

	return Let{
		e.BoundIds,
		bvs,
		e.Body.Copy(),
		e.meta,
	}
}

func NewLet(ids []Identifier, vs []Expression, b Expression) (Let, error) {
	return Let{
		ids,
		vs,
		b,
		nil,
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
	Id        Identifier
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

func (A Where) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Where:
		return A.Id.Equals(B.Id) && A.Condition.Equals(B.Condition)
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
