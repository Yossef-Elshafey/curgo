package ast

/*
since parsing by what should be there and what shouldn't is not the best but its the fasted to go with
there has some parsing and simple interpretation to avoid lack of flexibility and appropiate errors
TODO: pratt parsing of each tiny token
TODO: simple interpreter that is capable of the following
			bit operations
			real assignments
			syntax determination
			assignment handling for statements and exper
TODO: build an AST which in not that powerful for higher precedence but easy to move with for simple operations and assignments
*/

// Currently TODO: implment closure struct to handle ranges and bodies

type Ast struct {
	Global *Global
}

func NewAst() *Ast {
	return &Ast{
		Global: NewGlobal(),
	}
}
