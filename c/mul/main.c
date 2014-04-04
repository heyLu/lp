#include <stdio.h>

#include <editline/readline.h>

#include "mpc.h"

int main(int argc, char** argv) {
	mpc_parser_t* Number = mpc_new("number");
	mpc_parser_t* Operator = mpc_new("operator");
	mpc_parser_t* Expr = mpc_new("expr");
	mpc_parser_t* Lang = mpc_new("lang");

	mpca_lang(MPC_LANG_DEFAULT,
		"\
number   : /-?[0-9]+(\\.[0-9]+)?/ ; \
operator : '+' | '-' | '*' | '/' | /[a-zA-Z-]+/ ; \
expr     : <number> | '(' <operator> <expr>+ ')' ; \
lang     : /^/ <expr>+ /$/ ; \
",
		Number, Operator, Expr, Lang);

	puts("mul v0.0.1\n");

	while (1) {
		char* input = readline("> ");
		add_history(input);

		mpc_result_t r;
		if (mpc_parse("<stdin>", input, Lang, &r)) {
			mpc_ast_print(r.output);
			mpc_ast_delete(r.output);
		} else {
			mpc_err_print(r.error);
			mpc_err_delete(r.error);
		}

		free(input);
	}

	mpc_cleanup(4, Number, Operator, Expr, Lang);
	return 0;
}
