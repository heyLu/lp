/* a simple driver for scheme_entry */

#include <stdio.h>

#define fixnum_mask  3 // 11
#define fixnum_tag   0 // 00
#define fixnum_shift 2

#define empty_list   47 // 00101111

int scheme_entry();

int main(int argc, char **argv) {
	int val = scheme_entry();
	if ((val & fixnum_mask) == fixnum_tag) {
		printf("%d\n", val >> fixnum_shift);
	} else if (val == empty_list) {
		printf("()\n");
	} else {
		printf("Error: unhandled value: %d\n", val);
		return 1;
	}
	return 0;
}
