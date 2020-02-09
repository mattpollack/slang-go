package ast

// PROBLEM:
// List of values is impossible...

/*
import (
	"fmt"
	"reflect"
)

var builtins = []Builtin{
	{"+", func(a Value, env *Environment) Value {
		switch A := a.ast.(type) {
		case Number:
			return Value{Builtin{"+ val", func(b Value, env *Environment) Value {
				switch B := b.ast.(type) {
				case Number:
					return Value{Number{A.Value + B.Value}, env}
				}

				panic("Cannot + non-number")
			}}, env}
		}

		panic("Cannot + non-number")
	}},

	{"-", func(a Value, env *Environment) Value {
		switch A := a.ast.(type) {
		case Number:
			return Value{Builtin{"- val", func(b Value, env *Environment) Value {
				switch B := b.ast.(type) {
				case Number:
					return Value{Number{A.Value - B.Value}, env}
				}

				panic("Cannot - non-number")
			}}, env}
		}

		panic("Cannot - non-number")
	}},

	{"++", func(a Value, env *Environment) Value {
		switch A := a.ast.(type) {
		case List:
			return Value{Builtin{"++ val", func(b Value, env *Environment) Value {
				switch B := b.ast.(type) {
				case List:
					return Value{List{append(append([]AST{}, A.Values...), B.Values...)}, env}
				}

				panic("Cannot ++ non-list")
			}}, env}
		}

		panic("Cannot ++ non-list")
	}},

	{"print", func(val Value, env *Environment) Value {
		Print(val.ast)

		return val
	}},
}

func Interpret(ast AST) Value {
	env := newEnv(nil)

	for _, b := range builtins {
		env.Set(b.name, Value{b, env})
	}

	return visit(ast, env)
}

func visit(ast AST, env *Environment) Value {
	switch node := ast.(type) {
	case Let:
		env = newEnv(env)

		for i, id := range node.BoundIds {
			env.Set(id.Value, visit(node.BoundValues[i], env))
		}

		return visit(node.Body, env).Eval()

	case Identifier:
		if v, ok := env.Get(node.Value); ok {
			return v
		}

		panic(fmt.Sprintf("Cannot find value with identifier %s", node.Value))

	case Pattern:
		for range node.Bodies {
			node.Envs = append(node.Envs, newEnv(env))
		}

		return Value{node, env}

	}

	// Default evaluation
	return (&Value{ast, env}).Eval()
}

type Value struct {
	ast AST
	env *Environment
}

func (v Value) Print() {
	Print(v.ast)
	// EXTRA INFO
	// fmt.Println()
	// v.env.Print()
}

func (v Value) Eval() Value {
	switch node := v.ast.(type) {
	case Application:
		res := visit(node.Body[0], v.env)

		for _, argAst := range node.Body[1:] {
			res = res.Apply(visit(argAst, v.env))
		}

		return res

	default:
		return v
	}
}

func (v Value) Apply(arg Value) Value {
	switch node := v.ast.(type) {
	case Builtin:
		return node.apply(arg, v.env)

	case Pattern:
		return PatternApply(v, arg)

	default:
		panic(fmt.Sprintf("Cannot apply values of types %s and %s", reflect.TypeOf(v.ast), reflect.TypeOf(arg.ast)))
	}
}

// --------------------------------------------------------

func patternMatch(env *Environment, ast Pattern, res Pattern, m AST, val Value) bool {

	switch match := m.(type) {
	case ListConstructor:
		// Match a list
		if list, ok := val.ast.(List); ok {
			if len(list.Values) > 0 {
				return patternMatch(env, ast, res, match.Head, Value{list.Values[0], env}) && patternMatch(env, ast, res, match.Tail, Value{List{list.Values[1:]}, env})
			}
		}

		// Match a string
		if str, ok := val.ast.(String); ok {
			if len(str.Value) > 0 {
				return patternMatch(env, ast, res, match.Head, Value{String{string(str.Value[0])}, env}) && patternMatch(env, ast, res, match.Tail, Value{String{str.Value[1:]}, env})
			}
		}

	case Identifier:
		if v, ok := env.Get(match.Value); ok {
			return v.ast.Equals(val.ast)
		} else {
			env.Set(match.Value, val)

			return true
		}

	default:
		return match.Equals(val.ast)
	}

	return false
}

func PatternApply(pattern Value, val Value) Value {
	ast := pattern.ast.(Pattern)
	res, _ := NewPattern([][]AST{}, []AST{})

	for i, matchGroup := range ast.Matches {
		env := newEnv(ast.Envs[i])

		if patternMatch(env, ast, res, matchGroup[0], val) {
			res.Matches = append(res.Matches, matchGroup[1:])
			res.Bodies = append(res.Bodies, ast.Bodies[i])
			res.Envs = append(res.Envs, env)
		}
	}

	if len(res.Matches) == 0 {
		return Value{Label{"no_match"}, pattern.env}
	}

	if len(res.Matches[0]) == 0 {
		return visit(res.Bodies[0], res.Envs[0])
	}

	return Value{res, pattern.env}
}
*/
