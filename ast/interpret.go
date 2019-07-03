package ast

import (
	"fmt"
)

func Interpret(expr Expression) Expression {
	env := NewEnvironment()

	for k, fn := range builtin {
		env = env.Set(k, NewNative(fn))
	}

	return expr.Eval(env, false)
}

// --------------------------------------------------------

var True = Label{Value: "true"}
var False = Label{Value: "false"}

func Panic(safe bool, msg string) Label {
	if !safe {
		panic(msg)
	} else {
		return Label{Value: "false"}
	}
}

var builtin = map[string]func(Expression, bool) Expression{
	"+": func(arg Expression, safe bool) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression, safe bool) Expression {
				switch A1 := arg.(type) {
				case Number:
					return Number{Value: A0.Value + A1.Value}
				}

				return Panic(safe, "Mismatching types passed to '+'")
			})
		}

		return Panic(safe, "Mismatching types passed to '+'")
	},
	"-": func(arg Expression, safe bool) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression, safe bool) Expression {
				switch A1 := arg.(type) {
				case Number:
					return Number{Value: A0.Value - A1.Value}
				}

				return Panic(safe, "Mismatching types passed to '-'")
			})
		}

		return Panic(safe, "Mismatching types passed to '-'")
	},
	"*": func(arg Expression, safe bool) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression, safe bool) Expression {
				switch A1 := arg.(type) {
				case Number:
					return Number{Value: A0.Value * A1.Value}
				}

				return Panic(safe, "Mismatching types passed to '*'")
			})
		}

		return Panic(safe, "Mismatching types passed to '*'")
	},
	">": func(arg Expression, safe bool) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression, safe bool) Expression {
				switch A1 := arg.(type) {
				case Number:
					if A0.Value > A1.Value {
						return True
					} else {
						return False
					}
				}

				return Panic(safe, "Mismatching types passed to '>'")
			})
		}

		return Panic(safe, "Mismatching types passed to '>'")
	},
	">=": func(arg Expression, safe bool) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression, safe bool) Expression {
				switch A1 := arg.(type) {
				case Number:
					if A0.Value >= A1.Value {
						return True
					} else {
						return False
					}
				}

				return Panic(safe, "Mismatching types passed to '>='")
			})
		}

		return Panic(safe, "Mismatching types passed to '>='")
	},
	"<": func(arg Expression, safe bool) Expression {
		switch A0 := arg.(type) {
		case Number:
			return NewNative(func(arg Expression, safe bool) Expression {
				switch A1 := arg.(type) {
				case Number:
					if A0.Value < A1.Value {
						return True
					} else {
						return False
					}
				}

				return Panic(safe, "Mismatching types passed to '<'")
			})
		}

		return Panic(safe, "Mismatching types passed to '<'")
	},
	"abs": func(arg Expression, safe bool) Expression {
		switch A0 := arg.(type) {
		case Number:
			if A0.Value < 0 {
				A0.Value = A0.Value * -1
			}

			return A0
		}

		return Panic(safe, "Mismatching types passed to 'abs'")
	},
	"&&": func(arg Expression, safe bool) Expression {
		switch A0 := arg.(type) {
		case Label:
			return NewNative(func(arg Expression, safe bool) Expression {
				switch A1 := arg.(type) {
				case Label:
					if A0.Equals(True) && A1.Equals(True) {
						return True
					} else {
						return False
					}
				}

				return Panic(safe, "Mismatching types passed to '-'")
			})
		}

		return Panic(safe, "Mismatching types passed to '-'")
	},
	"==": func(a0 Expression, safe bool) Expression {
		return NewNative(func(a1 Expression, safe bool) Expression {
			if a0.Equals(a1) {
				return True
			} else {
				return False
			}
		})
	},
	"!=": func(a0 Expression, safe bool) Expression {
		return NewNative(func(a1 Expression, safe bool) Expression {
			if a0.Equals(a1) {
				return False
			} else {
				return True
			}
		})
	},
	"++": func(arg Expression, safe bool) Expression {
		switch A0 := arg.(type) {
		case List:
			return NewNative(func(arg Expression, safe bool) Expression {
				switch A1 := arg.(type) {
				case List:
					return List{Values: append(append([]Expression{}, A0.Values...), A1.Values...)}
				}

				return Panic(safe, "Non-slice passed to '++'")
			})
		}

		return Panic(safe, "Non-slice passed to '++'")
	},
	"print": func(arg Expression, safe bool) Expression {
		switch A := arg.(type) {
		case Number:
			fmt.Print(A.Value)
		case String:
			fmt.Print(A.Value)
		case Label:
			fmt.Print(A.Value)

		// Only print patterns if they respond to ".print"
		case Pattern:
			arg.Print(0)
			fmt.Println()

		default:
			panic("TODO print value of this type")
		}

		return arg
	},
	"print_ast": func(arg Expression, safe bool) Expression {
		arg.Print(0)

		return arg
	},
}

type Native struct {
	fn  func(Expression, bool) Expression
	env Environment
}

func NewNative(fn func(Expression, bool) Expression) Native {
	return Native{fn, NewEnvironment()}
}

func (e Native) IsExpression() {}
func (e Native) IsMatch()      {}
func (e Native) Equals(b interface{}) bool {
	return false
}

func (e Native) Eval(Environment, bool) Expression {
	panic("e :(")
}

func (e Native) Apply(arg Expression, safe bool) Expression {
	return e.fn(arg, safe)
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

func (e Slice) Eval(env Environment, safe bool) Expression {
	if e.Low != nil {
		e.Low = e.Low.Eval(env, safe)

		switch e.Low.(type) {
		case Number:
		default:
			return Panic(safe, "Slice low value must be a number")
		}
	}

	if e.High != nil {
		e.High = e.High.Eval(env, safe)

		switch e.High.(type) {
		case Number:
		default:
			return Panic(safe, "Slice high value must be a number")
		}
	}

	return e
}

func (e Slice) Apply(arg Expression, safe bool) Expression {
	// Apply the reverse
	return arg.Apply(e, safe)
}

func (e Application) Eval(env Environment, safe bool) Expression {
	result := e.Body[0].Eval(env, safe)

	if len(e.Body) == 1 {
		return result
	}

	for _, expr := range e.Body[1:] {
		result = result.Apply(expr.Eval(env, safe), safe)
	}

	return result
}

func (e Application) Apply(arg Expression, safe bool) Expression {
	panic("Applications aren't values")
}

func (e List) Eval(env Environment, safe bool) Expression {
	for i, _ := range e.Values {
		e.Values[i] = e.Values[i].Eval(env, safe)
	}

	return e
}

func (e List) Apply(arg Expression, safe bool) Expression {
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

func (e If) Eval(env Environment, safe bool) Expression {
	if e.Condition.Eval(env, safe).Equals(True) {
		return e.Tbody.Eval(env, safe)
	} else {
		return e.Fbody.Eval(env, safe)
	}
}

func (e If) Apply(arg Expression, safe bool) Expression {
	panic("a if")
	return nil
}

func (e Pattern) Eval(env Environment, safe bool) Expression {
	e.env = env

	return e
}

func (e Pattern) Apply(arg Expression, safe bool) Expression {
	res, _ := NewPattern([][]Match{}, []Expression{})

	// Copy env (move to method?)
	for k, v := range e.env.Bound {
		res.env = res.env.Set(k, v)
	}

	for i, match := range e.Matches {
		switch M := match[0].(type) {
		case Identifier:
			// NOTE: identifiers are always reassigned (since letrec support isn't in yet)
			// Apply an identifier like (a) to match it's value
			res.env = res.env.Set(M.Value, arg)
			res.Bodies = append(res.Bodies, e.Bodies[i])
			res.Matches = append(res.Matches, e.Matches[i][1:])

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
			if !M.Condition.Eval(res.env.Set(M.Id.Value, arg), true).Equals(False) {
				res.env = res.env.Set(M.Id.Value, arg)
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			}

		case Label:
			if M.Equals(arg) {
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			}

		case List:
			if M.Equals(arg) {
				res.Bodies = append(res.Bodies, e.Bodies[i])
				res.Matches = append(res.Matches, e.Matches[i][1:])
			}

		case ListConstructor:
			// Evaluate non-special cases for the head value
			switch M.Head.(type) {
			case Application:
				M.Head = M.Head.Eval(res.env, true)
			}

			// Evaluate non-special cases for the tail value
			switch M.Tail.(type) {
			case Application:
				M.Tail = M.Tail.Eval(res.env, true)
			}

			// NOTE: possibly redo this, it's a bit everywhere

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
				Panic(true, "Cannot match list constructor to value of this type")
			}

		default:
			panic("TODO implement match of type")
		}
	}

	if len(res.Matches) == 0 {
		if !safe {
			fmt.Println("Argument:")
			arg.Print(1)
			fmt.Println()
			fmt.Println("Pattern:")
			e.Print(1)
			fmt.Println()
			fmt.Println()
		}

		return Panic(safe, "Failed to match argument to pattern")
	}

	if len(res.Matches[0]) == 0 {
		return res.Bodies[0].Eval(res.env, false)
	}

	return res
}

func (e Identifier) Eval(env Environment, safe bool) Expression {
	res := env.Get(e.Value)

	if res == nil {
		return Panic(safe, fmt.Sprintf("Cannot find identifier '%s'", e.Value))
	}

	return res
}

func (e Identifier) Apply(arg Expression, safe bool) Expression {
	panic("a identifier")
	return nil
}

func (e Label) Eval(Environment, bool) Expression {
	return e
}

func (e Label) Apply(arg Expression, safe bool) Expression {
	arg.Print(0)
	e.Print(0)

	panic("a label")
	return nil
}

func (e String) Eval(Environment, bool) Expression {
	return e
}

func (e String) Apply(arg Expression, safe bool) Expression {
	switch A := arg.(type) {
	case Label:
		switch A.Value {
		case "len":
			num, _ := NewNumber(len(e.Value))
			return num
		default:
			return Panic(safe, "Invalid label applied to string")
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
		return Panic(safe, "Invalid type applied to string")
	}
}

func (e Number) Eval(Environment, bool) Expression {
	return e
}

func (e Number) Apply(arg Expression, safe bool) Expression {
	panic("a number")
	return nil
}

func (e Let) Eval(env Environment, safe bool) Expression {
	for i, id := range e.BoundIds {
		env = env.Set(id.Value, e.BoundValues[i].Eval(env, safe))
	}

	return e.Body.Eval(env, safe)
}

func (e Let) Apply(arg Expression, safe bool) Expression {
	panic("a let")
	return nil
}

func (m Where) Eval(Environment, bool) {
	panic("e where")
}

func (m Where) Apply(arg Expression, safe bool) Expression {
	panic("a where")
}
