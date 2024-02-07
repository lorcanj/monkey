package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {

	// want to check what the below is doing
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are set
	p.NextToken()
	p.NextToken()

	return p
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
	default:
		return nil
	}
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
