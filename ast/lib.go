package ast

var libFns = []Builtin{
	{
		"+",
		func(a AST, env *Environment) (AST, *RuntimeError) {
			switch A := a.(type) {
			case Number:
				return Builtin{
					"+ curried",
					func(b AST, env *Environment) (AST, *RuntimeError) {
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
		"-",
		func(a AST, env *Environment) (AST, *RuntimeError) {
			switch A := a.(type) {
			case Number:
				return Builtin{
					"- curried",
					func(b AST, env *Environment) (AST, *RuntimeError) {
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
		"++",
		func(a AST, env *Environment) (AST, *RuntimeError) {
			switch A := a.(type) {
			case List:
				return Builtin{
					"++ curried",
					func(b AST, env *Environment) (AST, *RuntimeError) {
						switch B := b.(type) {
						case List:
							return List{append(append([]AST{}, A.Values...), B.Values...)}, nil
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
		func(a AST, env *Environment) (AST, *RuntimeError) {
			Print(a)

			return a, nil
		},
	},
}

var StdLib = NewEnv(nil)

func init() {
	for _, builtin := range libFns {
		StdLib.Set(builtin.name, builtin)
	}
}
