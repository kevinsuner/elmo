package config

import "bytes"

type node interface {
    tokenLiteral()  string
    value()         string
}

type statement interface {
    node
    statementNode()
}

type expression interface {
    node
    expressionNode()
}

type program struct {
    statements []statement
}

func (p *program) tokenLiteral() string {
    if len(p.statements) > 0 {
        return p.statements[0].tokenLiteral()
    } else {
        return ""
    }
}

func (p *program) value() string {
    var out bytes.Buffer
    for _, s := range p.statements {
        out.WriteString(s.value())
    }

    return out.String()
}

type identifier struct {
    token token
    val string
}

func (i *identifier) expressionNode() {}
func (i *identifier) tokenLiteral() string { return i.token.literal }
func (i *identifier) value() string { return i.val }

type stmt struct {
    token   token
    name    *identifier
    val     expression
}

func (s *stmt) statementNode() {}
func (s *stmt) tokenLiteral() string { return s.token.literal }
func (s *stmt) value() string {
    var out bytes.Buffer
    out.WriteString(s.tokenLiteral())
    out.WriteString(" = ")

    if s.val != nil {
        out.WriteString(s.val.value())
    }

    return out.String()
}

type expressionStmt struct {
    token token
    expression expression
}

func (e *expressionStmt) expressionNode() {}
func (e *expressionStmt) tokenLiteral() string { return e.token.literal }
func (e *expressionStmt) value() string {
    if e.expression != nil {
        return e.expression.value()
    }

    return ""
}
