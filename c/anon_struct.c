/*
 * Anonymous/Unnamed structs in structs (also works w/ unions).
 *
 * See https://lwn.net/SubscriberLink/548560/26d15e832d21a483/ and
 *  http://gcc.gnu.org/onlinedocs/gcc/Unnamed-Fields.html for details.
 */
#include <stdio.h>

// if I could do this, that would allow composition (over inheritance)
// in c?
//struct awesome {
//	int the_answer;
//};

struct composable {
	char *name;
	struct awesome {
		int the_answer;
	};
};

int main(void) {
	printf("-- anonymous structures in C11\n\n");
	struct composable c = {
		.name = "helo",
		.the_answer = 42
	};
	printf("composable: %s, %d\n", c.name, c.the_answer);
	return 0;
}
