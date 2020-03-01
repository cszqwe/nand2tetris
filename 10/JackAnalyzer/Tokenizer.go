package main

import "fmt"

type jackTokenType int

const (
	keyword jackTokenType = iota
	symbol
	intergerConstant
	stringConstant
	identifier
)

type keywordTokenType string

const (
	notAKeyWord keywordTokenType = "this is not a keyword"
	class       keywordTokenType = "class"
	constructor keywordTokenType = "constructor"
	function    keywordTokenType = "function"
	method      keywordTokenType = "method"
	field       keywordTokenType = "field"
	static      keywordTokenType = "static"
	_var        keywordTokenType = "var"
	_int        keywordTokenType = "int"
	char        keywordTokenType = "char"
	boolean     keywordTokenType = "boolean"
	void        keywordTokenType = "void"
	_true       keywordTokenType = "true"
	_false      keywordTokenType = "false"
	null        keywordTokenType = "null"
	this        keywordTokenType = "this"
	let         keywordTokenType = "let"
	do          keywordTokenType = "do"
	_if         keywordTokenType = "if"
	_else       keywordTokenType = "else"
	while       keywordTokenType = "while"
	_return     keywordTokenType = "return"
)

type symbolTokenType byte

const (
	lsb       symbolTokenType = '('
	rsb       symbolTokenType = ')'
	lmb       symbolTokenType = '['
	rmb       symbolTokenType = ']'
	llb       symbolTokenType = '{'
	rlb       symbolTokenType = '}'
	dot       symbolTokenType = '.'
	comma     symbolTokenType = ','
	semicolon symbolTokenType = ';'
	plus      symbolTokenType = '+'
	minus     symbolTokenType = '-'
	mul       symbolTokenType = '*'
	div       symbolTokenType = '/'
	and       symbolTokenType = '&'
	or        symbolTokenType = '|'
	bigger    symbolTokenType = '>'
	lesser    symbolTokenType = '<'
	not       symbolTokenType = '~'
	equal     symbolTokenType = '='
)

type token interface {
	getType() jackTokenType
	print() string
}

type keywordToken struct {
	keywordType keywordTokenType
}

func (token keywordToken) getType() jackTokenType {
	return keyword
}

func (token keywordToken) print() string {
	return fmt.Sprintf("<keyword> %s </keyword>", string(token.keywordType))
}

type symbolToken struct {
	symbolType symbolTokenType
}

func (token symbolToken) getType() jackTokenType {
	return symbol
}

func (token symbolToken) print() string {
	if token.symbolType == bigger {
		return fmt.Sprintf("<symbol> &gt; </symbol>")
	}
	if token.symbolType == lesser {
		return fmt.Sprintf("<symbol> &lt; </symbol>")
	}
	return fmt.Sprintf("<symbol> %s </symbol>", string(token.symbolType))
}

type integerConstantToken struct {
	val int
}

func (token integerConstantToken) getType() jackTokenType {
	return intergerConstant
}

func (token integerConstantToken) print() string {
	return fmt.Sprintf("<integerConstant> %d </integerConstant>", token.val)
}

type stringConstantToken struct {
	val string
}

func (token stringConstantToken) getType() jackTokenType {
	return stringConstant
}

func (token stringConstantToken) print() string {
	return fmt.Sprintf("<stringConstant> %s </stringConstant>", token.val)
}

type identifierToken struct {
	val string
}

func (token identifierToken) getType() jackTokenType {
	return identifier
}

func (token identifierToken) print() string {
	return fmt.Sprintf("<identifier> %s </identifier>", token.val)
}

func getTokens(s string) []token {
	tokens := []token{}
	var t token
	for len(s) != 0 {
		if s[0] == '(' || s[0] == ')' ||
			s[0] == '[' || s[0] == ']' ||
			s[0] == '{' || s[0] == '}' ||
			s[0] == '.' || s[0] == ',' || s[0] == ';' ||
			s[0] == '+' || s[0] == '-' ||
			s[0] == '*' || s[0] == '/' ||
			s[0] == '&' || s[0] == '|' ||
			s[0] == '~' || s[0] == '=' ||
			s[0] == '>' || s[0] == '<' {
			t, s = tokenizeSymbol(s)
			tokens = append(tokens, t)
		} else if s[0] >= '0' && s[0] <= '9' {
			t, s = tokenizeIntegerConstant(s)
			tokens = append(tokens, t)
		} else if s[0] == '"' {
			t, s = tokenizeStringConstant(s)
			tokens = append(tokens, t)
		} else if s[0] == '_' || s[0] >= 'a' && s[0] <= 'z' || s[0] >= 'A' && s[0] <= 'Z' {
			t, s = tokenizeIdentifierOrKeyword(s)
			tokens = append(tokens, t)
		} else if s[0] == ' ' || s[0] == 9 {
			s = s[1:]
		} else {
			panic(fmt.Sprintf("cannot tokenize:%d", s[0]))
		}
	}
	return tokens
}

func tokenizeSymbol(s string) (token, string) {
	ch := s[0]
	s = s[1:]
	switch ch {
	case '(':
		return symbolToken{lsb}, s
	case ')':
		return symbolToken{rsb}, s
	case '[':
		return symbolToken{lmb}, s
	case ']':
		return symbolToken{rmb}, s
	case '{':
		return symbolToken{llb}, s
	case '}':
		return symbolToken{rlb}, s
	case '.':
		return symbolToken{dot}, s
	case ',':
		return symbolToken{comma}, s
	case ';':
		return symbolToken{semicolon}, s
	case '+':
		return symbolToken{plus}, s
	case '-':
		return symbolToken{minus}, s
	case '*':
		return symbolToken{mul}, s
	case '/':
		return symbolToken{div}, s
	case '&':
		return symbolToken{and}, s
	case '|':
		return symbolToken{or}, s
	case '>':
		return symbolToken{bigger}, s
	case '<':
		return symbolToken{lesser}, s
	case '~':
		return symbolToken{not}, s
	case '=':
		return symbolToken{equal}, s
	default:
		panic(fmt.Sprintf("invalid symbol: %s", string(s[0])))
	}
}

func tokenizeIntegerConstant(s string) (token, string) {
	number := 0
	for len(s) > 0 && s[0] >= '0' && s[0] <= '9' {
		number = number*10 + int(s[0]-'0')
		s = s[1:]
	}
	return integerConstantToken{number}, s
}

func tokenizeStringConstant(s string) (token, string) {
	str := []byte{}
	s = s[1:]
	for len(s) > 0 && s[0] != '"' {
		str = append(str, s[0])
		s = s[1:]
	}
	s = s[1:]
	return stringConstantToken{string(str)}, s
}

func tokenizeIdentifierOrKeyword(s string) (token, string) {
	str := []byte{}
	for len(s) > 0 && (isNumber(s[0]) || isAlphaBeta(s[0]) || s[0] == '_') {
		str = append(str, s[0])
		s = s[1:]
	}
	curKeywordTokenType := getKeyword(string(str))
	if curKeywordTokenType == notAKeyWord {
		return identifierToken{string(str)}, s
	}
	return keywordToken{curKeywordTokenType}, s
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlphaBeta(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func getKeyword(s string) keywordTokenType {
	switch s {
	case "class":
		return class
	case "constructor":
		return constructor
	case "function":
		return function
	case "method":
		return method
	case "field":
		return field
	case "static":
		return static
	case "var":
		return _var
	case "int":
		return _int
	case "char":
		return char
	case "boolean":
		return boolean
	case "void":
		return void
	case "true":
		return _true
	case "false":
		return _false
	case "null":
		return null
	case "this":
		return this
	case "let":
		return let
	case "do":
		return do
	case "if":
		return _if
	case "else":
		return _else
	case "while":
		return while
	case "return":
		return _return
	default:
		return notAKeyWord
	}
}
