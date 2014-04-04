#include <stdio.h>

#include <editline/readline.h>

int main(int argc, char** argv) {
	puts("mul v0.0.1\n");

	while (1) {
		char* input = readline("> ");
		add_history(input);

		printf("%s, yes.", input);

		free(input);
	}

	return 0;
}
