package config

import "testing"

func Test_nextToken(t *testing.T) {
    input := `title = "Elmo";
port = 8080;
debug = true;`

    tests := []struct{
        expectedKind    tokenKind
        expectedLiteral string
    }{
        {IDENT, "title"},
        {ASSIGN, "="},
        {STRING, "Elmo"},
        {SEMICOLON, ";"},
        {IDENT, "port"},
        {ASSIGN, "="},
        {INT, "8080"},
        {SEMICOLON, ";"},
        {IDENT, "debug"},
        {ASSIGN, "="},
        {TRUE, "true"},
        {SEMICOLON, ";"},
        {EOF, ""},
    }

    l := NewLexer(input)
    for i, tc := range tests {
        tok := l.nextToken()
        if tok.kind != tc.expectedKind {
            t.Fatalf("tests[%d] - token kind wrong. expected=%q, got=%q",
                i, tc.expectedKind, tok.kind)
        }

        if tok.literal != tc.expectedLiteral {
            t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
                i, tc.expectedLiteral, tok.literal)
        }
    }
}
