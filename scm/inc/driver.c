/* a simple driver for scheme_entry */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

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

#define cons_tag  1 // 001
#define cons_mask 7 // 111

#define MAX_MEMORY (1 << 20)

typedef long ptr;

ptr scheme_entry(int *heap);

void print_value(ptr val) {
	if ((val & fixnum_mask) == fixnum_tag) {
		printf("%d", (int)val >> fixnum_shift);
	} else if ((val & char_mask) == char_tag) {
		printf("#\\%c", (int)val >> char_shift);
	} else if ((val & boolean_mask) == boolean_tag) {
		printf("#%s", (val >> boolean_shift) == 1 ? "t" : "f");
	} else if (val == empty_list) {
		printf("()");
	} else if ((val & cons_mask) == cons_tag) {
		printf("(");
		print_value(*(int*)(val-1));
		printf(" ");
		print_value(*(int*)(val+3));
		printf(")");
	} else {
		printf("\nError: unhandled value: %d 0x%x 0x%p\n", (int)val, (int)val, (int*)val);
		exit(1);
	}
}

int main(int argc, char **argv) {
	int *mem = (int*)malloc(MAX_MEMORY * sizeof(int));
	if (mem == NULL) {
		perror("malloc");
		exit(1);
	}
	memset(mem, 0, MAX_MEMORY * sizeof(int));
	ptr val = scheme_entry(mem);
	print_value(val);
	printf("\n");
	return 0;
}
