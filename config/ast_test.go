package config

import "testing"

func Test_value(t *testing.T) {
    ast := &program{
        statements: []statement{
            &stmt{
                token: token{kind: IDENT, literal: "title"},
                name: &identifier{
                    token: token{kind: IDENT, literal: "title"},
                    val: "title",
                },
                val: &identifier{
                    token: token{kind: STRING, literal: "Elmo"},
                    val: "Elmo",
                },
            },
        },
    }

    if ast.value() != "title = Elmo" {
        t.Errorf("ast.value() wrong. got=%q", ast.value())
    }
}
