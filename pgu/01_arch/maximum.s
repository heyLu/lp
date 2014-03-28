# maximum - calculates the maximum of a fixed list of numbers
#   returns the maximum via exit

.section .data
numbers:
	.long 3,67,34,222,65,45,75,54,34,44,33,22,11,66,0

.section .text
.globl _start

# %eax current number
# %ebx current maximum
# %edi current offset

_start:
	movl $0, %edi # initial offset is 0
	movl numbers(,%edi,4), %eax # load first number into %eax
	movl %eax, %ebx # first number is also first maximum

	# apparrently then the program just goes to the next thing? what
	# would be if we'd put loop_start and loop_end before _start?

loop_start: 
	# 0 marks the end of the numbers, finish if %eax == 0
	cmpl $0, %eax
	je loop_exit

	# load next number
	incl %edi
	movl numbers(,%edi,4), %eax

	# if the current number is smaller (or equal?), go back to loop_start
	cmpl %ebx, %eax
	jle loop_start

	# otherwise store the current number as the maximum number
	movl %eax, %ebx
	jmp loop_start

loop_exit:
	# exit. conveniently the maximum is already stored in %ebx
	movl $1, %eax
	int $0x80
