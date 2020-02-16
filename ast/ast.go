package ast

import (
	"errors"
	"fmt"
	"strings"
)

/*

COPY PASTE NODES LIST

Application
Pattern
Identifier
Label
String
List
ListConstructor
Number
Let
Where
*/

type Strings interface {
	String() []string
}

type Equals interface {
	Equals(interface{}) bool
}

type AST interface {
	Strings
	Equals

	Eval(*Environment) (AST, error)
	Apply(AST) (AST, error)

	Copy() AST
}

// Package name and imports
type SourceFileImport struct {
	Path string
	Name string
}

type SourceFile struct {
	PackageName string
	Imports     []SourceFileImport
	Definition  AST
}

func (s *SourceFile) Print() {
	fmt.Printf("package %s\n", s.PackageName)
	fmt.Println()

	for _, imp := range s.Imports {
		fmt.Printf("import \"%s\"", imp.Path)

		if imp.Name != "" {
			fmt.Printf(" as %s", imp.Name)
		}

		fmt.Println()
	}

	fmt.Println()

	Print(s.Definition)
}

type Application struct {
	Body []AST
}

func NewApplication(body []AST) (Application, error) {
	return Application{
		body,
	}, nil
}

type Pattern struct {
	Matches [][]AST
	Bodies  []AST
	Envs    []*Environment
}

func NewPattern(matchGroups [][]AST, bodies []AST) (Pattern, error) {
	return Pattern{
		matchGroups,
		bodies,
		[]*Environment{},
	}, nil
}

type Identifier struct {
	Value string
}

func NewIdentifier(v string) (Identifier, error) {
	return Identifier{
		v,
	}, nil
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

type String struct {
	Value string
}

func NewString(v string) (String, error) {
	return String{
		v,
	}, nil
}

type List struct {
	Values []AST
}

type ListConstructor struct {
	Head AST
	Tail AST
}

type Number struct {
	Value int
}

func NewNumber(v int) (Number, error) {
	return Number{
		v,
	}, nil
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

type Where struct {
	Match        AST
	Condition    AST
	ConstantTime bool
}

func NewWhere(m AST, c AST, ct bool) (Where, error) {
	return Where{
		m,
		c,
		ct,
	}, nil
}

func Print(ast AST) {
	fmt.Println(strings.Join(ast.String(), "\n"))
}

func addTab(str string) string {
	return fmt.Sprintf("  %s", str)
}

type Builtin struct {
	name  string
	apply func(AST, *Environment) (AST, error)
}

// -- Strings -------------------------

func (e Application) String() []string {
	res := []string{}
	multiLine := false

	for _, expr := range e.Body {
		strs := expr.String()

		if len(strs) == 1 {
			res = append(res, strs[0])
		} else {
			multiLine = true
			res = append(res, strs...)
		}
	}

	res = append([]string{"("}, res...)
	res = append(res, ")")

	if multiLine {
		for i := range res[1:][:len(res)-2] {
			res[i+1] = addTab(res[i+1])
		}

		return res
	} else {
		return []string{strings.Join(res, " ")}
	}
}

func (e Pattern) String() []string {
	res := []string{"{"}

	for i, matchGroup := range e.Matches {
		matchStrs := []string{}

		for _, match := range matchGroup {
			matchStrs = append(matchStrs, strings.Join(match.String(), ";"))
		}

		str := addTab(fmt.Sprintf("%s -> ", strings.Join(matchStrs, " ")))

		if body := e.Bodies[i].String(); len(body) == 1 {
			str += body[0]

			res = append(res, str)
		} else {
			for i := range body {
				body[i] = addTab(addTab(body[i]))
			}

			res = append(res, str)
			res = append(res, body...)
		}
	}

	res = append(res, "}")

	return res
}

func (e Identifier) String() []string {
	return []string{e.Value}
}

func (e Label) String() []string {
	return []string{fmt.Sprintf(".%s", e.Value)}
}

func (e String) String() []string {
	return []string{fmt.Sprintf(`"%s"`, e.Value)}
}

func (e List) String() []string {
	isMultiLine := false
	res := []string{}

	for _, val := range e.Values {
		strs := val.String()

		if len(strs) > 1 {
			isMultiLine = true
		}

		res = append(res, strs...)
	}

	if len(res) == 0 {
		return []string{"[]"}
	}

	if !isMultiLine {
		return []string{"[ " + strings.Join(res, ", ") + " ]"}
	} else {
		for i := range res {
			res[i] = addTab(res[i])
		}

		res = append([]string{"["}, res...)
		res = append(res, "]")
	}

	return res
}

func (e ListConstructor) String() []string {
	res := []string{"["}
	res = append(res, e.Head.String()...)
	res = append(res, ":")
	res = append(res, e.Tail.String()...)
	res = append(res, "]")

	return []string{strings.Join(res, " ")}
}

func (e Number) String() []string {
	return []string{fmt.Sprintf("%d", e.Value)}
}

func (e Let) String() []string {
	res := []string{}

	for i := range e.BoundIds {
		id := e.BoundIds[i]
		body := e.BoundValues[i].String()

		if len(body) == 1 {
			res = append(res, fmt.Sprintf("%s = %s", id.Value, body[0]))
		} else {
			for i := range body {
				body[i] = addTab(body[i])
			}

			res = append(res, fmt.Sprintf("%s = ", id.Value))
			res = append(res, body...)
		}
	}

	res = append(res, e.Body.String()...)

	return res
}

func (e Where) String() []string {
	return []string{strings.Join(append(append(append([]string{}, e.Match.String()...), ":"), e.Condition.String()...), " ")}
}

func (e Builtin) String() []string {
	return []string{fmt.Sprintf("<%s>", e.name)}
}

// -- Equals --------------------------

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

func (A Pattern) Equals(b interface{}) bool {
	panic("TODO pattern equal")
}

func (A Identifier) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Identifier:
		return A.Value == B.Value
	}

	return false
}

func (A Label) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Label:
		return A.Value == B.Value
	}

	return false
}

func (A String) Equals(b interface{}) bool {
	switch B := b.(type) {
	case String:
		return A.Value == B.Value
	}

	return false
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

func (A ListConstructor) Equals(b interface{}) bool {
	panic("TODO list constructor equals")
}

func (A Number) Equals(b interface{}) bool {
	switch B := b.(type) {
	case Number:
		return A.Value == B.Value
	}

	return false
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

func (A Where) Equals(b interface{}) bool {
	panic("TODO")
}

func (A Builtin) Equals(b interface{}) bool {
	panic("TODO")
}

// -- Copy --------------------------

func (a Application) Copy() AST {
	res := Application{}

	for _, ast := range a.Body {
		res.Body = append(res.Body, ast.Copy())
	}

	return res
}

func (a Pattern) Copy() AST {
	res := Pattern{}

	for _, matchGroup := range a.Matches {
		matchGroupCopy := []AST{}

		for _, ast := range matchGroup {
			matchGroupCopy = append(matchGroupCopy, ast.Copy())
		}

		res.Matches = append(res.Matches, matchGroupCopy)
	}

	for _, ast := range a.Bodies {
		res.Bodies = append(res.Bodies, ast.Copy())
	}

	return res
}

func (a Identifier) Copy() AST {
	return a
}

func (a Label) Copy() AST {
	return a
}

func (a String) Copy() AST {
	return String{"" + a.Value}
}

func (a List) Copy() AST {
	res := List{}

	for _, ast := range a.Values {
		res.Values = append(res.Values, ast.Copy())
	}

	return res
}

func (a ListConstructor) Copy() AST {
	return ListConstructor{
		a.Head.Copy(),
		a.Tail.Copy(),
	}
}

func (a Number) Copy() AST {
	return a
}

func (a Let) Copy() AST {
	res := Let{}

	for _, id := range res.BoundIds {
		res.BoundIds = append(res.BoundIds, id.Copy().(Identifier))
	}

	for _, ast := range res.BoundValues {
		res.BoundValues = append(res.BoundValues, ast.Copy())
	}

	res.Body = a.Body.Copy()

	return res
}

func (a Where) Copy() AST {
	return Where{
		a.Match.Copy(),
		a.Condition.Copy(),
		a.ConstantTime,
	}
}

func (a Builtin) Copy() AST {
	return a
}

// -- RUNTIME ---------------------------

type RuntimeError struct {
	wraps   error
	message string
}

func NewRuntimeError(wrapped error, message string) RuntimeError {
	return RuntimeError{
		wraps:   wrapped,
		message: message,
	}
}

func (r RuntimeError) Unwrap() error {
	return r.wraps
}

func (r RuntimeError) Error() string {
	messages := []string{}

	var err error = r

	for ; err != nil; err = errors.Unwrap(err) {
		if rerr, ok := err.(RuntimeError); ok {
			messages = append(messages, rerr.message)
		} else {
			messages = append(messages, err.Error())
		}
	}

	message := "RUNTIME ERROR:\n"

	for i := len(messages) - 1; i >= 0; i-- {
		message += fmt.Sprintf(" > %s\n", messages[i])
	}

	return message
}

type Environment struct {
	parent *Environment
	bound  map[string]AST
}

func NewEnv(parent *Environment) *Environment {
	return &Environment{parent, map[string]AST{}}
}

func (e *Environment) Set(key string, val AST) {
	e.bound[key] = val
}

func (e *Environment) Get(key string) (AST, bool) {
	if v, ok := e.bound[key]; ok {
		return v, true
	} else {
		if e.parent == nil {
			return nil, false
		}

		return e.parent.Get(key)
	}
}

// -- EVAL ------------------------------

func (a Application) Eval(env *Environment) (AST, error) {
	res, err := a.Body[0].Eval(env)

	if err != nil {
		return nil, NewRuntimeError(err, "Unable to evaluate for application")
	}

	for _, argAst := range a.Body[1:] {
		arg, err := argAst.Eval(env)

		if err != nil {
			return nil, NewRuntimeError(err, "Unable to evaluate for application")
		}

		res, err = res.Apply(arg)

		if err != nil {
			return nil, NewRuntimeError(err, "Unable to apply for application")
		}
	}

	return res, nil
}

func (a Pattern) Eval(env *Environment) (AST, error) {
	for range a.Bodies {
		a.Envs = append(a.Envs, NewEnv(env))
	}

	return a, nil
}

func (a Identifier) Eval(env *Environment) (AST, error) {
	if v, ok := env.Get(a.Value); ok {
		return v, nil
	}

	return nil, NewRuntimeError(nil, fmt.Sprintf("Cannot get value for identifier '%s'", a.Value))
}

func (a Number) Eval(*Environment) (AST, error) { return a, nil }
func (a Label) Eval(*Environment) (AST, error)  { return a, nil }
func (a String) Eval(*Environment) (AST, error) { return a, nil }
func (a List) Eval(env *Environment) (AST, error) {
	res := List{}

	for _, valAst := range a.Values {
		val, err := valAst.Eval(env)

		if err != nil {
			return nil, NewRuntimeError(err, "Unable to evaluate for list")
		}

		res.Values = append(res.Values, val)
	}

	return res, nil
}
func (a ListConstructor) Eval(*Environment) (AST, error) { panic("TODO eval list con") }
func (a Let) Eval(env *Environment) (AST, error) {
	env = NewEnv(env)

	for i, id := range a.BoundIds {
		val, err := a.BoundValues[i].Eval(env)

		if err != nil {
			return nil, NewRuntimeError(err, "Unable to evaluate for let")
		}

		env.Set(id.Value, val)
	}

	return a.Body.Eval(env)
}
func (a Where) Eval(*Environment) (AST, error)   { panic("TODO eval where") }
func (a Builtin) Eval(*Environment) (AST, error) { panic("TODO eval built") }

// -- APPLY -----------------------------

func (a Application) Apply(b AST) (AST, error)     { panic("TODO apply app") }
func (a Identifier) Apply(b AST) (AST, error)      { panic("TODO apply id") }
func (a Label) Apply(b AST) (AST, error)           { panic("TODO apply lab") }
func (a String) Apply(b AST) (AST, error)          { panic("TODO apply str") }
func (a List) Apply(b AST) (AST, error)            { panic("TODO apply list") }
func (a ListConstructor) Apply(b AST) (AST, error) { panic("TODO apply list con") }
func (a Number) Apply(b AST) (AST, error)          { panic("TODO apply num") }
func (a Let) Apply(b AST) (AST, error)             { panic("TODO apply let") }
func (a Where) Apply(b AST) (AST, error)           { panic("TODO apply where") }
func (a Builtin) Apply(b AST) (AST, error) {
	return a.apply(b, nil)
}

func (a Pattern) Apply(val AST) (AST, error) {
	res, _ := NewPattern([][]AST{}, []AST{})

	for i, matchGroup := range a.Matches {
		env := NewEnv(a.Envs[i])

		if patternMatch(env, a, res, matchGroup[0], val) {
			res.Matches = append(res.Matches, matchGroup[1:])
			res.Bodies = append(res.Bodies, a.Bodies[i])
			res.Envs = append(res.Envs, env)
		}
	}

	if len(res.Matches) == 0 {
		return nil, NewRuntimeError(nil, "Unable to match any values to pattern")
	}

	if len(res.Matches[0]) == 0 {
		return res.Bodies[0].Eval(res.Envs[0])
	}

	return res, nil
}

func patternMatch(env *Environment, ast Pattern, res Pattern, m AST, val AST) bool {
	switch match := m.(type) {
	case Where:
		if patternMatch(env, ast, res, match.Match, val) {
			res, err := match.Condition.Eval(env)

			return err == nil && res.Equals(Label{Value: "true"})
		}

	case ListConstructor:
		// Match for a list
		if list, ok := val.(List); ok {
			if len(list.Values) > 0 {
				return patternMatch(env, ast, res, match.Head, list.Values[0]) && patternMatch(env, ast, res, match.Tail, List{list.Values[1:]})
			}
		}

		// Match for a string
		if str, ok := val.(String); ok {
			if len(str.Value) > 0 {
				return patternMatch(env, ast, res, match.Head, String{string(str.Value[0])}) && patternMatch(env, ast, res, match.Tail, String{str.Value[1:]})
			}
		}

	case Identifier:
		if v, ok := env.Get(match.Value); ok {
			return v.Equals(val)
		} else {
			env.Set(match.Value, val)

			return true
		}

	default:
		return match.Equals(val)
	}

	return false
}
