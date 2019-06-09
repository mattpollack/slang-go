package ast

import (
	"fmt"
)

func Interpret(expr Expression) Expression {
	env := NewEnvironment()

	for k, fn := range builtin {
		env = env.Set(k, NewNative(fn))
	}

	return expr.Eval(env)
}

// --------------------------------------------------------

var True = Label{Value: "true"}
var False = Label{Value: "false"}

var builtin = map[string]func(Expression) Expression{
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
	"*": func(arg Expression) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression) Expression {
				switch A1 := arg.(type) {
				case Number:
					return Number{Value: A0.Value * A1.Value}
				}

				panic("Mismatching types passed to '*'")
			})
		}

		panic("Mismatching types passed to '*'")
	},
	">": func(arg Expression) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression) Expression {
				switch A1 := arg.(type) {
				case Number:
					if A0.Value > A1.Value {
						return True
					} else {
						return False
					}
				}

				panic("Mismatching types passed to '>'")
			})
		}

		panic("Mismatching types passed to '>'")
	},
	">=": func(arg Expression) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression) Expression {
				switch A1 := arg.(type) {
				case Number:
					if A0.Value >= A1.Value {
						return True
					} else {
						return False
					}
				}

				panic("Mismatching types passed to '>='")
			})
		}

		panic("Mismatching types passed to '>='")
	},
	"<": func(arg Expression) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression) Expression {
				switch A1 := arg.(type) {
				case Number:
					if A0.Value < A1.Value {
						return True
					} else {
						return False
					}
				}

				panic("Mismatching types passed to '<'")
			})
		}

		panic("Mismatching types passed to '<'")
	},
	"abs": func(arg Expression) Expression {
		switch A0 := arg.(type) {
		case Number:
			if A0.Value < 0 {
				A0.Value = A0.Value * -1
			}

			return A0
		}

		panic("Mismatching types passed to 'abs'")
	},
	"&&": func(arg Expression) Expression {
		switch A0 := arg.(type) {
		case Label:
			return NewNative(func(arg Expression) Expression {
				switch A1 := arg.(type) {
				case Label:
					if A0.Equals(True) && A1.Equals(True) {
						return True
					} else {
						return False
					}
				}

				panic("Mismatching types passed to '-'")
			})
		}

		panic("Mismatching types passed to '-'")
	},
	"==": func(a0 Expression) Expression {
		return NewNative(func(a1 Expression) Expression {
			if a0.Equals(a1) {
				return True
			} else {
				return False
			}
		})
	},
	"print": func(arg Expression) Expression {
		switch A := arg.(type) {
		case Number:
			fmt.Print(A.Value)
		case String:
			fmt.Print(A.Value)
		case Label:
			fmt.Print(A.Value)

		// Only print patterns if they respond to ".print"
		case Pattern:
			l, _ := NewLabel("print")

			if p := A.Apply(l); p != nil {
				return p
			}

		default:
			arg.Print(0)
			panic("TODO print value of this type")
		}

		return arg
	},
	"print_ast": func(arg Expression) Expression {
		arg.Print(0)

		return arg
	},
}

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

func (e Slice) Eval(env Environment) Expression {
	if e.Low != nil {
		e.Low = e.Low.Eval(env)

		switch e.Low.(type) {
		case Number:
		default:
			panic("Slice low value must be a number")
		}
	}

	if e.High != nil {
		e.High = e.High.Eval(env)

		switch e.High.(type) {
		case Number:
		default:
			panic("Slice high value must be a number")
		}
	}

	return e
}

func (e Slice) Apply(arg Expression) Expression {
	// Apply the reverse
	return arg.Apply(e)
}

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

func (e List) Eval(env Environment) Expression {
	for i, _ := range e.Values {
		e.Values[i] = e.Values[i].Eval(env)
	}

	return e
}

func (e List) Apply(arg Expression) Expression {
	switch A := arg.(type) {
	case Label:
		if A.Value == "len" {
			n, _ := NewNumber(len(e.Values))
			return n
		}

		if A.Value == "head" {
			return e.Values[0]
		}

		if A.Value == "tail" {
			return List{Values: e.Values[1:]}
		}

		panic(fmt.Sprintf("List doesn't respond to .%s", A.Value))
	default:
		panic("Cannot apply this type to list")
	}
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
			if v := res.env.Get(M.Value); v == nil {
				res.env = res.env.Set(M.Value, arg)
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			} else if v.Equals(arg) {
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			}

		case Number:
			if M.Equals(arg) {
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			}

		case String:
			if M.Equals(arg) {
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			}

		case Where:
			if M.Condition.Eval(res.env.Set(M.Id.Value, arg)).Equals(True) {
				res.env = res.env.Set(M.Id.Value, arg)
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			}

		case Label:
			if M.Equals(arg) {
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			}

		case ListConstructor:
			switch A := arg.(type) {
			case String:
				sets := []func(){}
				head := false
				tail := false
				headOffset := 0

				// Match first value to head
				switch Head := M.Head.(type) {
				case Identifier:
					if v := res.env.Get(Head.Value); v == nil && len(A.Value) > 0 {
						sets = append(sets, func() {
							res.env = res.env.Set(Head.Value, String{Value: A.Value[:1]})
						})

						headOffset = 1
						head = true
					} else {
						switch V := v.(type) {
						case String:
							if len(A.Value) >= len(V.Value) && V.Equals(String{Value: A.Value[:len(V.Value)]}) {
								headOffset = len(V.Value)
								head = true
							}
						}
					}

				case String:
					if len(A.Value) >= len(Head.Value) && M.Head.Equals(String{Value: A.Value[:len(Head.Value)]}) {
						headOffset = len(Head.Value)
						head = true
					}
				}

				if !head {
					continue
				}

				// Match the rest of the values to tail
				switch Tail := M.Tail.(type) {
				case Identifier:
					if v := res.env.Get(Tail.Value); v == nil {
						sets = append(sets, func() {
							res.env = res.env.Set(Tail.Value, String{Value: A.Value[headOffset:]})
						})

						tail = true
					} else {
						switch V := v.(type) {
						case String:
							if V.Equals(String{Value: A.Value[headOffset:]}) {
								headOffset = len(V.Value)
								head = true
							}
						}
					}
				default:
					if Tail.Equals(String{Value: A.Value[headOffset:]}) {
						tail = true
					}
				}

				if tail {
					for _, fn := range sets {
						fn()
					}

					res.Bodies = append(res.Bodies, e.Bodies[i])
					res.Matches = append(res.Matches, e.Matches[i][1:])
				}

			case List:
				if len(A.Values) > 0 {
					sets := []func(){}
					head := false
					tail := false

					// Match first value to head
					switch Head := M.Head.(type) {
					case Identifier:
						sets = append(sets, func() {
							res.env = res.env.Set(Head.Value, A.Values[0])
						})

						head = true
					default:
						if M.Head.Equals(A.Values[0]) {
							head = true
						}
					}

					if !head {
						continue
					}

					// Match the rest of the values to tail
					switch Tail := M.Tail.(type) {
					case Identifier:
						sets = append(sets, func() {
							res.env = res.env.Set(Tail.Value, List{Values: A.Values[1:]})
						})

						tail = true
					default:
						if M.Tail.Equals(List{Values: A.Values[1:]}) {
							tail = true
						}
					}

					if tail {
						for _, fn := range sets {
							fn()
						}

						res.Bodies = append(res.Bodies, e.Bodies[i])
						res.Matches = append(res.Matches, e.Matches[i][1:])
					}
				}

			default:
				panic("Cannot match list construtor to value of this type")
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
	return e
}

func (e Label) Apply(arg Expression) Expression {
	panic("a label")
	return nil
}

func (e String) Eval(env Environment) Expression {
	return e
}

func (e String) Apply(arg Expression) Expression {
	switch A := arg.(type) {
	case Label:
		switch A.Value {
		case "len":
			num, _ := NewNumber(len(e.Value))
			return num
		default:
			panic("Invalid label applied to string")
		}

	case Slice:
		l := -1
		h := -1

		if A.Low != nil {
			l = A.Low.(Number).Value
		}

		if A.High != nil {
			h = A.High.(Number).Value
		}

		if l != -1 {
			e.Value = e.Value[l:]
		}

		if h != -1 {
			e.Value = e.Value[:h]
		}

		return e

	default:
		panic("Invalid type applied to string")
	}
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
