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
