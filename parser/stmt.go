package parser

import (
	"curgo/ast"
	"curgo/lexer"
)

func parseStmt(p *parser) ast.Stmt {
	stmtFn, exits := stmtLu[p.currentToken().Type]
	if exits {
		return stmtFn(p)
	}
	expression := parseExpr(p, defalt)
	p.expect(lexer.SEMI_COLON)
	return ast.ExprStmt{
		Expression: expression,
	}
}
