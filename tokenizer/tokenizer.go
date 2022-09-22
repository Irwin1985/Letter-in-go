package tokenizer

import (
	"regexp"
)

type TokenType uint8

const (
	LPAREN    TokenType = iota // (
	RPAREN                     // )
	LBRACE                     // {
	RBRACE                     // }
	LBRACKET                   // [
	RBRACKET                   // ]
	COMMA                      // ,
	SEMICOLON                  // ;
	DOT                        // .

	// Keywords
	LET
	IF
	ELSE
	TRUE
	FALSE
	NULL
	DEF
	RETURN

	// Iterators keywords
	WHILE
	DO
	FOR

	// OOP Keywords
	CLASS
	THIS
	EXTENDS
	SUPER
	NEW

	// Literals
	NUMBER
	STRING
	IDENTIFIER

	// Operators
	SIMPLE_ASSIGN
	COMPLEX_ASSIGN
	RELATIONAL_OPERATOR
	EQUALITY_OPERATOR
	ADDITIVE_OPERATOR
	MULTIPLICATIVE_OPERATOR
	LOGICAL_AND
	LOGICAL_OR
	LOGICAL_NOT

	// NONE
	IGNORE

	EOF
	ERROR
)

var TokenNames = [40]string{
	"LPAREN",
	"RPAREN",
	"LBRACE",
	"RBRACE",
	"LBRACKET",
	"RBRACKET",
	"COMMA",
	"SEMICOLON",
	"DOT",

	// Keywords
	"LET",
	"IF",
	"ELSE",
	"TRUE",
	"FALSE",
	"NULL",
	"DEF",
	"RETURN",

	// Iterators keywords
	"WHILE",
	"DO",
	"FOR",

	// OOP Keywords
	"CLASS",
	"THIS",
	"EXTENDS",
	"SUPER",
	"NEW",

	// Literals
	"NUMBER",
	"STRING",
	"IDENTIFIER",

	// Operators
	"SIMPLE_ASSIGN",
	"COMPLEX_ASSIGN",
	"RELATIONAL_OPERATOR",
	"EQUALITY_OPERATOR",
	"ADDITIVE_OPERATOR",
	"MULTIPLICATIVE_OPERATOR",
	"LOGICAL_AND",
	"LOGICAL_OR",
	"LOGICAL_NOT",

	// NONE
	"IGNORE",

	"EOF",
	"ERROR",
}

type Token struct {
	Type  TokenType
	Value string
}

type Tuple struct {
	RegEx string
	Type  TokenType
}

var Spec = [41]Tuple{
	// Whitespace
	{RegEx: `^\s+`, Type: IGNORE},
	// Single comment
	{RegEx: `^//.*`, Type: IGNORE},
	// Multiple comments
	{RegEx: `^\/\*[\s\S]*?\*\/`, Type: IGNORE},
	// Symbols and delimiters
	{RegEx: `^\{`, Type: LBRACE},
	{RegEx: `^\}`, Type: RBRACE},
	{RegEx: `^\(`, Type: LPAREN},
	{RegEx: `^\)`, Type: RPAREN},
	{RegEx: `^\[`, Type: LBRACKET},
	{RegEx: `^\]`, Type: RBRACKET},
	{RegEx: `^,`, Type: COMMA},
	{RegEx: `^;`, Type: SEMICOLON},
	{RegEx: `^\.`, Type: DOT},
	// Relational Operators
	{RegEx: `^[<>]=?`, Type: RELATIONAL_OPERATOR},
	{RegEx: `^[!=]=`, Type: EQUALITY_OPERATOR},

	// Logical Operators
	{RegEx: `^&&`, Type: LOGICAL_AND},
	{RegEx: `^\|\|`, Type: LOGICAL_OR},
	{RegEx: `^!`, Type: LOGICAL_NOT},

	// Assignment Operators
	{RegEx: `^=`, Type: SIMPLE_ASSIGN},
	{RegEx: `^[\+\-\*\/]=`, Type: COMPLEX_ASSIGN},

	// Math operators: +, -, *, /
	{RegEx: `^[+\-]`, Type: ADDITIVE_OPERATOR},
	{RegEx: `^[\*\/]`, Type: MULTIPLICATIVE_OPERATOR},

	// Keywords
	{RegEx: `^\blet\b`, Type: LET},
	{RegEx: `^\bif\b`, Type: IF},
	{RegEx: `^\belse\b`, Type: ELSE},
	{RegEx: `^\btrue\b`, Type: TRUE},
	{RegEx: `^\bfalse\b`, Type: FALSE},
	{RegEx: `^\bnull\b`, Type: NULL},
	{RegEx: `^\bdef\b`, Type: DEF},
	{RegEx: `^\breturn\b`, Type: RETURN},

	// OOP keywords
	{RegEx: `^\bclass\b`, Type: CLASS},
	{RegEx: `^\bthis\b`, Type: THIS},
	{RegEx: `^\bextends\b`, Type: EXTENDS},
	{RegEx: `^\bsuper\b`, Type: SUPER},
	{RegEx: `^\bnew\b`, Type: NEW},

	// Iterators keywords
	{RegEx: `^\bwhile\b`, Type: WHILE},
	{RegEx: `^\bdo\b`, Type: DO},
	{RegEx: `^\bfor\b`, Type: FOR},

	// Literals
	{RegEx: `^\d+`, Type: NUMBER},
	{RegEx: `^"[^"]*"`, Type: STRING},
	{RegEx: `^'[^']*'`, Type: STRING},
	{RegEx: `^\w+`, Type: IDENTIFIER},
}

type Tokenizer struct {
	input  string
	cursor int
}

// New Lazily pulls a token from a stream.
func New(input string) *Tokenizer {
	scanner := &Tokenizer{input: input, cursor: 0}
	return scanner
}

// hasMoreTokens Whether we still have more tokens.
func (t *Tokenizer) hasMoreTokens() bool {
	return t.cursor < len(t.input)
}

// GetNextToken Obtains the next token.
func (t *Tokenizer) GetNextToken() Token {
	if !t.hasMoreTokens() {
		return Token{Type: EOF, Value: ""}
	}
	input := t.input[t.cursor:]

	for _, regExp := range Spec {
		tokenValue := t.match(regExp.RegEx, input)
		// Couldn't match this rule, so continue to the next one...
		if len(tokenValue) == 0 {
			continue
		}
		// Check for IGNORE token which means it could be a whitespace or skipped token
		if regExp.Type == IGNORE {
			// call GetNextToken again to repeat all regular expression rules.
			return t.GetNextToken()
		}
		// Finally we return the token
		return Token{Type: regExp.Type, Value: tokenValue}
	}

	return Token{Type: ERROR, Value: ""}
}

// match: Matches a token for a regular expression.
func (t *Tokenizer) match(regExp string, input string) string {
	pattern := regexp.MustCompile(regExp)
	matched := pattern.FindString(input)
	if len(matched) > 0 {
		t.cursor += len(matched)
	}
	return matched
}
