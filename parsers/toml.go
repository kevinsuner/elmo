package parsers

type TokenType string

type Token struct {
    Type TokenType
    Literal string
}

const (
    ILLEGAL = "ILLEGAL"
    EOF     = "EOF"

    // Identifiers + literals
    IDENT   = "IDENT"
    INT     = "INT"
    STRING  = "STRING"

    // Operators
    ASSIGN = "="

    // Delimiters
    COMMA       = ","
)

type Lexer struct {
    input       string
    position    int // current position in input (points to current char)
    readPos     int // current reading position in input (after current char)
    ch          byte // current char under examination
}

func New(input string) *Lexer {
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
    case '"':
        tok.Type = STRING
        tok.Literal = l.readStr()
    case 0:
        tok.Literal = ""
        tok.Type = EOF
    default:
        if isLetter(l.ch) {
            tok.Type = IDENT
            tok.Literal = l.readIdent()
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

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
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

