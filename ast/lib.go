package ast

var isLib = map[string]bool{}
var StdLib Let

func init() {
	for _, b := range libFns {
		isLib[b.name] = true
	}

	StdLib, _ = NewLet([]Identifier{}, []AST{}, nil)

	for _, b := range libFns {
		StdLib.Bind(Identifier{b.name}, b)
	}

	StdLib.Body = Identifier{"NO BODY"}
}

func IsBuiltin(v string) bool {
	if _, ok := isLib[v]; ok {
		return true
	}

	return false
}

var libFns = []Builtin{
	{
		"+",
		func(a AST, env *Environment) (AST, error) {
			switch A := a.(type) {
			case Number:
				return Builtin{
					"+ curried",
					func(b AST, env *Environment) (AST, error) {
						switch B := b.(type) {
						case Number:
							return Number{A.Value + B.Value}, nil
						}

						return nil, NewRuntimeError(nil, "Can't add non-integer type")
					},
				}, nil
			}

			return nil, NewRuntimeError(nil, "Can't add non-integer type")
		},
	},

	{
		"*",
		func(a AST, env *Environment) (AST, error) {
			switch A := a.(type) {
			case Number:
				return Builtin{
					"* curried",
					func(b AST, env *Environment) (AST, error) {
						switch B := b.(type) {
						case Number:
							return Number{A.Value * B.Value}, nil
						}

						return nil, NewRuntimeError(nil, "Can't multiply non-integer type")
					},
				}, nil
			}

			return nil, NewRuntimeError(nil, "Can't multiply non-integer type")
		},
	},

	{
		"/",
		func(a AST, env *Environment) (AST, error) {
			switch A := a.(type) {
			case Number:
				return Builtin{
					"/ curried",
					func(b AST, env *Environment) (AST, error) {
						switch B := b.(type) {
						case Number:
							return Number{A.Value / B.Value}, nil
						}

						return nil, NewRuntimeError(nil, "Can't divide non-integer type")
					},
				}, nil
			}

			return nil, NewRuntimeError(nil, "Can't divide non-integer type")
		},
	},

	{
		"%",
		func(a AST, env *Environment) (AST, error) {
			switch A := a.(type) {
			case Number:
				return Builtin{
					"% curried",
					func(b AST, env *Environment) (AST, error) {
						switch B := b.(type) {
						case Number:
							return Number{A.Value % B.Value}, nil
						}

						return nil, NewRuntimeError(nil, "Can't modulo non-integer type")
					},
				}, nil
			}

			return nil, NewRuntimeError(nil, "Can't modulo non-integer type")
		},
	},

	{
		"-",
		func(a AST, env *Environment) (AST, error) {
			switch A := a.(type) {
			case Number:
				return Builtin{
					"- curried",
					func(b AST, env *Environment) (AST, error) {
						switch B := b.(type) {
						case Number:
							return Number{A.Value - B.Value}, nil
						}

						return nil, NewRuntimeError(nil, "Can't subtract non-integer type")
					},
				}, nil
			}

			return nil, NewRuntimeError(nil, "Can't add non-integer type")
		},
	},

	{
		">",
		func(a AST, env *Environment) (AST, error) {
			switch A := a.(type) {
			case Number:
				return Builtin{
					"> curried",
					func(b AST, env *Environment) (AST, error) {
						switch B := b.(type) {
						case Number:
							if A.Value > B.Value {
								return Label{"true"}, nil
							}

							return Label{"false"}, nil
						}

						return nil, NewRuntimeError(nil, "Can't apply greater than on non-integer type")
					},
				}, nil
			}

			return nil, NewRuntimeError(nil, "Can't apply greater than on non-integer type")
		},
	},

	{
		"||",
		func(a AST, env *Environment) (AST, error) {
			switch A := a.(type) {
			case Label:
				return Builtin{
					"|| curried",
					func(b AST, env *Environment) (AST, error) {
						switch B := b.(type) {
						case Label:
							if A.Value == "true" || B.Value == "true" {
								return Label{"true"}, nil
							}

							return Label{"false"}, nil
						}

						return nil, NewRuntimeError(nil, "Can't apply boolean or on non-boolean type")
					},
				}, nil
			}

			return nil, NewRuntimeError(nil, "Can't apply boolean or on non-boolean type")
		},
	},

	{
		"==",
		func(a AST, env *Environment) (AST, error) {
			return Builtin{
				"== curried",
				func(b AST, env *Environment) (AST, error) {
					if a.Equals(b) {
						return Label{"true"}, nil
					}

					return Label{"false"}, nil
				},
			}, nil
		},
	},

	{
		"!=",
		func(a AST, env *Environment) (AST, error) {
			return Builtin{
				"!= curried",
				func(b AST, env *Environment) (AST, error) {
					if !a.Equals(b) {
						return Label{"true"}, nil
					}

					return Label{"false"}, nil
				},
			}, nil
		},
	},

	{
		"++",
		func(a AST, env *Environment) (AST, error) {
			switch A := a.(type) {
			case List:
				return Builtin{
					"++ curried-list",
					func(b AST, env *Environment) (AST, error) {
						switch B := b.(type) {
						case List:
							return List{append(append([]AST{}, A.Values...), B.Values...)}, nil
						}

						return nil, NewRuntimeError(nil, "Can't concatenate non-list type")
					},
				}, nil

			case String:
				return Builtin{
					"++ curried-string",
					func(b AST, env *Environment) (AST, error) {
						switch B := b.(type) {
						case String:
							return String{A.Value + B.Value}, nil
						}

						return nil, NewRuntimeError(nil, "Can't concatenate non-list type")
					},
				}, nil

			}

			return nil, NewRuntimeError(nil, "Can't concatenate non-list type")
		},
	},

	{
		"print",
		func(a AST, env *Environment) (AST, error) {
			Print(a)

			return a, nil
		},
	},

	{
		"len",
		func(a AST, env *Environment) (AST, error) {
			switch A := a.(type) {
			case String:
				return Number{len(A.Value)}, nil
			}

			return nil, NewRuntimeError(nil, "Can't find the length of non-string type")
		},
	},
}
