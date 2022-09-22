package parser

import (
	"Letter/ast"
	"Letter/tokenizer"
	"fmt"
	"strconv"
)

type Parser struct {
	lookAhead tokenizer.Token
	scanner   *tokenizer.Tokenizer
}

func New(source string) *Parser {
	p := &Parser{
		lookAhead: tokenizer.Token{
			Type:  tokenizer.EOF,
			Value: "",
		},
		scanner: tokenizer.New(source),
	}

	return p
}

func (p *Parser) Parse() *ast.Program {
	p.lookAhead = p.scanner.GetNextToken()
	return p.program()
}

// program ::= statementList
func (p *Parser) program() *ast.Program {
	return &ast.Program{Statements: p.statementList(tokenizer.EOF)}
}

// statementList ::= (statement)*
func (p *Parser) statementList(stopLookAhead tokenizer.TokenType) []ast.Statement {
	var statements []ast.Statement

	for p.lookAhead.Type != stopLookAhead {
		statements = append(statements, p.statement())
	}

	return statements
}

// statement ::= expressionStatement
func (p *Parser) statement() ast.Statement {
	return p.expressionStatement()
}

// expressionStatement ::= expression
func (p *Parser) expressionStatement() *ast.ExpressionStatement {
	expression := p.expression()
	return &ast.ExpressionStatement{Expression: expression}
}

// expression ::= assignmentExpression
func (p *Parser) expression() ast.Expression {
	return p.assignmentExpression()
}

// assignmentExpression ::= logicalOrExp | leftHandSideExp ('=','+=','-=','*=','/=') assignmentExpression
func (p *Parser) assignmentExpression() ast.Expression {
	left := p.logicalOrExpression()
	if !p.isAssignmentOperator(p.lookAhead.Type) {
		return left
	}

	return &ast.AssignmentExpression{
		Operator: p.assignmentOperator().Value,
		Left:     left,
		Right:    p.assignmentExpression(),
	}
}

// assignmentOperator
func (p *Parser) assignmentOperator() tokenizer.Token {
	if p.lookAhead.Type == tokenizer.SIMPLE_ASSIGN {
		return p.eat(tokenizer.SIMPLE_ASSIGN)
	}
	return p.eat(tokenizer.COMPLEX_ASSIGN)
}

// isAssignmentOperator
func (p *Parser) isAssignmentOperator(operator tokenizer.TokenType) bool {
	return operator == tokenizer.SIMPLE_ASSIGN || operator == tokenizer.COMPLEX_ASSIGN
}

// logicalOrExpression ::= logicalAndExpression ('||' logicalAndExpression)*
func (p *Parser) logicalOrExpression() ast.Expression {
	left := p.logicalAndExpression()
	for p.lookAhead.Type == tokenizer.LOGICAL_OR {
		operator := p.eat(tokenizer.LOGICAL_OR).Value
		left = &ast.LogicalExpression{
			Operator: operator,
			Left:     left,
			Right:    p.logicalAndExpression(),
		}
	}
	return left
}

// logicalAndExpression ::= equalityExpression ('&&' equalityExpression)*
func (p *Parser) logicalAndExpression() ast.Expression {
	left := p.equality()
	for p.lookAhead.Type == tokenizer.LOGICAL_AND {
		operator := p.eat(tokenizer.LOGICAL_AND).Value
		right := p.equality()
		left = &ast.LogicalExpression{
			Operator: operator,
			Left:     left,
			Right:    right,
		}
	}
	return left
}

// equality ::= comparison ('=='|'!=' comparison)*
func (p *Parser) equality() ast.Expression {
	left := p.comparison()
	for p.lookAhead.Type == tokenizer.EQUALITY_OPERATOR {
		operator := p.eat(tokenizer.EQUALITY_OPERATOR).Value
		right := p.comparison()
		left = &ast.LogicalExpression{
			Operator: operator,
			Left:     left,
			Right:    right,
		}
	}
	return left
}

// comparison ::= term ('<'|'>'|'<='|'>=' term)*
func (p *Parser) comparison() ast.Expression {
	left := p.term()

	for p.lookAhead.Type == tokenizer.RELATIONAL_OPERATOR {
		operator := p.eat(tokenizer.RELATIONAL_OPERATOR).Value
		right := p.term()
		left = &ast.BinaryExpression{
			Operator: operator,
			Left:     left,
			Right:    right,
		}
	}

	return left
}

// term ::= factor ('+'|'-' factor)*
func (p *Parser) term() ast.Expression {
	left := p.factor()

	for p.lookAhead.Type == tokenizer.ADDITIVE_OPERATOR {
		operator := p.eat(tokenizer.ADDITIVE_OPERATOR).Value
		right := p.factor()
		left = &ast.BinaryExpression{
			Operator: operator,
			Left:     left,
			Right:    right,
		}
	}

	return left
}

// factor ::= unary ('*'|'/' unary)*
func (p *Parser) factor() ast.Expression {
	left := p.unary()

	for p.lookAhead.Type == tokenizer.MULTIPLICATIVE_OPERATOR {
		operator := p.eat(tokenizer.MULTIPLICATIVE_OPERATOR).Value
		right := p.unary()
		left = &ast.BinaryExpression{
			Operator: operator,
			Left:     left,
			Right:    right,
		}
	}
	return left
}

// unary ::= ('+'|'-'|'!' unary) | leftHandSide
func (p *Parser) unary() ast.Expression {
	operator := ""
	switch p.lookAhead.Type {
	case tokenizer.ADDITIVE_OPERATOR:
		operator = p.eat(tokenizer.ADDITIVE_OPERATOR).Value
	case tokenizer.LOGICAL_NOT:
		operator = p.eat(tokenizer.LOGICAL_NOT).Value
	}
	if len(operator) > 0 {
		return &ast.UnaryExpression{
			Operator: operator,
			Right:    p.unary(), // right recursive e.g: +++foo, -+!bar
		}
	}

	return p.leftHandSideExp()
}

// leftHandSideExp ::= callMemberExp
func (p *Parser) leftHandSideExp() ast.Expression {
	return p.callMemberExp()
}

// callMemberExp ::= 'super' callExpression | memberExp | callExpression
func (p *Parser) callMemberExp() ast.Expression {
	// primero revisamos si es una llamada a super e.g super()
	if p.lookAhead.Type == tokenizer.SUPER {
		return p.callExpression(p.super())
	}

	// obtenemos el miembro
	member := p.memberExp()

	// revisar si se trata de un call e.g: person.getAge()
	if p.lookAhead.Type == tokenizer.LPAREN {
		return p.callExpression(member)
	}

	return member
}

// memberExp ::= primary | (primary '.' identifier) | (primary '[' expression ']')*
func (p *Parser) memberExp() ast.Expression {
	object := p.primary()
	for p.lookAhead.Type == tokenizer.DOT || p.lookAhead.Type == tokenizer.LBRACKET {
		if p.lookAhead.Type == tokenizer.DOT {
			p.eat(tokenizer.DOT)
			property := p.identifier()
			object = &ast.MemberExpression{Object: object, Property: property, Computed: false}
		} else {
			p.eat(tokenizer.LBRACKET)
			property := p.expression()
			p.eat(tokenizer.RBRACKET)
			object = &ast.MemberExpression{Object: object, Property: property, Computed: true}
		}
	}

	return object
}

// super
func (p *Parser) super() ast.Expression {
	p.eat(tokenizer.SUPER)
	return &ast.SuperExpression{}
}

// callExpression ::= Generic callExpression
func (p *Parser) callExpression(callee ast.Expression) ast.Expression {
	var callExpression ast.Expression
	callExpression = &ast.CallExpression{Callee: callee, Arguments: p.arguments()}

	if p.lookAhead.Type == tokenizer.LPAREN {
		callExpression = p.callExpression(callExpression)
	}

	return callExpression
}

// arguments ::= argumentList
func (p *Parser) arguments() []ast.Expression {
	var argumentList []ast.Expression
	p.eat(tokenizer.LPAREN)
	if p.lookAhead.Type != tokenizer.RPAREN {
		argumentList = p.argumentList()
	}
	p.eat(tokenizer.RPAREN)

	return argumentList
}

// argumentList ::= assignment (',' assignment)*
func (p *Parser) argumentList() []ast.Expression {
	var argumentList []ast.Expression

	argumentList = append(argumentList, p.assignmentExpression())
	for p.lookAhead.Type == tokenizer.COMMA {
		p.eat(tokenizer.COMMA)
		argumentList = append(argumentList, p.assignmentExpression())
	}

	return argumentList
}

// primary ::= literal | this | new | grouped | identifier
func (p *Parser) primary() ast.Expression {
	if p.isLiteral(p.lookAhead.Type) {
		return p.literal()
	}
	switch p.lookAhead.Type {
	case tokenizer.LPAREN:
		p.eat(tokenizer.LPAREN)
		exp := p.expression()
		p.eat(tokenizer.RPAREN)
		return exp
	case tokenizer.IDENTIFIER:
		return p.identifier()
	case tokenizer.THIS:
		p.eat(tokenizer.THIS)
		return &ast.ThisExpression{}
	case tokenizer.NEW:
		p.eat(tokenizer.NEW)
		return &ast.NewExpression{Callee: p.memberExp(), Arguments: p.arguments()}
	default:
		panic("Unexpected primary expression.")
	}
	return nil
}

func (p *Parser) identifier() ast.Expression {
	name := p.eat(tokenizer.IDENTIFIER).Value
	return &ast.Identifier{Name: name}
}

func (p *Parser) isLiteral(tokenType tokenizer.TokenType) bool {
	return tokenType == tokenizer.NUMBER ||
		tokenType == tokenizer.STRING ||
		tokenType == tokenizer.TRUE ||
		tokenType == tokenizer.FALSE ||
		tokenType == tokenizer.NULL
}

func (p *Parser) literal() ast.Expression {
	switch p.lookAhead.Type {
	case tokenizer.NUMBER:
		value, err := strconv.ParseInt(p.eat(tokenizer.NUMBER).Value, 0, 64)
		if err != nil {
			panic(err)
		}
		return &ast.NumericLiteral{Value: value}
	case tokenizer.STRING:
		value := p.eat(tokenizer.STRING).Value
		return &ast.StringLiteral{Value: value}
	case tokenizer.TRUE:
		p.eat(tokenizer.TRUE)
		return &ast.BooleanLiteral{Value: true}
	case tokenizer.FALSE:
		p.eat(tokenizer.FALSE)
		return &ast.BooleanLiteral{Value: false}
	case tokenizer.NULL:
		p.eat(tokenizer.NULL)
		return &ast.NullLiteral{}
	}
	panic("Literal: unexpected literal production.")
}

func (p *Parser) eat(tokenType tokenizer.TokenType) tokenizer.Token {
	token := p.lookAhead

	if token.Type == tokenizer.EOF {
		panic(fmt.Sprintf("Unexpected end of input, expected: %s", tokenizer.TokenNames[tokenType]))
	}

	if token.Type != tokenType {
		panic(fmt.Sprintf("Unexpected token: %s, expected: %s", tokenizer.TokenNames[token.Type], tokenizer.TokenNames[tokenType]))
	}

	newToken := p.scanner.GetNextToken()
	p.lookAhead = newToken

	return token
}
