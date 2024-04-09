package config

import "testing"

func Test_statements(t *testing.T) {
    input := `title = "Elmo"
port = 8080
debug = true`

    l := newLexer(input)
    p := newParser(l)

    program := p.parseProgram()
    checkParserErrors(t, p)

    for _, stmt := range program.statements {
        t.Log(stmt)
    }

    if len(program.statements) != 3 {
        t.Fatalf("program.statements does not contain 3 statements, got=%d",
            len(program.statements))
    }

    tests := []struct{
        expectedIdentifier string
    }{
        {"title"},
        {"port"},
        {"debug"},
    }

    for i, tc := range tests {
        stmt := program.statements[i]
        if !testStatement(t, stmt, tc.expectedIdentifier) {
            return
        }
    }
}

func checkParserErrors(t *testing.T, p *parser) {
    errs := p.errors()
    if len(errs) == 0 {
        return
    }

    t.Errorf("parser has %d errors", len(errs))
    for _, err := range errs {
        t.Errorf("parser error: %v", err)
    }

    t.FailNow()
}

func testStatement(t *testing.T, s statement, name string) bool {
    if len(s.tokenLiteral()) == 0 {
        t.Error("s.tokenLiteral of zero length")
        return false
    }

    stmt, ok := s.(*stmt)
    if !ok {
        t.Errorf("s not *stmt. got=%T", s)
        return false
    }

    if stmt.name.val != name {
        t.Errorf("stmt.name.val not '%s'. got=%s", name, stmt.name.val)
        return false
    }

    if stmt.name.tokenLiteral() != name {
        t.Errorf("stmt.name.tokenLiteral not '%s'. got=%s",
            name, stmt.name.tokenLiteral())
        return false
    }

    return true
}

func Test_integerLiteralExpression(t *testing.T) {
    l := newLexer("5")
    p := newParser(l)
    program := p.parseProgram()
    checkParserErrors(t, p)

    if len(program.statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d",
            len(program.statements))
    }

    stmt, ok := program.statements[0].(*expressionStmt)
    if !ok {
        t.Fatalf("program.statements[0] is not expressionStmt. got=%T", 
            program.statements[0])
    }

    literal, ok := stmt.expression.(*integerLiteral)
    if !ok {
        t.Fatalf("exp not *integerLiteral. got=%T", stmt.expression)
    }

    if literal.val != 5 {
        t.Errorf("literal.val not %d. got=%d", 5, literal.val)
    }

    if literal.tokenLiteral() != "5" {
        t.Errorf("literal.tokenLiteral not %s. got=%s", "5",
            literal.tokenLiteral())
    }
}
