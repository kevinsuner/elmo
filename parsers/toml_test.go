package parsers

import (
	"testing"
)

func TestNextToken(t *testing.T) {
    input := `title = "TOML Example"
    port = 8000
    ports = [ 8001, 8002, 8003 ]
    data = [ ["delta", "phi"],  [3] ]
    user = { name = "john", age = 28 }`

    tests := []struct{
        expectedType    TokenType
        expectedLiteral string
    }{
        {IDENT, "title"},
        {ASSIGN, "="},
        {STRING, "TOML Example"},
        {EOL, ""},
        {IDENT, "port"},
        {ASSIGN, "="},
        {INT, "8000"},
        {EOL, ""},
        {IDENT, "ports"},
        {ASSIGN, "="},
        {LBRACKET, "["},
        {INT, "8001"},
        {COMMA, ","},
        {INT, "8002"},
        {COMMA, ","},
        {INT, "8003"},
        {RBRACKET, "]"},
        {EOL, ""},
        {IDENT, "data"},
        {ASSIGN, "="},
        {LBRACKET, "["},
        {LBRACKET, "["},
        {STRING, "delta"},
        {COMMA, ","},
        {STRING, "phi"},
        {RBRACKET, "]"},
        {COMMA, ","},
        {LBRACKET, "["},
        {INT, "3"},
        {RBRACKET, "]"},
        {RBRACKET, "]"},
        {EOL, ""},
        {IDENT, "user"},
        {ASSIGN, "="},
        {LBRACE, "{"},
        {IDENT, "name"},
        {ASSIGN, "="},
        {STRING, "john"},
        {COMMA, ","},
        {IDENT, "age"},
        {ASSIGN, "="},
        {INT, "28"},
        {RBRACE, "}"},
        {EOF, ""},
    }

    l := NewLexer(input)

    for i, tc := range tests {
        tok := l.NextToken()
        t.Log(tok)

        if tok.Type != tc.expectedType {
            t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
                i, tc.expectedType, tok.Type)
        }

        if tok.Literal != tc.expectedLiteral {
            t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
                i, tc.expectedLiteral, tok.Literal)
        }
    }
}

func TestStatements(t *testing.T) {
    input := `title = "TOML Example"
port = 8000
name = "John"
age = 28
    `

    l := NewLexer(input)
    p := NewParser(l)

    program := p.ParseProgram()
    if program == nil {
        t.Fatalf("ParseProgram() returned nil")
    }

    if len(program.Statements) != 4 {
        t.Fatalf("program.Statements does not contain 4 statements. got=%d",
            len(program.Statements))
    }

    tests := []struct{
        expectedIdentifier string
    }{
        {"title"},
        {"port"},
        {"name"},
        {"age"},
    }

    for i, tc := range tests {
        stmt := program.Statements[i]
        if !testStatement(t, stmt, tc.expectedIdentifier) {
            return
        }
    }
}

func testStatement(t *testing.T, s Statement, name string) bool {
    if len(s.TokenLiteral()) == 0 {
        t.Error("s.TokenLiteral of zero length")
        return false
    }

    stmt, ok := s.(*Stmt)
    if !ok {
        t.Errorf("s not *Stmt. got=%T", s)
        return false
    }

    if stmt.Name.Value != name {
        t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
        return false
    }

    if stmt.Name.TokenLiteral() != name {
        t.Errorf("stmt.Name.TokenLiteral() not '%s'. got=%s",
            name, stmt.Name.TokenLiteral())
        return false
    }

    return true
}
