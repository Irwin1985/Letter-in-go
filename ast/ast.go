package ast

import (
	"bytes"
	"fmt"
	"strings"
)

type NodeType int

const (
	PROGRAM NodeType = iota
	EXPRESSION
	NUMBER
	STRING
	BOOLEAN
	NULL
	ASSIGNMENT
	LOGICAL
	BINARY
	UNARY
	CALL
	MEMBER
	SUPER
	IDENTIFIER
	THIS
	NEW
)

type Node interface {
	Type() NodeType
	ToString() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (s *Program) Type() NodeType {
	return PROGRAM
}

func (s *Program) ToString() string {
	var out bytes.Buffer

	for _, s := range s.Statements {
		out.WriteString(s.ToString())
	}

	return out.String()
}

type ExpressionStatement struct {
	Expression Expression
}

func (e *ExpressionStatement) Type() NodeType {
	return EXPRESSION
}

func (e *ExpressionStatement) ToString() string {
	return e.Expression.ToString()
}

func (e *ExpressionStatement) statementNode() {}

// NumericLiteral e.g: 123
type NumericLiteral struct {
	Value int64
}

func (e *NumericLiteral) Type() NodeType {
	return NUMBER
}

func (e *NumericLiteral) ToString() string {
	return fmt.Sprintf("%v", e.Value)
}

func (e *NumericLiteral) expressionNode() {}

// StringLiteral e.g "foo", 'bar'
type StringLiteral struct {
	Value string
}

func (e *StringLiteral) Type() NodeType {
	return STRING
}

func (e *StringLiteral) ToString() string {
	return e.Value
}

func (e *StringLiteral) expressionNode() {}

// BooleanLiteral e.g: true, false
type BooleanLiteral struct {
	Value bool
}

func (e *BooleanLiteral) Type() NodeType {
	return BOOLEAN
}

func (e *BooleanLiteral) ToString() string {
	if e.Value {
		return "true"
	}
	return "false"
}

func (e *BooleanLiteral) expressionNode() {}

// NullLiteral e.g: null
type NullLiteral struct {
	// nothing
}

func (e *NullLiteral) Type() NodeType {
	return NULL
}

func (e *NullLiteral) ToString() string {
	return "null"
}

func (e *NullLiteral) expressionNode() {}

// AssignmentExpression e.g foo = bar | foo += bar
type AssignmentExpression struct {
	Operator string
	Left     Expression
	Right    Expression
}

func (e *AssignmentExpression) Type() NodeType {
	return ASSIGNMENT
}

func (e *AssignmentExpression) ToString() string {
	return fmt.Sprintf("%s %s %s", e.Left.ToString(), e.Operator, e.Right.ToString())
}

func (e *AssignmentExpression) expressionNode() {}

// LogicalExpression e.g foo || bar
type LogicalExpression struct {
	Operator string
	Left     Expression
	Right    Expression
}

func (e *LogicalExpression) Type() NodeType {
	return LOGICAL
}

func (e *LogicalExpression) ToString() string {
	return fmt.Sprintf("(%s%s%s)", e.Left.ToString(), e.Operator, e.Right.ToString())
}

func (e *LogicalExpression) expressionNode() {}

// BinaryExpression e.g foo + bar
type BinaryExpression struct {
	Operator string
	Left     Expression
	Right    Expression
}

func (e *BinaryExpression) expressionNode() {}

func (e *BinaryExpression) Type() NodeType {
	return BINARY
}

func (e *BinaryExpression) ToString() string {
	return fmt.Sprintf("(%s%s%s)", e.Left.ToString(), e.Operator, e.Right.ToString())
}

// UnaryExpression e.g -foo, +bar, !baz
type UnaryExpression struct {
	Operator string
	Right    Expression
}

func (e *UnaryExpression) Type() NodeType {
	return UNARY
}

func (e *UnaryExpression) ToString() string {
	return fmt.Sprintf("(%s%s)", e.Operator, e.Right.ToString())
}

func (e *UnaryExpression) expressionNode() {}

// CallExpression e.g foo(), bar(a, b)
type CallExpression struct {
	Callee    Expression
	Arguments []Expression
}

func (e *CallExpression) Type() NodeType {
	return CALL
}

func (e *CallExpression) ToString() string {
	var arguments []string
	for _, a := range e.Arguments {
		arguments = append(arguments, a.ToString())
	}
	return fmt.Sprintf("%s(%s)", e.Callee.ToString(), strings.Join(arguments, ", "))
}

func (e *CallExpression) expressionNode() {}

// MemberExpression e.g foo.bar, bar.baz(), foo.bar[2]
type MemberExpression struct {
	Computed bool
	Object   Expression
	Property Expression
}

func (e *MemberExpression) Type() NodeType {
	return MEMBER
}

func (e *MemberExpression) ToString() string {
	var format string
	if !e.Computed {
		format = "%s.%s"
	} else {
		format = "%s[%s]"
	}
	return fmt.Sprintf(format, e.Object.ToString(), e.Property.ToString())
}

func (e *MemberExpression) expressionNode() {}

// SuperExpression e.g Super()
type SuperExpression struct {
	// nothing
}

func (e *SuperExpression) Type() NodeType {
	return SUPER
}

func (e *SuperExpression) ToString() string {
	return "super"
}

func (e *SuperExpression) expressionNode() {}

// Identifier e.g foo
type Identifier struct {
	Name string
}

func (e *Identifier) Type() NodeType {
	return IDENTIFIER
}

func (e *Identifier) ToString() string {
	return e.Name
}

func (e *Identifier) expressionNode() {}

// ThisExpression e.g this
type ThisExpression struct {
	// nothing
}

func (e *ThisExpression) Type() NodeType {
	return THIS
}

func (e *ThisExpression) ToString() string {
	return "this"
}

func (e *ThisExpression) expressionNode() {}

// NewExpression e.g this
type NewExpression struct {
	Callee    Expression
	Arguments []Expression
}

func (e *NewExpression) Type() NodeType {
	return NEW
}

func (e *NewExpression) ToString() string {
	var arguments []string
	for _, a := range e.Arguments {
		arguments = append(arguments, a.ToString())
	}
	return fmt.Sprintf("new %s(%s)", e.Callee.ToString(), strings.Join(arguments, ", "))
}

func (e *NewExpression) expressionNode() {}
