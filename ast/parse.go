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
	TOKEN_KIND_PACKAGE
	TOKEN_KIND_IMPORT
	TOKEN_KIND_MODULE
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
	TOKEN_KIND_DOUBLE_COLON

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
	TOKEN_KIND_OP_MODULO
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
		TOKEN_KIND_OP_MODULO,
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

func pad(len int) string {
	res := ""

	for ;len > 0; len-- {
		res += " "
	}

	return res
}

type ParseError struct {
	wrapped error
	message string
}

/*

Formatted like:

   | package sketch
   |
 3 | test = {
          ^
   ERROR: msg

*/
func NewParseError(p *parser, wrapped error, err string) *ParseError {
	lineNum := fmt.Sprintf("%d", p.line)
	message := "\n"
	padStr := fmt.Sprintf(pad(len(lineNum)+2))

	for i := numberOfSavedLines-1; i >= 0; i-- {
		if i == 0 {
			message += fmt.Sprintf(" %s ", lineNum)
		} else {
			message += padStr
		}

		message += fmt.Sprintf("| %s\n", string(p.savedLines[i]))
	}

	message += pad(p.char+1) + "^"
	message += "\n"+padStr
	message += fmt.Sprintf("ERROR: %s\n", err)

	return &ParseError{
		wrapped,
		message,
	}
}

func (e *ParseError) Error() string {
	if e.wrapped == nil {
		return e.message
	}

	return e.wrapped.Error()
}

func (e *ParseError) Unwrap() error {
	return e.wrapped
}

const numberOfSavedLines = 4

type parser struct {
	src []byte

	line int
	char int

	savedLines [numberOfSavedLines][]byte
	currLine []byte
}

func (p *parser) nextLine()  {
	if len(p.src) == 0 {
		p.currLine = []byte{}

		return
	}

	i := 0

	for ; i < len(p.src) && p.src[i] != '\n'; i++ {}

	var line []byte

	if i < len(p.src) {
		i++
		line = p.src[:i-1]

		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
	} else {
		line = p.src[:i]
	}

	for j := numberOfSavedLines-1; j >= 1; j-- {
		p.savedLines[j] = p.savedLines[j-1]
	}

	p.savedLines[0] = line
	p.src = p.src[i:]
	p.line++
	p.char = 1
	p.currLine = line
}

func (p *parser) currLineIsEmpty() bool {
	return len(p.currLine) == 0
}

func (p *parser) isEOF() bool {
	return len(p.src) == 0 && len(p.currLine) == 0
}

func (p *parser) SkipWS() {
	i := 0

	for ; i < len(p.currLine) && (unicode.IsSpace(rune(p.currLine[i])) || p.currLine[i] == '\r'); i++ {}

	p.char += i
	p.currLine = p.currLine[i:]
}

func (p *parser) Next() token {
	oldLen := len(p.src)
	oldLine := p.line

	for !p.isEOF() {
		p.SkipWS()

		if !p.currLineIsEmpty() && p.currLine[0] != '#' {
			break
		}

		p.nextLine()
	}

	firstOfLine := oldLine != p.line
	skippedSpace := oldLen != len(p.src)

	if p.isEOF() {
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
		"package",
		"import",
		"module",
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
		"::",
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
		"%",
		">",
		"<",
		";",
	} {
		if bytes.HasPrefix(p.currLine, []byte(s)) {
			token := token{
				i + 1,
				[]byte(s),
				p.line,
				p.char,
				firstOfLine,
				skippedSpace,
			}

			p.char += len(s)
			p.currLine = p.currLine[len(s):]

			return token
		}
	}

	// Parse identifier
	if unicode.IsLetter(rune(p.currLine[0])) || p.currLine[0] == '_' {
		i := 0

		for ; i < len(p.currLine) && (unicode.IsLetter(rune(p.currLine[i])) || unicode.IsNumber(rune(p.currLine[i])) || p.currLine[i] == '_'); i++ {
		}

		token := token{
			TOKEN_KIND_IDENTIFIER,
			p.currLine[:i],
			p.line,
			p.char,
			firstOfLine,
			skippedSpace,
		}

		p.char += i
		p.currLine = p.currLine[i:]

		return token
	}

	// Parse label
	if p.currLine[0] == '.' {
		i := 1

		for ; i < len(p.currLine) && (unicode.IsLetter(rune(p.currLine[i])) || unicode.IsNumber(rune(p.currLine[i])) || p.currLine[i] == '_'); i++ {
		}

		token := token{
			TOKEN_KIND_LABEL,
			p.currLine[1:i],
			p.line,
			p.char,
			firstOfLine,
			skippedSpace,
		}

		p.char += i
		p.currLine = p.currLine[i:]

		return token
	}

	// Parse string
	if p.currLine[0] == '"' {
		i := 1

		for ; i < len(p.currLine) && p.currLine[i] != '"'; i++ {
		}

		if i < len(p.currLine) && p.currLine[i] == '"' {
			i++

			token := token{
				TOKEN_KIND_STRING,
				p.currLine[1 : i-1],
				p.line,
				p.char,
				firstOfLine,
				skippedSpace,
			}

			p.char += i
			p.currLine = p.currLine[i:]

			return token
		}
	}

	// Parse number
	if unicode.IsNumber(rune(p.currLine[0])) {
		i := 0

		for ; i < len(p.currLine) && unicode.IsNumber(rune(p.currLine[i])); i++ {
		}

		token := token{
			TOKEN_KIND_NUMBER,
			p.currLine[:i],
			p.line,
			p.char,
			firstOfLine,
			skippedSpace,
		}

		p.char += i
		p.currLine = p.currLine[i:]

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
	savedLines := [numberOfSavedLines][]byte{}

	for i := 0; i < numberOfSavedLines; i++ {
		savedLines[i] = p.savedLines[i]
	}

	currLine := p.currLine

	token := p.Next()

	p.src = src
	p.line = line
	p.char = char
	p.savedLines = savedLines
	p.currLine = currLine

	return token
}

func (p *parser) ConsumeIfNext(kind int) bool {
	if p.Peek().kind == kind {
		p.Next()
		return true
	}

	return false
}

func (p *parser) Identifier() (AST, error) {
	if p.Peek().kind == TOKEN_KIND_IDENTIFIER {
		return NewIdentifier(string(p.Next().value))
	}

	return nil, errors.New("Cannot parse identifier")
}

func (p *parser) Module() (AST, error) {
	if !p.ConsumeIfNext(TOKEN_KIND_MODULE) {
		return nil, NewParseError(p, nil, "Module must begin with 'module'")
	}

	if !p.ConsumeIfNext(TOKEN_KIND_BRACE_OPEN) {
		return nil, NewParseError(p, nil, "Module must be surrounded by braces")
	}

	boundIds := []Identifier{}
	boundValues := []AST{}

	for !p.ConsumeIfNext(TOKEN_KIND_BRACE_CLOSE) {
		id, err := p.Identifier()

		if err != nil {
			return nil, NewParseError(p, err, "Failed to parse bound id in module")
		}

		if !p.ConsumeIfNext(TOKEN_KIND_EQUAL) {
			return nil, NewParseError(p, nil, "Bound id in module must be followed by '='")
		}

		value, err := p.Expression([]int{})

		if err != nil {
			return nil, NewParseError(p, err, "Failed to parse bound value in module")
		}

		boundIds = append(boundIds, id.(Identifier))
		boundValues = append(boundValues, value)
	}

	if len(boundIds) == 0 {
		return nil, NewParseError(p, nil, "Module must have at least one bound value")
	}

	matchGroups := [][]AST{}
	bodies := []AST{}

	for _, id := range boundIds {
		matchGroups = append(matchGroups, []AST{Label{id.Value}})
		bodies = append(bodies, id)
	}

	body, _ := NewPattern(matchGroups, bodies)

	return NewLet(boundIds, boundValues, body)
}

func (p *parser) Label() (AST, error) {
	if p.Peek().kind == TOKEN_KIND_LABEL {
		return NewLabel(string(p.Next().value))
	}

	return nil, errors.New("Cannot parse label")
}

func (p *parser) String() (AST, error) {
	if p.Peek().kind == TOKEN_KIND_STRING {
		str := string(p.Next().value)

		for i := 0; i < len(str); i++ {
			switch {
			case strings.HasPrefix(str[i:], "\\r"):
				str = str[:i] + "\r" + str[i+2:]
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

func (p *parser) Number() (AST, error) {
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

func (p *parser) MatchExpr() (AST, error) {
	if !p.ConsumeIfNext(TOKEN_KIND_MATCH) {
		return nil, NewParseError(p, nil, "Match must begin with 'match'")
	}

	toMatch, err := p.Expression([]int{TOKEN_KIND_BRACE_OPEN})

	if err != nil {
		return nil, NewParseError(p, err, "Cannot parse expression in match expression")
	}

	with, err := p.Pattern()

	if err != nil {
		return nil, NewParseError(p, err, "Cannot parse pattern in match expression")
	}

	return NewApplication([]AST{with, toMatch})
}

func (p *parser) Let(identifier Identifier) (AST, error) {
	if !p.ConsumeIfNext(TOKEN_KIND_EQUAL) {
		return nil, NewParseError(p, nil, ("'=' must follow identifier in let"))
	}

	value, err := p.Expression([]int{})

	if err != nil {
		return nil, NewParseError(p, err, ("Let must be assigned a value"))
	}

	body, err := p.Expression([]int{})

	if err != nil {
		return nil, NewParseError(p, err, ("Let must have a body"))
	}

	// If the body is a let expression, merge them
	switch b := body.(type) {
	case Let:
		return NewLet(
			append([]Identifier{identifier}, b.BoundIds...),
			append([]AST{value}, b.BoundValues...),
			b.Body,
		)
	default:
		return NewLet(
			[]Identifier{identifier},
			[]AST{value},
			body,
		)
	}
}

func (p *parser) List(first AST) (AST, error) {
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
				return nil, NewParseError(p, err, ("Cannot parse expression in list"))
			}
			list.Values = append(list.Values, val)

			if p.ConsumeIfNext(TOKEN_KIND_BRACKET_CLOSE) {
				break
			}

			// NOTE: !p.Peek().firstOfLine &&
			if !p.ConsumeIfNext(TOKEN_KIND_COMMA) {
				return nil, NewParseError(p, nil, ("Cannot parse comma in list"))
			}
		}
	}

	return list, nil
}

func (p *parser) ListConstructor(head AST) (AST, error) {
	if !p.ConsumeIfNext(TOKEN_KIND_COLON) {
		return nil, NewParseError(p, nil, ("List constructors must have a colon ':' separating the head and tail expressions"))
	}

	tail, err := p.Expression([]int{TOKEN_KIND_BRACKET_CLOSE})

	if err != nil {
		return nil, NewParseError(p, err, ("Cannot parse expression in list constructor"))
	}

	if !p.ConsumeIfNext(TOKEN_KIND_BRACKET_CLOSE) {
		return nil, NewParseError(p, nil, ("List constructors must be enclosed by brackets '[' ']'"))
	}

	return ListConstructor{Head: head, Tail: tail}, nil
}

func (p *parser) Where(match AST) (AST, error) {
	constantTime := false

	switch p.Peek().kind {
	case TOKEN_KIND_COLON:
		p.Next()
	case TOKEN_KIND_DOUBLE_COLON:
		p.Next()
		constantTime = true
	default:
		return nil, NewParseError(p, nil, ("Where must start with a colon"))
	}

	if !p.ConsumeIfNext(TOKEN_KIND_PAREN_OPEN) {
		return nil, NewParseError(p, nil, ("Where body must be enclosed by parenthesis '(' ')'"))
	}

	body, err := p.Expression([]int{TOKEN_KIND_PAREN_CLOSE})

	if err != nil {
		return nil, NewParseError(p, err, ("Cannot parse expression in where"))
	}

	if !p.ConsumeIfNext(TOKEN_KIND_PAREN_CLOSE) {
		return nil, NewParseError(p, nil, ("Where body must be enclosed by parenthesis '(' ')'"))
	}

	return NewWhere(match, body, constantTime)
}

func (p *parser) Pattern() (AST, error) {
	if !p.ConsumeIfNext(TOKEN_KIND_BRACE_OPEN) {
		return nil, NewParseError(p, nil, ("Pattern must be enclosed by braces '{' '}'"))
	}

	matchBodies := [][]AST{}
	bodies := []AST{}

	for !p.ConsumeIfNext(TOKEN_KIND_BRACE_CLOSE) {
		matches := []AST{}
		isImplicitBody := false

		if p.ConsumeIfNext(TOKEN_KIND_FAT_ARROW) {
			if len(matchBodies) == 0 {
				return nil, NewParseError(p, nil, ("Default pattern match cannot be the first body"))
			}

			for range matchBodies[0] {
				id, _ := NewIdentifier("_")
				matches = append(matches, id)
			}
		} else {
			for !p.ConsumeIfNext(TOKEN_KIND_ARROW) {
				match, err := p.Match()

				if err != nil {
					return nil, NewParseError(p, err, ("Cannot parse match in pattern"))
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
				return nil, NewParseError(p, nil, ("Pattern can only have one implicit true match"))
			}

			falseMatches := []AST{}

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
				return nil, NewParseError(p, err, ("Cannot parse expression in pattern"))
			}

			if len(matchBodies) > 0 && len(matches) != len(matchBodies[0]) {
				return nil, NewParseError(p, err, ("Pattern cannot take varying arguments"))
			}

			matchBodies = append(matchBodies, matches)
			bodies = append(bodies, body)
		}
	}

	return NewPattern(matchBodies, bodies)
}

// REFACTOR: where should apply to all match exprs
func (p *parser) Match() (AST, error) {
	var match AST
	var err error

	switch p.Peek().kind {
	case TOKEN_KIND_IDENTIFIER:
		match, err = p.Identifier()

		if err != nil {
			return nil, NewParseError(p, err, ("Cannot parse identifier in match"))
		}

	case TOKEN_KIND_LABEL:
		match, err = p.Label()

		if err != nil {
			return nil, NewParseError(p, err, ("Cannot parse label in match"))
		}

	case TOKEN_KIND_STRING:
		match, err = p.String()

		if err != nil {
			return nil, NewParseError(p, err, ("Cannot parse string in match"))
		}

	case TOKEN_KIND_NUMBER:
		match, err = p.Number()

		if err != nil {
			return nil, NewParseError(p, err, ("Cannot parse number in match"))
		}

	// Lists and list constructors: open bracket and an expression
	case TOKEN_KIND_BRACKET_OPEN:
		p.Next()

		if p.Peek().kind == TOKEN_KIND_BRACKET_CLOSE {
			match, err = p.List(nil)

			if err != nil {
				return nil, NewParseError(p, err, ("Cannot parse list in match"))
			}

			break
		}

		if p.Peek().kind == TOKEN_KIND_COLON {
			return nil, NewParseError(p, nil, ("List constructor cannot have an empty head value"))
		}

		expr, err := p.Expression([]int{TOKEN_KIND_COLON, TOKEN_KIND_COMMA, TOKEN_KIND_BRACKET_CLOSE})

		if err != nil {
			return nil, NewParseError(p, err, ("Cannot parse expression in match"))
		}

		if p.Peek().kind == TOKEN_KIND_COLON {
			match, err = p.ListConstructor(expr)

			if err != nil {
				return nil, NewParseError(p, err, ("Cannot parse list constructor in match"))
			}
		} else if p.ConsumeIfNext(TOKEN_KIND_COMMA) {
			match, err = p.List(expr)

			if err != nil {
				return nil, NewParseError(p, err, ("Cannot parse list in match"))
			}
		} else if p.ConsumeIfNext(TOKEN_KIND_BRACKET_CLOSE) {
			match = List{Values: []AST{expr}}
		} else {
			return nil, NewParseError(p, nil, ("Unable to parse list or list constructor in match expression"))
		}

	default:
		return nil, NewParseError(p, nil, ("Unexpected error occured when parsing match"))
	}

	// Match possible where after identifier
	if p.Peek().kind == TOKEN_KIND_COLON || p.Peek().kind == TOKEN_KIND_DOUBLE_COLON {
		matchWhere, err := p.Where(match)

		if err != nil {
			return nil, NewParseError(p, err, ("Cannot parse where in match"))
		}

		return matchWhere, nil
	}

	return match, nil
}

func (p *parser) PrimaryExpr(endTokenKinds []int) (AST, error) {
	if p.Peek().kind == TOKEN_KIND_IDENTIFIER {
		id, _ := NewIdentifier(string(p.Next().value))

		if p.Peek().kind == TOKEN_KIND_EQUAL {
			return p.Let(id)
		} else {
			return id, nil
		}
	}

	if p.Peek().kind == TOKEN_KIND_MODULE {
		return p.Module()
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
			return nil, NewParseError(p, err, ("Cannot parse expression in primary expression"))
		}

		if !p.ConsumeIfNext(TOKEN_KIND_PAREN_CLOSE) {
			return nil, NewParseError(p, nil, ("Application which starts with an open parenthesis must close with one"))
		}

		return res, nil
	}

	// Lists and list constructors: open bracket and an expression
	if p.Peek().kind == TOKEN_KIND_BRACKET_OPEN {
		p.Next()
		return p.List(nil)
	}

	return nil, NewParseError(p, nil, "Unexpected error occured when parsing an expression")
}

func (p *parser) OpExpr(precedence int, endTokenKinds []int) (AST, error) {
	if precedence >= len(opPrecedence) {
		exprs := []AST{}

		for {
			next, err := p.PrimaryExpr(endTokenKinds)

			if err != nil {
				return nil, NewParseError(p, err, ("Cannot parse primary expression in op expression"))
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
			return nil, NewParseError(p, nil, ("Cannot parse primary expression"))
		}

		if len(exprs) == 1 {
			return exprs[0], nil
		}

		return Application{Body: exprs}, nil
	} else {
		var head AST
		var err error
		head, err = p.OpExpr(precedence+1, append(append([]int{}, opPrecedence[precedence]...), endTokenKinds...))

		if err != nil {
			return nil, NewParseError(p, err, ("Cannot parse op expression in primary expression"))
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
						return nil, NewParseError(p, err, ("Cannot parse op expression in primary expression"))
					}

					head = Application{Body: []AST{Identifier{Value: string(t.value)}, head, next}}
					break
				}
			}

			if !foundOp {
				return nil, NewParseError(p, nil, ("Unexpected token found when attempting to parse operator"))
			}
		}

		return head, nil
	}
}

func (p *parser) Expression(endTokenKinds []int) (AST, error) {
	// Start parsing with default ender tokens
	expr, err := p.OpExpr(0, append(endTokenKinds, []int{TOKEN_KIND_EOF, TOKEN_KIND_SEMI_COLON}...))

	if err != nil {
		return nil, NewParseError(p, err, ("Cannot parse op expression in expression"))
	}

	if p.Peek().kind == TOKEN_KIND_SEMI_COLON {
		p.Next()
	}

	return expr, nil
}

func Parse(src []byte) (*SourceFile, error) {
	p := &parser{
		src,
		0,
		0,
		[numberOfSavedLines][]byte{},
		[]byte{},
	}

	for i := range p.savedLines {
		p.savedLines[i] = []byte{}
	}

	// Parse package then imports
	if !p.ConsumeIfNext(TOKEN_KIND_PACKAGE) {
		return nil, NewParseError(p, nil, "Must begin with a package name")
	}

	if p.Peek().kind != TOKEN_KIND_IDENTIFIER {
		return nil, NewParseError(p, nil, "Package name must be an identifier")
	}

	file := &SourceFile{}
	file.PackageName = string(p.Next().value)

	// Parse all imports
	for p.ConsumeIfNext(TOKEN_KIND_IMPORT) {
		if p.Peek().kind != TOKEN_KIND_STRING {
			return nil, NewParseError(p, nil, "Import path must be a string")
		}

		path := string(p.Next().value)

		if len(path) < 3 || path[len(path)-3:] != ".sl" {
			return nil, NewParseError(p, nil, "Invalid path string specified")
		}

		name := ""

		/* Import renaming removed for now
		if p.Peek().kind == TOKEN_KIND_IDENTIFIER {
			return nil, NewParseError(p, nil, "Import name must be an identifier")
		}

		name = string(p.Next().value)

		if name == "_" {
			return nil, NewParseError(p, nil, "Import name can't be discarded (_)")
		}*/

		file.Imports = append(file.Imports, SourceFileImport{path, name})
	}

	ast, err := p.Expression([]int{})

	if err != nil {
		return nil, NewParseError(p, err, "Cannot parse base expression")
	}

	if p.Next().kind != TOKEN_KIND_EOF {
		return nil, NewParseError(p, nil, fmt.Sprintf("[%d:%d] Unexpected end of parsing", p.line, p.char))
	}

	file.Definition = ast

	return file, nil
}
