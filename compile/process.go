package compile

import (
	"fmt"

	"../ast"
)

// --------------------------------------------------------

type BoundValue struct {
	Index int

	meta interface{}
}

func (e BoundValue) IsExpression() {}

func NewBoundValue(s int) BoundValue {
	return BoundValue{
		s,
		nil,
	}
}

func (A BoundValue) Equals(b interface{}) bool {
	switch B := b.(type) {
	case BoundValue:
		return A.Index == B.Index
	}

	return false
}

func (e BoundValue) MetaGet() interface{} {
	return e.meta
}

func (e BoundValue) MetaSet(meta interface{}) interface{} {
	e.meta = meta

	return e
}

func (e BoundValue) Eval(ast.Environment) ast.Expression {
	return e
}

func (e BoundValue) Apply(ast.Expression) ast.Expression {
	return e
}

func (e BoundValue) Print(tab int) {
	ast.PrintTab(tab)
	fmt.Printf("arg[%d]\n", e.Index)
}

// --------------------------------------------------------

type Function struct {
	// entry
	// bodies
	// final body
}

// --------------------------------------------------------

func ReplaceMatch(in ast.Match, replace ast.Expression, with ast.Expression) (ast.Match, bool) {
	var contains bool
	var tmp bool

	switch In := in.(type) {
	case ast.Where:
		In.Condition, tmp = Replace(In.Condition, replace, with)
		contains = contains || tmp
	}

	return in, contains
}

// Returns the result expression and whether it replaced anything
func Replace(in ast.Expression, replace ast.Expression, with ast.Expression) (ast.Expression, bool) {
	if in.Equals(replace) {
		return with, true
	}

	var contains bool
	var tmp bool

	switch In := in.(type) {
	case ast.Application:
		for i, _ := range In.Body {
			In.Body[i], tmp = Replace(In.Body[i], replace, with)
			contains = contains || tmp
		}

	case ast.Pattern:
		for _, group := range In.Matches {
			for j, _ := range group {
				group[j], tmp = ReplaceMatch(group[j], replace, with)
				contains = contains || tmp
			}
		}

		for i, _ := range In.Bodies {
			In.Bodies[i], tmp = Replace(In.Bodies[i], replace, with)
			contains = contains || tmp
		}

	case ast.Let:
		for i, _ := range In.BoundValues {
			// If the replace is an identifier, and the let binding reassigns it, break
			if In.BoundIds[i].Equals(replace) {
				break
			}

			In.BoundValues[i], tmp = Replace(In.BoundValues[i], replace, with)
			contains = contains || tmp
		}

	case ast.If:
		In.Condition, tmp = Replace(In.Condition, replace, with)
		contains = contains || tmp

		In.Tbody, tmp = Replace(In.Tbody, replace, with)
		contains = contains || tmp

		In.Fbody, tmp = Replace(In.Fbody, replace, with)
		contains = contains || tmp

	}

	return in, contains
}

// --------------------------------------------------------

func defaultMeta() interface{} {
	return &Meta{
		name: "__NO_NAME",
	}
}

func buildEnv(matches int) ast.Pattern {
	env := ast.Pattern{}

	for i := 0; i < matches; i++ {
		env.Matches = append(env.Matches, []ast.Match{
			ast.Label{Value: fmt.Sprintf("match_%d", i)},
		})
		env.Bodies = append(env.Bodies, ast.Label{Value: "true"})
	}

	return env
}

var ID = 0

func nextId() string {
	ID++

	return fmt.Sprintf("__gen__%d", ID)
}

func pass0(e ast.Expression) ast.Expression {
	switch E := e.(type) {
	case ast.Application:
		for i, _ := range E.Body {
			E.Body[i] = pass0(E.Body[i])
		}

		return E

	case ast.Pattern:
		// Visit bodies
		for i, _ := range E.Bodies {
			E.Bodies[i] = pass0(E.Bodies[i])
		}

		// 1) Generate entry function
		name := nextId()

		entry := ast.Pattern{
			Matches: [][]ast.Match{[]ast.Match{ast.Identifier{Value: "n"}}},
			Bodies: []ast.Expression{
				ast.Application{
					Body: []ast.Expression{
						ast.Identifier{Value: fmt.Sprintf("%s_0", name)},
						buildEnv(len(E.Bodies)),
						NewBoundValue(0),
					},
				},
			},
		}

		// 2) Generate body functions
		bodies := []ast.Pattern{}

		for i := 0; i < len(E.Matches[0]); i++ {
			// Create pattern which accepts an environment and a param
			match := ast.Pattern{
				Matches: [][]ast.Match{[]ast.Match{
					ast.Identifier{Value: "n"},
					ast.Identifier{Value: "env"},
				}},
				Bodies: []ast.Expression{},
			}

			env := ast.Pattern{
				Matches: [][]ast.Match{[]ast.Match{ast.Label{Value: "next"}}},
				Bodies:  []ast.Expression{ast.Identifier{Value: fmt.Sprintf("%s_%d", name, i+1)}},
			}

			// Add already captured values to environment
			for j := 0; j < i; j++ {
				env.Matches = append(env.Matches,
					[]ast.Match{ast.Label{Value: fmt.Sprintf("arg_%d", j)}},
				)
				env.Bodies = append(env.Bodies,
					ast.Application{Body: []ast.Expression{
						NewBoundValue(1),
						ast.Label{Value: fmt.Sprintf("arg_%d", j)},
					}},
				)
			}

			// Add new argument value to environment
			env.Matches = append(env.Matches,
				[]ast.Match{ast.Label{Value: fmt.Sprintf("arg_%d", i)}},
			)
			env.Bodies = append(env.Bodies, NewBoundValue(0))

			// Add match state to environment
			for j, matchGroup := range E.Matches {
				env.Matches = append(env.Matches,
					[]ast.Match{ast.Label{Value: fmt.Sprintf("match_%d", j)}},
				)

				// TODO: move match generation to its own function?
				// Generate match expression for the i'th match group
				switch matchGroup[i].(type) {
				case ast.Identifier:
					env.Bodies = append(env.Bodies, ast.Label{Value: "true"})
				case ast.Where:
					panic("TODO: generate match for where")
				default:
					env.Bodies = append(env.Bodies, ast.Application{
						Body: []ast.Expression{
							ast.Identifier{Value: "&&"},
							ast.Application{
								Body: []ast.Expression{
									NewBoundValue(1),
									ast.Label{Value: fmt.Sprintf("match_%d", j)},
								},
							},
							ast.Application{
								Body: []ast.Expression{
									ast.Identifier{Value: "=="},
									NewBoundValue(0),
									matchGroup[i].(ast.Expression),
								},
							},
						},
					})
				}
			}

			// Add sub expressions to parent expressions
			match.Bodies = append(match.Bodies, env)
			bodies = append(bodies, match)
		}

		// Alter the last match to call the final body
		bodies[len(bodies)-1].Bodies = []ast.Expression{
			ast.Application{
				Body: []ast.Expression{
					ast.Identifier{Value: fmt.Sprintf("%s_%d", name, len(E.Matches[0]))},
					bodies[len(bodies)-1].Bodies[0],
				},
			},
		}

		// 3) Generate final body
		finalBody := ast.Pattern{}

		for i, body := range E.Bodies {
			finalBody.Matches = append(finalBody.Matches, []ast.Match{
				ast.Where{
					Id: ast.Identifier{Value: "env"},
					Condition: ast.Application{
						Body: []ast.Expression{
							NewBoundValue(0),
							ast.Label{Value: fmt.Sprintf("match_%d", i)},
						},
					},
				},
			})

			// Replace bound identifiers in pattern bodies with arguments
			for j, match := range E.Matches[i] {
				with := ast.Application{
					Body: []ast.Expression{
						NewBoundValue(0),
						ast.Label{Value: fmt.Sprintf("arg_%d", j)},
					},
				}

				switch M := match.(type) {
				case ast.Where:
					body, _ = Replace(body, M.Id, with)
					M.Condition, _ = Replace(M.Condition, M.Id, with)

				case ast.Identifier:
					body, _ = Replace(body, M, with)
				}
			}

			finalBody.Bodies = append(finalBody.Bodies, body)
		}

		// 4) Build let to attach all the new pattern components
		res := ast.Let{}

		for _, b := range bodies {
			res.BoundValues = append(res.BoundValues, b)
		}

		res.BoundValues = append(res.BoundValues, finalBody)
		res.Body = entry

		for i := 0; i <= len(E.Matches[0]); i++ {
			res.BoundIds = append(res.BoundIds, ast.Identifier{Value: fmt.Sprintf("%s_%d", name, i)})
		}

		return res

	case ast.Let:
		for i, _ := range E.BoundValues {
			E.BoundValues[i] = pass0(E.BoundValues[i])

			/* NOTE: maybe replace ids with parent id?????
			val := pass0(E.BoundValues[i])

			switch V := val.(type) {
			case ast.Let:
				E.BoundIds =
					append(append(append(
						[]ast.Identifier{},
						E.BoundIds[:i]...),
						V.BoundIds...),
						E.BoundIds[i+1:]...)
				E.BoundValues =
					append(append(append(
						[]ast.Expression{},
						E.BoundValues[:i]...),
						V.BoundValues...),
						E.BoundValues[i+1:]...)
			default:
				E.BoundValues[i] = val
			}
			*/
		}

		return E

	case ast.If:
		E.Condition = pass0(E.Condition)
		E.Tbody = pass0(E.Tbody)
		E.Fbody = pass0(E.Fbody)

		return E
	}

	return e
}
