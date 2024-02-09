package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/helper_functions"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
		`
	// below will tokenise the input
	l := lexer.New(input)
	// then call the parser on the tokens
	// generated by the lexer
	p := New(l)
	// program will be a pointer to an instance of
	// the ast which holds a list of statements
	program := p.ParseProgram()
	checkParserErrors(t, p)
	// shouldn't be nil because our parser should be able to parse
	// let statements at this point
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// added for readability
	// could add something that parses the number of statements by counting the number of newline characters maybe?
	expectedLength := 3
	helper_functions.CheckProgramLength(t, len(program.Statements), expectedLength)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let. got=%q", s.TokenLiteral())
		return false
	}

	// ast.Statement is an interface so below we are checking that the type of
	// statement that has been passed to this function is of the concrete type LetStatement
	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	// want to check letStmt.Name and what stuff it has access to
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

// TODO: might be worth at some point adding line and col numbers to the output
// to make it easier to spot where the issues are
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 989389382;
		`
	l := lexer.New(input)
	p := New(l)

	// what is the output of the below?
	program := p.ParseProgram()
	// can do this but won't currently do anything yet
	checkParserErrors(t, p)
	// probably should create a separate function as duplicating code here
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// added for readability
	// could add something that parses the number of statements by counting the number of newline characters maybe?
	expectedLength := 3
	helper_functions.CheckProgramLength(t, len(program.Statements), expectedLength)

	for _, stmt := range program.Statements {
		// want to check what the below is doing
		// is this creating an ast.ReturnStatement node from the stmt
		// and if it fails then the bool ok is set to false
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got%T", stmt)
			continue
		}
		// so a return statment should
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
		// currently this does not include expressions
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// added for readability
	// could add something that parses the number of statements by counting the number of newline characters maybe?
	expectedLength := 1
	helper_functions.CheckProgramLength(t, len(program.Statements), expectedLength)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("expression not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	// added for readability
	// could add something that parses the number of statements by counting the number of newline characters maybe?
	expectedLength := 1
	helper_functions.CheckProgramLength(t, len(program.Statements), expectedLength)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt is not ast.IntegerLiteral, got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %s. got=%s", "5", literal.TokenLiteral())
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		expectedLength := 1
		helper_functions.CheckProgramLength(t, len(program.Statements), expectedLength)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T",
				stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

// problem with this is that can fail first test and so won't check any of the
// tests below this
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}
