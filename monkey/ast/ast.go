package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string
}

// Distinction between Statements and Expressions

// A Statement does not produce a value, for example let x = 5
type Statement interface {
	Node
	statementNode()
}

// An Expression does produce a value, for example add(5, 5)
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Although the Name in a let statment is being defined as an Identifier
// this contradicts what we wrote above in the example let x = 5, as the
// Identifier implements the Expression interface but 'x' in the example
// is not producing a value and so in this example should be an Expression
type LetStatement struct {
	Token token.Token // e.g. the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // e.g. the token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type Identifier struct {
	Token token.Token // e.g. the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
