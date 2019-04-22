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

func (A BoundValue) EqualsExpr(b ast.Expression) bool {
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

// --------------------------------------------------------

func ReplaceMatch(in ast.Match, replace ast.Expression, with ast.Expression) ast.Match {
	switch In := in.(type) {
	case ast.Where:
		In.Condition = Replace(In.Condition, replace, with)
	}

	return in
}

func Replace(in ast.Expression, replace ast.Expression, with ast.Expression) ast.Expression {
	if in.EqualsExpr(replace) {
		return with
	}

	switch In := in.(type) {
	case ast.Application:
		for i, _ := range In.Body {
			In.Body[i] = Replace(In.Body[i], replace, with)
		}

	case ast.Pattern:
		for _, group := range In.Matches {
			for j, _ := range group {
				group[j] = ReplaceMatch(group[j], replace, with)
			}
		}

		for i, _ := range In.Bodies {
			In.Bodies[i] = Replace(In.Bodies[i], replace, with)
		}

	case ast.Let:
		for i, _ := range In.BoundValues {
			// If the replace is an identifier, and the let binding reassigns it, break
			if In.BoundIds[i].EqualsExpr(replace) {
				break
			}

			In.BoundValues[i] = Replace(In.BoundValues[i], replace, with)
		}

	case ast.If:
		In.Condition = Replace(In.Condition, replace, with)
		In.Tbody = Replace(In.Tbody, replace, with)
		In.Fbody = Replace(In.Fbody, replace, with)

	}

	return in
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
						ast.Identifier{Value: "n"},
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
						ast.Identifier{Value: "env"},
						ast.Label{Value: fmt.Sprintf("arg_%d", j)},
					}},
				)
			}

			// Add new argument value to environment
			env.Matches = append(env.Matches,
				[]ast.Match{ast.Label{Value: fmt.Sprintf("arg_%d", i)}},
			)
			env.Bodies = append(env.Bodies, ast.Identifier{Value: "n"})

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
									ast.Identifier{Value: "env"},
									ast.Label{Value: fmt.Sprintf("match_%d", j)},
								},
							},
							ast.Application{
								Body: []ast.Expression{
									ast.Identifier{Value: "=="},
									ast.Identifier{Value: "n"},
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
							ast.Identifier{Value: "env"},
							ast.Label{Value: fmt.Sprintf("match_%d", i)},
						},
					},
				},
			})

			// Replace bound identifiers in pattern bodies with arguments
			for j, match := range E.Matches[i] {
				with := ast.Application{
					Body: []ast.Expression{
						ast.Identifier{Value: "env"},
						ast.Label{Value: fmt.Sprintf("arg_%d", j)},
					},
				}

				switch M := match.(type) {
				case ast.Where:
					body = Replace(body, M.Id, with)
					M.Condition = Replace(M.Condition, M.Id, with)

				case ast.Identifier:
					body = Replace(body, M, with)
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

		for i, _ := range E.Bodies {
			res.BoundIds = append(res.BoundIds, ast.Identifier{Value: fmt.Sprintf("%s_%d", name, i)})
		}

		return res

	case ast.Let:
		for i, _ := range E.BoundValues {
			E.BoundValues[i] = pass0(E.BoundValues[i])
		}

	case ast.If:
		E.Condition = pass0(E.Condition)
		E.Tbody = pass0(E.Tbody)
		E.Fbody = pass0(E.Fbody)
	}

	return e
}

/*
func pass0(e ast.Expression) ast.Expression {
	// Set default meta information
	e = e.MetaSet(
		&Meta{
			name: "__NO_NAME",
		},
	).(ast.Expression)

	switch E := e.(type) {
	case ast.Application:
		for i, _ := range E.Body {
			E.Body[i] = pass0(E.Body[i])
		}

	case ast.Pattern:
		for i, group := range E.Matches {
			j := 0

			for k, m := range group {
				bv := NewBoundValue(j)

				switch M := m.(type) {
				case ast.Where:
					E.Bodies[i] = Replace(E.Bodies[i], M.Id, bv)
					M.Condition = Replace(M.Condition, M.Id, bv)
					M.Id = bv
					group[k] = M
					j++
				case ast.Identifier:
					E.Bodies[i] = Replace(E.Bodies[i], M, NewBoundValue(j))
					j++
				}
			}
		}

		for i, _ := range E.Bodies {
			E.Bodies[i] = pass0(E.Bodies[i])
		}

	case ast.Let:
		for i, _ := range E.BoundValues {
			E.BoundValues[i] = pass0(E.BoundValues[i])
		}

	case ast.If:
		E.Condition = pass0(E.Condition)
		E.Tbody = pass0(E.Tbody)
		E.Fbody = pass0(E.Fbody)
	}

	return e
}
*/
