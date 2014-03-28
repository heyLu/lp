# exit: a program that calls `exit()`

# an assembler program has sections, this section normally has data in
# it that the program uses. (i think this would be variables in any
# other language, but probably there's more.)
.section .data

# the program text lives here.
.section .text
.globl _start # a special symbol that the linker uses to start the program

# call `exit(0)`
_start:
	# $<something> is the syntax for the immediate addressing mode
	movl $1, %eax # exit has the syscall number 1, defined in /usr/include/asm/unistd_32.h

	movl $0, %ebx # pass 0 as an argument

	int $0x80     # 0x80 is the syscall interupt for the kernel
