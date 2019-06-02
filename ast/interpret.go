package ast

import (
	"fmt"
)

var True = Label{Value: "true"}
var False = Label{Value: "false"}

var global = map[string]func(Expression) Expression{
	"+": func(arg Expression) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression) Expression {
				switch A1 := arg.(type) {
				case Number:
					return Number{Value: A0.Value + A1.Value}
				}

				panic("Mismatching types passed to '+'")
			})
		}

		panic("Mismatching types passed to '+'")
	},
	"-": func(arg Expression) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression) Expression {
				switch A1 := arg.(type) {
				case Number:
					return Number{Value: A0.Value - A1.Value}
				}

				panic("Mismatching types passed to '-'")
			})
		}

		panic("Mismatching types passed to '-'")
	},
	"eql": func(a0 Expression) Expression {
		return NewNative(func(a1 Expression) Expression {
			if a0.Equals(a1) {
				return True
			} else {
				return False
			}
		})
	},
	"print": func(arg Expression) Expression {
		arg.Print(0)

		return arg
	},
}

func Interpret(expr Expression) Expression {
	env := NewEnvironment()

	for k, fn := range global {
		env = env.Set(k, NewNative(fn))
	}

	return expr.Eval(env)
}

// --------------------------------------------------------

type Native struct {
	fn  func(Expression) Expression
	env Environment
}

func NewNative(fn func(Expression) Expression) Native {
	return Native{fn, NewEnvironment()}
}

func (e Native) IsExpression() {}
func (e Native) IsMatch()      {}
func (e Native) Equals(b interface{}) bool {
	return false
}

func (e Native) Eval(env Environment) Expression {
	panic("e :(")
	return nil
}

func (e Native) Apply(arg Expression) Expression {
	return e.fn(arg)
}

func (e Native) MetaGet() interface{} {
	return nil
}

func (e Native) MetaSet(meta interface{}) interface{} {
	return nil
}

func (e Native) Print(tab int) {
	printTab(tab)
	fmt.Println("<native>")
}

// --------------------------------------------------------

func (e Application) Eval(env Environment) Expression {
	result := e.Body[0].Eval(env)

	if len(e.Body) == 1 {
		return result
	}

	for _, expr := range e.Body[1:] {
		result = result.Apply(expr.Eval(env))
	}

	return result
}

func (e Application) Apply(arg Expression) Expression {
	panic("Applications aren't values")
}

func (e If) Eval(env Environment) Expression {
	if e.Condition.Eval(env).Equals(True) {
		return e.Tbody.Eval(env)
	} else {
		return e.Fbody.Eval(env)
	}
}

func (e If) Apply(arg Expression) Expression {
	panic("a if")
	return nil
}

func (e Pattern) Eval(env Environment) Expression {
	e.env = env

	return e
}

func (e Pattern) Apply(arg Expression) Expression {
	res, _ := NewPattern([][]Match{}, []Expression{})

	// Copy env (move to method?)
	for k, v := range e.env.Bound {
		res.env = res.env.Set(k, v)
	}

	for i, match := range e.Matches {
		switch M := match[0].(type) {
		case Identifier:
			res.env = res.env.Set(M.Value, arg)
			res.Bodies = append(res.Bodies, e.Bodies[i])
			res.Matches = append(res.Matches, e.Matches[i][1:])
		case Number:
			if M.Equals(arg) {
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			}

		default:
			panic("TODO implement match of type")
		}
	}

	if len(res.Matches) == 0 {
		panic("Failed to match argument to pattern")
	}

	if len(res.Matches[0]) == 0 {
		return res.Bodies[0].Eval(res.env)
	}

	return res
}

func (e Identifier) Eval(env Environment) Expression {
	res := env.Get(e.Value)

	if res == nil {
		panic(fmt.Sprintf("Cannot find identifier '%s'", e.Value))
	}

	return res
}

func (e Identifier) Apply(arg Expression) Expression {
	panic("a identifier")
	return nil
}

func (e Label) Eval(env Environment) Expression {
	panic("e label")
	return nil
}

func (e Label) Apply(arg Expression) Expression {
	panic("a label")
	return nil
}

func (e String) Eval(env Environment) Expression {
	return e
}

func (e String) Apply(arg Expression) Expression {
	panic("a string")
	return nil
}

func (e Number) Eval(env Environment) Expression {
	return e
}

func (e Number) Apply(arg Expression) Expression {
	panic("a number")
	return nil
}

func (e Let) Eval(env Environment) Expression {
	for i, id := range e.BoundIds {
		env = env.Set(id.Value, e.BoundValues[i].Eval(env))
	}

	return e.Body.Eval(env)
}

func (e Let) Apply(arg Expression) Expression {
	panic("a let")
	return nil
}

func (m Where) Eval(env Environment) {
	panic("e where")
}

func (m Where) Apply(arg Expression) Expression {
	panic("a where")
}
