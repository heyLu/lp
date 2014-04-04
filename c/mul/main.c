#include <math.h>
#include <stdio.h>

#include <editline/readline.h>

#include "mpc.h"

long eval_binary_op(char* op, long x, long y) {
	if (strcmp(op, "+") == 0) { return x + y; }
	if (strcmp(op, "-") == 0) { return x - y; }
	if (strcmp(op, "*") == 0) { return x * y; }
	if (strcmp(op, "/") == 0) { return x / y; }
	if (strcmp(op, "%") == 0) { return x % y; }
	if (strcmp(op, "^") == 0) { return pow(x, y); }
	if (strcmp(op, "min") == 0) { return x < y ? x : y; }
	if (strcmp(op, "max") == 0) { return x > y ? x : y; }
	return 0;
}

long eval_unary_op(char* op, long x) {
	if (strcmp(op, "-") == 0) { return -x; }
	return x;
}

long eval(mpc_ast_t* t) {
	if (strstr(t->tag, "number")) {
		return atoi(t->contents);
	}

	char* op = t->children[1]->contents;
	long x = eval(t->children[2]);

	if (t->children_num == 4) {
		return eval_unary_op(op, x);
	} else {
		int i = 3;
		while (strstr(t->children[i]->tag, "expr")) {
			x = eval_binary_op(op, x, eval(t->children[i]));
			i++;
		}

		return x;
	}
}

int main(int argc, char** argv) {
	int debug_mode = 0;
	if (getenv("DEBUG")) {
		debug_mode = 1;
	}

	mpc_parser_t* Number = mpc_new("number");
	mpc_parser_t* Operator = mpc_new("operator");
	mpc_parser_t* Expr = mpc_new("expr");
	mpc_parser_t* Lang = mpc_new("lang");

	mpca_lang(MPC_LANG_DEFAULT,
		"\
number   : /-?[0-9]+(\\.[0-9]+)?/ ; \
operator : '+' | '-' | '*' | '/' | '%' | '^' | /[a-zA-Z-]+/ ; \
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
			mpc_ast_t* t = r.output;
			if (debug_mode) {
				mpc_ast_print(t->children[1]);
			}
			long result = eval(t->children[1]);
			printf("%li\n", result);
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
