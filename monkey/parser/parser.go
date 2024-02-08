package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	// here iota gives the below constants incrementing numbers as values
	// the underscore takes the first value (0) so the consts start from 1
	// can use these numbers to assign the order of operations
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	// these will store a map of tokenTypes to parsing functions
	// means that we might have to parse the same token differently depending on
	// whether it appears prefix or infix in our statements
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// we can create types in Go, as shown below
// of the form identifier Source
// in the example below, both of the types are functions that return
// an ast.Expression, however for infix this also requires an ast.Expression as input
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// slightly confused about whether this means that there is only one function
// per tokenType
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {

	// want to check what the below is doing
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// Read two tokens, so curToken and peekToken are set
	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	// do we now want to add it to the errors array?
	p.errors = append(p.errors, msg)
}

// on first call, curToken is empty and peek token is set to the first
// token for the input
// then on the second call, curToken is set to peekToken, which is the first token
// and peekToken is set to the next token in the list of tokens
func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// the below is returning a pointer to the ast.Program
// which is the root node of the ast
func (p *Parser) ParseProgram() *ast.Program {
	// the ampersand symbol points to the address of the stored value
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		// ParseStatement is passing the literal values and not
		// pointers
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}
	return program
}

// are we doing this just for let statements or for all statments?
// assume we want this to work for all statements eventually
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	// don't need to check whether current token is equal to return here
	// because this has been checked already further down the function stack
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.NextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

// why should this pass the pointer and the one above pass the value itself??
func (p *Parser) parseLetStatement() ast.Statement {
	// want to create a new node in the AST
	// the current token is 'let' as shown in the parseStatement method above
	stmt := &ast.LetStatement{Token: p.curToken}

	// the token after 'let' should be an identifier and so enforcing that here
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// want to double check whether the below is being correctly assigned
	// as not 100% sure
	// might want to add another unit test
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// could alter this because no point doing the assignment above when the below fails
	// but can't move this above as expectPeek increments to the next tokem
	// which doesn't seem clear from the name of the method
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// should take in the expression and create a new ast node of ExpressionStatement type
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	// passing in const LOWEST
	stmt.Expression = p.parseExpression(LOWEST)

	// advancing token if next token is a semi-colon
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	// not failing if expression doesn't end with a semi-colon
	// which means can add expression like 5 + 5 to REPL

	return stmt
}

// function to check whether there is a prefix function associated with
// the token type and calls and returns the result of that call if found
func (p *Parser) parseExpression(precedence int) ast.Expression {

	// now the parseExpression function will use the parsing function that we have
	// added to the prefixParse map that was added when Parse was initialised
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}
