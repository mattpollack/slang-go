package bytecode

import (
	"../ast"
)

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

func pass0(e ast.Expression) ast.Expression {
	// NOTE: small hack since all other interface methods are copy operations
	e = e.MetaSet(
		Meta{
			name: "__NO_NAME",
		},
	).(ast.Expression)

	switch E := e.(type) {
	case ast.Application:
		for i, _ := range E.Body {
			E.Body[i] = pass0(E.Body[i])
		}

	case ast.Pattern:
		// Closure conversion
		// Move all nested closures to a top level let block with auto generated names

		// Replace pattern matches which bind values with stack indices
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
