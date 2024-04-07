package parsers

import (
	"bytes"
	"fmt"
)

type TokenType string

type Token struct {
    Type TokenType
    Literal string
}

const (
    ILLEGAL = "ILLEGAL"
    EOL     = "EOL"
    EOF     = "EOF"

    // Keywords
    TRUE    = "TRUE"
    FALSE   = "FALSE"

    // Identifiers + literals
    IDENT   = "IDENT"
    INT     = "INT"
    STRING  = "STRING"

    // Operators
    ASSIGN = "="

    // Delimiters
    COMMA       = ","
    LBRACE      = "{"
    RBRACE      = "}"
    LBRACKET    = "["
    RBRACKET    = "]"
)

var keywords = map[string]TokenType{
    "true": TRUE,
    "false": FALSE,
}

type Lexer struct {
    input       string
    position    int // current position in input (points to current char)
    readPos     int // current reading position in input (after current char)
    ch          byte // current char under examination
}

func NewLexer(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) NextToken() Token {
    var tok Token

    l.skipWhitespace()

    switch l.ch {
    case '=':
        tok = newToken(ASSIGN, l.ch)
    case ',':
        tok = newToken(COMMA, l.ch)
    case '{':
        tok = newToken(LBRACE, l.ch)
    case '}':
        tok = newToken(RBRACE, l.ch)
    case '[':
        tok = newToken(LBRACKET, l.ch)
    case ']':
        tok = newToken(RBRACKET, l.ch)
    case '"':
        tok.Type = STRING
        tok.Literal = l.readStr()
    case '\n':
        tok.Literal = ""
        tok.Type = EOL
    case 0:
        tok.Literal = ""
        tok.Type = EOF
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdent()
            tok.Type = lookupIdent(tok.Literal)
            return tok
        } else if isDigit(l.ch) {
            tok.Type = INT
            tok.Literal = l.readNum()
            return tok
        } else {
            tok = newToken(ILLEGAL, l.ch)
        }
    }

    l.readChar()
    return tok
}

func lookupIdent(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }

    return IDENT
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) readStr() string {
    position := l.position + 1
    for {
        l.readChar()
        if l.ch == '"' || l.ch == 0 {
            break
        }
    }

    return l.input[position:l.position]
}

func (l *Lexer) readNum() string {
    position := l.position
    for isDigit(l.ch) {
        l.readChar()
    }

    return l.input[position:l.position]
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdent() string {
    position := l.position
    for isLetter(l.ch) {
        l.readChar()
    }

    return l.input[position:l.position]
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType TokenType, ch byte) Token {
    return Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readChar() {
    if l.readPos >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPos]
    }
    
    l.position = l.readPos
    l.readPos += 1
}

type Node interface {
    TokenLiteral() string
    String() string
}

type Statement interface {
    Node
    statementNode()
}

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

func (p *Program) String() string {
    var out bytes.Buffer
    for _, s := range p.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

type Stmt struct {
    Token Token
    Name *Ident
    Value Expression
}

func (s *Stmt) statementNode() {}
func (s *Stmt) TokenLiteral() string { return s.Token.Literal }
func (s *Stmt) String() string {
    var out bytes.Buffer
    out.WriteString(s.TokenLiteral())
    out.WriteString(" = ")

    if s.Value != nil {
        out.WriteString(s.Value.String())
    }

    return out.String()
}

type ExprStmt struct {
    Token Token
    Expression Expression
}

func (e *ExprStmt) statementNode() {}
func (e *ExprStmt) TokenLiteral() string { return e.Token.Literal }
func (e *ExprStmt) String() string {
    if e.Expression != nil {
        return e.Expression.String()
    }

    return ""
}

type Ident struct {
    Token Token
    Value string
}

func (i *Ident) expressionNode() {}
func (i *Ident) TokenLiteral() string { return i.Token.Literal }
func (i *Ident) String() string { return i.Value }

type Parser struct {
    l       *Lexer
    errors  []string

    curToken    Token
    peekToken   Token
}

func NewParser(l *Lexer) *Parser {
    p := &Parser{l: l, errors: []string{}}
    p.nextToken()
    p.nextToken()
    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) peekError(t TokenType) {
    msg := fmt.Sprintf("expected next token to be %s, got %s instead",
        t, p.peekToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
    program := &Program{}
    program.Statements = []Statement{}

    for p.curToken.Type != EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.nextToken()
    }

    return program
}

func (p *Parser) parseStatement() Statement {
    switch p.curToken.Type {
    case IDENT:
        return p.parseStmt()
    default:
        return nil 
    }
}

func (p *Parser) parseStmt() *Stmt {
    stmt := &Stmt{Token: p.curToken}
    stmt.Name = &Ident{Token: p.curToken, Value: p.curToken.Literal}

    if !p.expectPeek(ASSIGN) {
        return nil 
    }

    for !p.curTokenIs(EOL) {
        p.nextToken()
    }

    return stmt
}

func (p *Parser) curTokenIs(t TokenType) bool {
    return p.curToken.Type == t
} 

func (p *Parser) peekTokenIs(t TokenType) bool {
    return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
    if p.peekTokenIs(t) {
        p.nextToken()
        return true
    } else {
        p.peekError(t)
        return false
    }
}
