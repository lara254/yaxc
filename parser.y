%{
package main
import (
	"fmt"
	"text/scanner"
	"strconv"
	"strings"
)

type Expression interface{}

type IntLiteral struct {
	Value int
}

type Var struct {
	Name string
}

type Token struct {
    token int
    literal string
}

type IfExpr struct {
	Cond Expression
	Then Expression
	Else Expression
}

type LetExpr struct {
	Bindings []Binding
	Body     Expression
}

type Binding struct {
	Name  string
	Value Expression
}

type Application struct {
	Func Expression
	Args []Expression
}

type DefineExpr struct {
	Name  string
	Value Expression
}

type BinaryOp struct {
	Operator string
	Left     Expression
	Right    Expression
}

type WhileExpr struct {
    Cnd Expression
    Body Expression
  }
type SetExpr struct {
	Name  string
	Value Expression
}

type BeginExpr struct {
	Exprs []Expression
}

%}

%union {
	token  Token
	expr   Expression
	str    string
	intval int
}

%type<expr> program
%type<expr> expr
%type<expr> binding
%type<expr> expr_list

%token<str> NAME
%token<intval> INTEGER
%token LPAREN RPAREN PLUS LT GT

%token<str> LET
%token IF DEFINE LAMBDA EQ WHILE SET BEGIN

%start program

%%

program:
	expr {
		$$ = $1
		yylex.(*Lexer).result = $$
	}

expr:
	INTEGER {
		$$ = IntLiteral{Value: $1}
	}
	| NAME {
		$$ = Var{Name: $1}
	}
	| LPAREN LET LPAREN binding RPAREN expr RPAREN {
		$$ = LetExpr{Bindings: []Binding{$4.(Binding)}, Body: $6}
	}
	| LPAREN IF expr expr expr RPAREN {
		$$ = IfExpr{Cond: $3, Then: $4, Else: $5}
	}
	| LPAREN DEFINE NAME expr RPAREN {
		$$ = DefineExpr{Name: $3, Value: $4}
	}
	| LPAREN expr expr_list RPAREN {
		$$ = Application{Func: $2, Args: $3.([]Expression)}
	}
	| LPAREN PLUS expr expr RPAREN {
		$$ = BinaryOp{Operator: "+", Left: $3, Right: $4}
	}
	| LPAREN LT expr expr RPAREN {
		$$ = BinaryOp{Operator: "<", Left: $3, Right: $4}
	}
	| LPAREN GT expr expr RPAREN {
		$$ = BinaryOp{Operator: ">", Left: $3, Right: $4}
	}
        | LPAREN WHILE expr expr RPAREN {
		$$ = WhileExpr{Cnd: $3, Body: $4}
	}
        | LPAREN SET NAME expr RPAREN {
	       $$ = SetExpr{Name: $3, Value: $4}
	}
        | LPAREN BEGIN expr_list RPAREN {
               $$ = BeginExpr{Exprs: $3.([]Expression)}
        }


binding:
	LPAREN NAME expr RPAREN {
		$$ = Binding{Name: $2, Value: $3}
	}

expr_list:
	/* empty */ {
		$$ = []Expression{}
	}
	| expr expr_list {
		$$ = append([]Expression{$1}, $2.([]Expression)...)
	}

%%

type Lexer struct {
	scanner.Scanner
	result Expression
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok := l.Scan()
	lit := l.TokenText()

	switch tok {
	case scanner.Int:
		tokVal, _ := strconv.Atoi(lit)
		lval.intval = tokVal
		return INTEGER
	case '(':
		return LPAREN
	case ')':
		return RPAREN
	case scanner.Ident:
		switch lit {
		case "if":
			return IF
		case "let":
			return LET
		case "define":
			return DEFINE
	       	case "while":
			return WHILE
		case "set":
			return SET
		case "begin":
			return BEGIN
		default:
			lval.str = lit
			return NAME
		}
	case '+':
		return PLUS
	case '<':
		return LT
	case '>':
		return GT
	case '=':
		return EQ
	}

	return 0
}

func (l *Lexer) Error(e string) {
	fmt.Printf("Lex error: %s\n", e)
}

func Parse(input string) (Expression, error) {
	lexer := &Lexer{}
	lexer.Init(strings.NewReader(input))

	if yyParse(lexer) == 0 {
		return lexer.result, nil
	} else {
		return nil, fmt.Errorf("Parsing error")
	}
}
