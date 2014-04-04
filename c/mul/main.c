#include <stdio.h>

int main(int argc, char** argv) {
	puts("mul v0.0.1\n");

	while (1) {
		fputs("> ", stdout);

		fgets(input, 2048, stdin);

		printf("%s, yes.", input);
	}

	return 0;
}
