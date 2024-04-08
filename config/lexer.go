package config

type tokenKind string
type token struct {
    kind    tokenKind
    literal string
}

const (
    ILLEGAL = "ILLEGAL"
    EOF     = "EOF"

    // Keywords
    TRUE    = "TRUE"
    FALSE   = "FALSE"

    // Identifiers
    IDENT   = "IDENT"

    // Literals
    INT     = "INT"
    STRING  = "STRING"

    // Operators
    ASSIGN  = "="

    // Delimiters
    EOL     = "EOL"
)

var keywords = map[string]tokenKind{
    "true":     TRUE,
    "false":    FALSE,
}

type lexer struct {
    input           string
    position        int     // current position in input (points to current char)
    readPosition    int     // current reading position in input (after current char)
    char            byte    // current character under examination
}

func newLexer(input string) *lexer {
    l := &lexer{input: input}
    l.readCharacter()
    return l
}

func (l *lexer) readCharacter() {
    if l.readPosition >= len(l.input) {
        l.char = 0
    } else {
        l.char = l.input[l.readPosition]
    }

    l.position = l.readPosition
    l.readPosition += 1
}

func (l *lexer) nextToken() token {
    var tok token
    l.eatWhitespace()

    switch l.char {
    case '=':
        tok = newToken(ASSIGN, l.char)
    case '"':
        tok.kind = STRING
        tok.literal = l.readString()
    case '\n':
        tok.kind = EOL
        tok.literal = ""
    case 0:
        tok.kind = EOF
        tok.literal = ""
    default:
        if isLetter(l.char) {
            tok.literal = l.readIdentifier()
            tok.kind = lookupIdentifier(tok.literal)
            return tok
        } else if isDigit(l.char) {
            tok.kind = INT
            tok.literal = l.readNumber()
            return tok
        } else {
            tok = newToken(ILLEGAL, l.char)
        }
    }

    l.readCharacter()
    return tok
}

func newToken(kind tokenKind, char byte) token {
    return token{kind: kind, literal: string(char)}
}

func isLetter(char byte) bool {
    return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
    return '0' <= char && char <= '9'
}

func lookupIdentifier(identifier string) tokenKind {
    if tok, ok := keywords[identifier]; ok {
        return tok
    }

    return IDENT
}

func (l *lexer) eatWhitespace() {
    for l.char == ' ' || l.char == '\t' || l.char == '\r' {
        l.readCharacter()
    }
}

func (l *lexer) readString() string {
    position := l.position + 1
    for {
        l.readCharacter()
        if l.char == '"' || l.char == 0 {
            break
        }
    }

    return l.input[position:l.position]
}

func (l *lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.char) {
        l.readCharacter()
    }

    return l.input[position:l.position]
}

func (l *lexer) readNumber() string {
    position := l.position
    for isDigit(l.char) {
        l.readCharacter()
    }

    return l.input[position:l.position]
}
