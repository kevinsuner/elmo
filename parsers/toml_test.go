package parsers

import (
	"testing"
)

func TestNextToken(t *testing.T) {
    input := `title = "TOML Example"
    port = 8000
    ports = [ 8001, 8002, 8003 ]
    user = { name = "john", pass = "1234" }`

    tests := []struct{
        expectedType    TokenType
        expectedLiteral string
    }{
        {IDENT, "title"},
        {ASSIGN, "="},
        {STRING, "TOML Example"},
        {IDENT, "port"},
        {ASSIGN, "="},
        {INT, "8000"},
        {IDENT, "ports"},
        {ASSIGN, "="},
        {LBRACKET, "["},
        {INT, "8001"},
        {COMMA, ","},
        {INT, "8002"},
        {COMMA, ","},
        {INT, "8003"},
        {RBRACKET, "]"},
        {IDENT, "user"},
        {ASSIGN, "="},
        {LBRACE, "{"},
        {IDENT, "name"},
        {ASSIGN, "="},
        {STRING, "john"},
        {COMMA, ","},
        {IDENT, "pass"},
        {ASSIGN, "="},
        {STRING, "1234"},
        {RBRACE, "}"},
        {EOF, ""},
    }

    l := New(input)

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
