#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>

int main(int argc, char *argv[]) {
	// machine code for:
	//   mov eax, 0
	//   ret
	unsigned char code[] = {0xb8, 0x00, 0x00, 0x00, 0x00, 0xc3};

	if (argc < 2) {
		fprintf(stderr, "Usage: hello_jit <integer>\n");
		return 1;
	}

	// copy the given number into the code, e.g:
	//   mov eax, <user's value>
	//   ret
	int num = atoi(argv[1]);
	memcpy(&code[1], &num, 4);

	// allocate writable & executable memory.
	// note: security risk!
	void *mem = mmap(NULL, sizeof(code), PROT_WRITE | PROT_EXEC, MAP_ANON | MAP_PRIVATE, -1, 0);
	memcpy(mem, code, sizeof(code));

	int (*func)() = mem;
	return func();
}
