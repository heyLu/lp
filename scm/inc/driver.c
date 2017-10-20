/* a simple driver for scheme_entry */

#include <stdio.h>

#define fixnum_mask  3 // 11
#define fixnum_tag   0 // 00
#define fixnum_shift 2

#define char_mask 0xff  // 11111111
#define char_tag  15    // 00001111
#define char_shift 8

#define boolean_mask  0x3f // 1111111
#define boolean_tag   31   // 0011111
#define boolean_shift 7

#define empty_list   47 // 00101111

int scheme_entry();

int main(int argc, char **argv) {
	int val = scheme_entry();
	if ((val & fixnum_mask) == fixnum_tag) {
		printf("%d\n", val >> fixnum_shift);
	} else if ((val & char_mask) == char_tag) {
		printf("#\\%c\n", val >> char_shift);
	} else if ((val & boolean_mask) == boolean_tag) {
		printf("#%s\n", (val >> boolean_shift) == 1 ? "t" : "f");
	} else if (val == empty_list) {
		printf("()\n");
	} else {
		printf("Error: unhandled value: %d\n", val);
		return 1;
	}
	return 0;
}
