package config

import "fmt"

type parser struct {
    lexer   *lexer
    errs  []error

    currentToken    token
    peekToken       token
}

func newParser(lexer *lexer) *parser {
    p := &parser{lexer: lexer, errs: make([]error, 0)}
    p.nextToken()
    p.nextToken()
    return p
}

func (p *parser) errors() []error {
    return p.errs
}

func (p *parser) nextToken() {
    p.currentToken = p.peekToken
    p.peekToken = p.lexer.nextToken()
}

func (p *parser) parseProgram() *program {
    program := &program{}
    program.statements = make([]statement, 0)

    for p.currentToken.kind != EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            program.statements = append(program.statements, stmt)
        }
        p.nextToken()
    }

    return program
}

func (p *parser) parseStatement() statement {
    switch p.currentToken.kind {
    case IDENT:
        return p.parseStmt()
    default:
        return nil
    }
}

func (p *parser) parseStmt() *stmt {
    stmt := &stmt{token: p.currentToken}
    stmt.name = &identifier{token: p.currentToken, val: p.currentToken.literal}

    if !p.expectPeek(ASSIGN) {
        return nil
    }

    if !p.currentTokenIs(EOL) {
        p.nextToken()
    }

    return stmt
}

func (p *parser) expectPeek(kind tokenKind) bool {
    if p.peekTokenIs(kind) {
        p.nextToken()
        return true
    } else {
        p.peekError(kind)
        return false
    }
}

func (p *parser) currentTokenIs(kind tokenKind) bool {
    return p.currentToken.kind == kind
}

func (p *parser) peekTokenIs(kind tokenKind) bool {
    return p.peekToken.kind == kind
}

func (p *parser) peekError(kind tokenKind) {
    p.errs = append(p.errs, fmt.Errorf(
        "expected next token to be %s, got %s instead",
        kind, p.peekToken.kind,
    ))
}

