package ast

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const (
	_ = iota

	// Reserved sorted by length
	TOKEN_KIND_MATCH
	TOKEN_KIND_ELSE

	TOKEN_KIND_IF
	TOKEN_KIND_FAT_ARROW
	TOKEN_KIND_ARROW
	TOKEN_KIND_OP_EQUALS
	TOKEN_KIND_OP_NOT_EQUALS
	TOKEN_KIND_OP_GREATER_THAN_EQUAL
	TOKEN_KIND_OP_LESS_THAN_EQUAL
	TOKEN_KIND_OP_LOGICAL_AND
	TOKEN_KIND_OP_LOGICAL_OR
	TOKEN_KIND_OP_APPEND

	TOKEN_KIND_BRACE_OPEN
	TOKEN_KIND_BRACE_CLOSE
	TOKEN_KIND_PAREN_OPEN
	TOKEN_KIND_PAREN_CLOSE
	TOKEN_KIND_BRACKET_OPEN
	TOKEN_KIND_BRACKET_CLOSE
	TOKEN_KIND_EQUAL
	TOKEN_KIND_COLON
	TOKEN_KIND_COMMA
	TOKEN_KIND_OP_ADD
	TOKEN_KIND_OP_SUBTRACT
	TOKEN_KIND_OP_MULTIPLY
	TOKEN_KIND_OP_DIVIDE
	TOKEN_KIND_OP_GREATER_THAN
	TOKEN_KIND_OP_LESS_THAN

	TOKEN_KIND_SEMI_COLON

	TOKEN_KIND_IDENTIFIER
	TOKEN_KIND_LABEL
	TOKEN_KIND_STRING
	TOKEN_KIND_NUMBER

	TOKEN_KIND_ERROR
	TOKEN_KIND_EOF
)

var opPrecedence [][]int = [][]int{
	[]int{
		TOKEN_KIND_OP_ADD,
		TOKEN_KIND_OP_SUBTRACT,
	},
	[]int{
		TOKEN_KIND_OP_MULTIPLY,
		TOKEN_KIND_OP_DIVIDE,
	},
	[]int{
		TOKEN_KIND_OP_APPEND,
	},
	[]int{
		TOKEN_KIND_OP_LOGICAL_AND,
		TOKEN_KIND_OP_LOGICAL_OR,
	},
	[]int{
		TOKEN_KIND_OP_GREATER_THAN,
		TOKEN_KIND_OP_GREATER_THAN_EQUAL,
		TOKEN_KIND_OP_LESS_THAN,
		TOKEN_KIND_OP_LESS_THAN_EQUAL,
		TOKEN_KIND_OP_EQUALS,
		TOKEN_KIND_OP_NOT_EQUALS,
	},
}

type token struct {
	kind  int
	value []byte

	line int
	char int

	firstOfLine  bool // this token is first on its line
	skippedSpace bool // this token did not skip spaces before getting parsed
}

type parser struct {
	src []byte

	line int
	char int
}

type mark struct {
	b parser
	a *parser
}

func NewMark(p *parser) *mark {
	return &mark{
		*p,
		p,
	}
}

func (m *mark) Done() parser {
	p := parser{
		m.b.src[:len(m.b.src)-len(m.a.src)],
		m.b.line - m.a.line,
		m.b.char - m.a.char,
	}
	p.SkipWS()

	return p
}

func (m *mark) Error(msg string) error {
	msg = fmt.Sprintf("[%d:%d] %s", m.a.line, m.a.char, msg)
	src := fmt.Sprintf("%s\n", string(m.Done().src))

	for i := 1; i < m.a.char; i++ {
		src += " "
	}

	src += "^"

	msg = fmt.Sprintf("%s\n%s", src, msg)

	return errors.New(msg)
}

func (p *parser) SkipWS() {
	i := 0

	for ; i < len(p.src) && unicode.IsSpace(rune(p.src[i])); i++ {
		if p.src[i] == '\n' {
			p.line++
			p.char = 1
		} else {
			p.char++
		}
	}

	p.src = p.src[i:]
}

func (p *parser) Next() token {
	oldLen := len(p.src)
	oldLine := p.line

	// Handle spaces and comments
	for {
		p.SkipWS()

		if len(p.src) > 0 && p.src[0] == '#' {
			i := 0

			for ; i < len(p.src) && p.src[i] != '\n'; i++ {
			}

			p.src = p.src[i:]
		} else {
			break
		}
	}

	firstOfLine := oldLine != p.line
	skippedSpace := oldLen != len(p.src)

	if len(p.src) == 0 {
		return token{
			TOKEN_KIND_EOF,
			[]byte("EOF"),
			p.line,
			p.char,
			firstOfLine,
			skippedSpace,
		}
	}

	// Parse reserved
	for i, s := range []string{
		"match",
		"else",
		"if",
		"=>",
		"->",
		"==",
		"!=",
		">=",
		"<=",
		"&&",
		"||",
		"++",
		"{", "}",
		"(", ")",
		"[", "]",
		"=",
		":",
		",",
		"+",
		"-",
		"*",
		"/",
		">",
		"<",
		";",
	} {
		if bytes.HasPrefix(p.src, []byte(s)) {
			token := token{
				TOKEN_KIND_MATCH + i,
				[]byte(s),
				p.line,
				p.char,
				firstOfLine,
				skippedSpace,
			}

			p.char += len(s)
			p.src = p.src[len(s):]

			return token
		}
	}

	// Parse identifier
	if unicode.IsLetter(rune(p.src[0])) || p.src[0] == '_' {
		i := 0

		for ; i < len(p.src) && (unicode.IsLetter(rune(p.src[i])) || unicode.IsNumber(rune(p.src[i])) || p.src[i] == '_'); i++ {
		}

		token := token{
			TOKEN_KIND_IDENTIFIER,
			p.src[:i],
			p.line,
			p.char,
			firstOfLine,
			skippedSpace,
		}

		p.char += i
		p.src = p.src[i:]

		return token
	}

	// Parse label
	if p.src[0] == '.' {
		i := 1

		for ; i < len(p.src) && (unicode.IsLetter(rune(p.src[i])) || unicode.IsNumber(rune(p.src[i])) || p.src[i] == '_'); i++ {
		}

		token := token{
			TOKEN_KIND_LABEL,
			p.src[1:i],
			p.line,
			p.char,
			firstOfLine,
			skippedSpace,
		}

		p.char += i
		p.src = p.src[i:]

		return token
	}

	if p.src[0] == '"' {
		i := 1

		for ; i < len(p.src) && p.src[i] != '"'; i++ {
		}

		i++

		token := token{
			TOKEN_KIND_STRING,
			p.src[1 : i-1],
			p.line,
			p.char,
			firstOfLine,
			skippedSpace,
		}

		p.char += i
		p.src = p.src[i:]

		return token
	}

	// Parse number
	if unicode.IsNumber(rune(p.src[0])) {
		i := 0

		for ; i < len(p.src) && unicode.IsNumber(rune(p.src[i])); i++ {
		}

		token := token{
			TOKEN_KIND_NUMBER,
			p.src[:i],
			p.line,
			p.char,
			firstOfLine,
			skippedSpace,
		}

		p.char += i
		p.src = p.src[i:]

		return token
	}

	return token{
		TOKEN_KIND_ERROR,
		[]byte("Unexpected end of parsing"),
		p.line,
		p.char,
		firstOfLine,
		skippedSpace,
	}
}

func (p *parser) Peek() token {
	src := p.src
	line := p.line
	char := p.char

	token := p.Next()

	p.src = src
	p.line = line
	p.char = char

	return token
}

func (p *parser) ConsumeIfNext(kind int) bool {
	if p.Peek().kind == kind {
		p.Next()
		return true
	}

	return false
}

func (p *parser) Identifier() (Expression, error) {
	if p.Peek().kind == TOKEN_KIND_IDENTIFIER {
		return NewIdentifier(string(p.Next().value))
	}

	return nil, errors.New("Cannot parse identifier")
}

func (p *parser) Label() (Expression, error) {
	if p.Peek().kind == TOKEN_KIND_LABEL {
		return NewLabel(string(p.Next().value))
	}

	return nil, errors.New("Cannot parse label")
}

func (p *parser) String() (Expression, error) {
	if p.Peek().kind == TOKEN_KIND_STRING {
		str := string(p.Next().value)

		for i := 0; i < len(str); i++ {
			switch {
			case strings.HasPrefix(str[i:], "\\n"):
				str = str[:i] + "\n" + str[i+2:]
			case strings.HasPrefix(str[i:], "\\t"):
				str = str[:i] + "\t" + str[i+2:]
			}
		}

		return NewString(str)
	}

	return nil, errors.New("Cannot parse string")
}

func (p *parser) Number() (Expression, error) {
	if p.Peek().kind != TOKEN_KIND_NUMBER {
		return nil, errors.New("Cannot parse number")
	}

	str := p.Next().value
	value := 0

	for i := 0; i < len(str); i++ {
		value = value*10 + int(str[i]-'0')
	}

	return NewNumber(value)
}

/*
func (p *parser) If() (Expression, error) {
	m := NewMark(p)

	if !p.ConsumeIfNext(TOKEN_KIND_IF) {
		return nil, m.Error("If must begin with 'if'")
	}

	condition, err := p.Expression()

	if err != nil {
		return nil, m.Error("If must have a condition")
	}

	tbody, err := p.Expression()

	if err != nil {
		return nil, m.Error("If must have a body")
	}

	if !p.ConsumeIfNext(TOKEN_KIND_ELSE) {
		return nil, m.Error("If must be followed by an else")
	}

	fbody, err := p.Expression()

	if err != nil {
		return nil, m.Error("Else must have a body")
	}

	return NewIf(condition, tbody, fbody)
}
*/

func (p *parser) MatchExpr() (Expression, error) {
	m := NewMark(p)

	if !p.ConsumeIfNext(TOKEN_KIND_MATCH) {
		return nil, m.Error("Match must begin with 'match'")
	}

	toMatch, err := p.Expression([]int{TOKEN_KIND_BRACE_OPEN})

	if err != nil {
		return nil, err
	}

	with, err := p.Pattern()

	if err != nil {
		return nil, err
	}

	return NewApplication([]Expression{with, toMatch})
}

func (p *parser) Let(identifier Identifier) (Expression, error) {
	m := NewMark(p)

	if !p.ConsumeIfNext(TOKEN_KIND_EQUAL) {
		return nil, m.Error("'=' must follow identifier in let")
	}

	value, err := p.Expression([]int{})

	if err != nil {
		return nil, m.Error("Let must be assigned a value")
	}

	body, err := p.Expression([]int{})

	if err != nil {
		return nil, m.Error("Let must have a body")
	}

	// If the body is a let expression, merge them
	switch b := body.(type) {
	case Let:
		return NewLet(
			append([]Identifier{identifier}, b.BoundIds...),
			append([]Expression{value}, b.BoundValues...),
			b.Body,
		)
	default:
		return NewLet(
			[]Identifier{identifier},
			[]Expression{value},
			body,
		)
	}
}

func (p *parser) List(first Expression) (Expression, error) {
	list := List{}

	// Because lists and list constructors share their starts
	if first != nil {
		list.Values = append(list.Values, first)
	}

	// Check for an empty list
	if !p.ConsumeIfNext(TOKEN_KIND_BRACKET_CLOSE) {
		for {
			val, err := p.Expression([]int{TOKEN_KIND_COMMA, TOKEN_KIND_BRACKET_CLOSE})

			if err != nil {
				return nil, err
			}

			list.Values = append(list.Values, val)

			if p.ConsumeIfNext(TOKEN_KIND_BRACKET_CLOSE) {
				break
			}

			// NOTE: !p.Peek().firstOfLine &&
			if !p.ConsumeIfNext(TOKEN_KIND_COMMA) {
				return nil, errors.New("Values in a list must be separated by commas or new lines")
			}
		}
	}

	return list, nil
}

func (p *parser) ListConstructor(head Expression) (Match, error) {
	if !p.ConsumeIfNext(TOKEN_KIND_COLON) {
		return nil, errors.New("List constructors must have a colon ':' separating the head and tail expressions")
	}

	tail, err := p.Expression([]int{TOKEN_KIND_BRACKET_CLOSE})

	if err != nil {
		return nil, err
	}

	if !p.ConsumeIfNext(TOKEN_KIND_BRACKET_CLOSE) {
		return nil, errors.New("List constructors must be enclosed by brackets '[' ']'")
	}

	return ListConstructor{Head: head, Tail: tail}, nil
}

func (p *parser) Slice() (Expression, error) {
	m := NewMark(p)
	slice := Slice{}

	if !p.ConsumeIfNext(TOKEN_KIND_BRACKET_OPEN) {
		return nil, m.Error("Slice must be enclosed by brackets '[' ']'")
	}

	if p.Peek().kind != TOKEN_KIND_COLON {
		low, err := p.Expression([]int{TOKEN_KIND_COLON})

		if err != nil {
			return nil, err
		}

		slice.Low = low
	}

	if !p.ConsumeIfNext(TOKEN_KIND_COLON) {
		return nil, m.Error("Slice bounds must be separated by a colon ':'")
	}

	if p.Peek().kind != TOKEN_KIND_BRACKET_CLOSE {
		high, err := p.Expression([]int{TOKEN_KIND_BRACKET_CLOSE})

		if err != nil {
			return nil, err
		}

		slice.High = high
	}

	if !p.ConsumeIfNext(TOKEN_KIND_BRACKET_CLOSE) {
		return nil, m.Error("Slice must be enclosed by brackets '[' ']'")
	}

	return slice, nil
}

func (p *parser) Where(id Identifier) (Match, error) {
	m := NewMark(p)

	if !p.ConsumeIfNext(TOKEN_KIND_COLON) {
		return nil, m.Error("Colon must separate where identifier from body")
	}

	if !p.ConsumeIfNext(TOKEN_KIND_PAREN_OPEN) {
		return nil, m.Error("Where body must be enclosed by parenthesis '(' ')'")
	}

	body, err := p.Expression([]int{TOKEN_KIND_PAREN_CLOSE})

	if err != nil {
		return nil, err
	}

	if !p.ConsumeIfNext(TOKEN_KIND_PAREN_CLOSE) {
		return nil, m.Error("Where body must be enclosed by parenthesis '(' ')'")
	}

	return NewWhere(id, body)
}

func (p *parser) Pattern() (Expression, error) {
	m := NewMark(p)

	if !p.ConsumeIfNext(TOKEN_KIND_BRACE_OPEN) {
		return nil, m.Error("Pattern must be enclosed by braces '{' '}'")
	}

	matchBodies := [][]Match{}
	bodies := []Expression{}

	for !p.ConsumeIfNext(TOKEN_KIND_BRACE_CLOSE) {
		matches := []Match{}
		isImplicitBody := false

		if p.ConsumeIfNext(TOKEN_KIND_FAT_ARROW) {
			if len(matchBodies) == 0 {
				return nil, errors.New("Default pattern match cannot be the first body")
			}

			for range matchBodies[0] {
				id, _ := NewIdentifier("_")
				matches = append(matches, id)
			}
		} else {
			for !p.ConsumeIfNext(TOKEN_KIND_ARROW) {
				match, err := p.Match()

				if err != nil {
					return nil, err
				}

				matches = append(matches, match)

				if p.Peek().kind == TOKEN_KIND_BRACE_CLOSE {
					isImplicitBody = true
					break
				}
			}
		}

		// Add missing values for implicit body
		if isImplicitBody {
			if len(matchBodies) != 0 {
				return nil, errors.New("Pattern can only have one implicit true match")
			}

			falseMatches := []Match{}

			for range matches {
				id, _ := NewIdentifier("_")
				falseMatches = append(falseMatches, id)
			}

			bodies = append(bodies, True)
			bodies = append(bodies, False)
			matchBodies = append(matchBodies, matches)
			matchBodies = append(matchBodies, falseMatches)

		} else {
			// TODO: revisit cases
			body, err := p.Expression([]int{TOKEN_KIND_BRACE_CLOSE})

			if err != nil {
				return nil, err
			}

			if len(matchBodies) > 0 && len(matches) != len(matchBodies[0]) {
				return nil, errors.New("Pattern cannot take varying arguments")
			}

			matchBodies = append(matchBodies, matches)
			bodies = append(bodies, body)
		}
	}

	return NewPattern(matchBodies, bodies)
}

func (p *parser) Match() (Match, error) {
	m := NewMark(p)

	switch p.Peek().kind {
	// Could be where or identifier
	case TOKEN_KIND_IDENTIFIER:
		m, err := p.Identifier()

		if err != nil {
			return nil, err
		}

		if p.Peek().kind == TOKEN_KIND_COLON {
			m, err := p.Where(m.(Identifier))

			if err != nil {
				return nil, err
			}

			return m, nil
		}

		return m.(Identifier), nil

	case TOKEN_KIND_LABEL:
		m, err := p.Label()

		if err != nil {
			return nil, err
		}

		return m.(Label), nil

	case TOKEN_KIND_STRING:
		m, err := p.String()

		if err != nil {
			return nil, err
		}

		return m.(String), nil

	case TOKEN_KIND_NUMBER:
		m, err := p.Number()

		if err != nil {
			return nil, err
		}

		return m.(Number), nil

	// Lists and list constructors: open bracket and an expression
	case TOKEN_KIND_BRACKET_OPEN:
		p.Next()

		if p.Peek().kind == TOKEN_KIND_BRACKET_CLOSE {
			m, err := p.List(nil)

			if err != nil {
				return nil, err
			}

			return m.(List), nil
		}

		if p.Peek().kind == TOKEN_KIND_COLON {
			return nil, errors.New("List constructor cannot have an empty head value")
		}

		expr, err := p.Expression([]int{TOKEN_KIND_COLON, TOKEN_KIND_COMMA, TOKEN_KIND_BRACKET_CLOSE})

		if err != nil {
			return nil, err
		}

		if p.Peek().kind == TOKEN_KIND_COLON {
			return p.ListConstructor(expr)
		} else {
			m, err := p.List(expr)

			if err != nil {
				return nil, err
			}

			return m.(List), nil
		}

	default:
		return nil, m.Error("Unexpected error occured when parsing match")
	}
}

func (p *parser) PrimaryExpr(endTokenKinds []int) (Expression, error) {
	if p.Peek().kind == TOKEN_KIND_IDENTIFIER {
		id, _ := NewIdentifier(string(p.Next().value))

		if p.Peek().kind == TOKEN_KIND_EQUAL {
			return p.Let(id)
		} else {
			return id, nil
		}
	}

	if p.Peek().kind == TOKEN_KIND_LABEL {
		return p.Label()
	}

	if p.Peek().kind == TOKEN_KIND_STRING {
		return p.String()
	}

	if p.Peek().kind == TOKEN_KIND_NUMBER {
		return p.Number()
	}

	if p.Peek().kind == TOKEN_KIND_MATCH {
		return p.MatchExpr()
	}

	if p.Peek().kind == TOKEN_KIND_BRACE_OPEN {
		return p.Pattern()
	}

	if p.Peek().kind == TOKEN_KIND_PAREN_OPEN {
		p.Next()
		res, err := p.Expression([]int{TOKEN_KIND_PAREN_CLOSE})

		if err != nil {
			return nil, err
		}

		if !p.ConsumeIfNext(TOKEN_KIND_PAREN_CLOSE) {
			return nil, errors.New("Application which starts with an open parenthesis must close with one")
		}

		return res, nil
	}

	// Lists and list constructors: open bracket and an expression
	if p.Peek().kind == TOKEN_KIND_BRACKET_OPEN {
		p.Next()
		return p.List(nil)
	}

	return nil, errors.New("Unexpected error occured when parsing an expression")
}

func (p *parser) OpExpr(precedence int, endTokenKinds []int) (Expression, error) {
	if precedence >= len(opPrecedence) {
		exprs := []Expression{}

		for {
			next, err := p.PrimaryExpr(endTokenKinds)

			if err != nil {
				return nil, err
			}

			exprs = append(exprs, next)

			// Parse for end of op
			end := false

			for _, token := range endTokenKinds {
				if p.Peek().kind == token {
					end = true
					break
				}
			}

			if end {
				break
			}

			if p.Peek().firstOfLine {
				break
			}
		}

		if len(exprs) == 0 {
			return nil, errors.New("Cannot parse primary expression")
		}

		if len(exprs) == 1 {
			return exprs[0], nil
		}

		return Application{Body: exprs}, nil
	} else {
		var head Expression
		var err error
		head, err = p.OpExpr(precedence+1, append(append([]int{}, opPrecedence[precedence]...), endTokenKinds...))

		if err != nil {
			return nil, err
		}

		for !p.Peek().firstOfLine {
			// Parse for end of op
			end := false

			for _, token := range endTokenKinds {
				if p.Peek().kind == token {
					end = true
					break
				}
			}

			if end {
				break
			}

			// Parse op
			foundOp := false

			for _, token := range opPrecedence[precedence] {
				if p.Peek().kind == token {
					foundOp = true
					t := p.Next()

					next, err := p.OpExpr(precedence+1, append(append([]int{}, opPrecedence[precedence]...), endTokenKinds...))

					if err != nil {
						return nil, err
					}

					head = Application{Body: []Expression{Identifier{Value: string(t.value)}, head, next}}
					break
				}
			}

			if !foundOp {
				return nil, errors.New("Unexpected token found when attempting to parse operator")
			}
		}

		return head, nil
	}
}

func (p *parser) Expression(endTokenKinds []int) (Expression, error) {
	// Start parsing with default ender tokens
	expr, err := p.OpExpr(0, append(endTokenKinds, []int{TOKEN_KIND_EOF, TOKEN_KIND_SEMI_COLON}...))

	if err != nil {
		return nil, err
	}

	if p.Peek().kind == TOKEN_KIND_SEMI_COLON {
		p.Next()
	}

	return expr, nil
}

func Parse(src []byte) (Expression, error) {
	p := parser{
		src,
		1,
		1,
	}

	ast, err := p.Expression([]int{})

	if err != nil {
		return nil, err
	}

	if p.Next().kind != TOKEN_KIND_EOF {
		return nil, fmt.Errorf("[%d:%d] Unexpected end of parsing", p.line, p.char)
	}

	return ast, nil
}
