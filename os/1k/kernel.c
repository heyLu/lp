#include "common.h"
#include "kernel.h"

typedef unsigned char uint8_t;
typedef unsigned int uint32_t;
typedef uint32_t size_t;

extern char __bss[], __bss_end[], __stack_top[];
extern char __free_ram[], __free_ram_end[];

struct process procs[PROCS_MAX]; // all processes
struct process *current_proc;    // currently running
struct process *idle_proc;       // default idle

struct sbiret sbi_call(long arg0, long arg1, long arg2, long arg3, long arg4,
                       long arg5, long fid, long eid) {
  register long a0 __asm__("a0") = arg0;
  register long a1 __asm__("a1") = arg1;
  register long a2 __asm__("a2") = arg2;
  register long a3 __asm__("a3") = arg3;
  register long a4 __asm__("a4") = arg4;
  register long a5 __asm__("a5") = arg5;
  register long a6 __asm__("a6") = fid;
  register long a7 __asm__("a7") = eid;

  __asm__ __volatile__("ecall"
                       : "=r"(a0), "=r"(a1)
                       : "r"(a0), "r"(a1), "r"(a2), "r"(a3), "r"(a4), "r"(a5),
                         "r"(a6), "r"(a7)
                       : "memory");
  return (struct sbiret){.error = a0, .value = a1};
}

void putchar(char ch) {
  sbi_call(ch, 0, 0, 0, 0, 0, 0, 1 /* Console Putchar */);
}

__attribute__((naked)) __attribute__((aligned(4))) void kernel_entry(void) {
  __asm__ __volatile__(

      // retrieve kernel stack of running process from sscratch
      "csrrw sp, sscratch, sp\n"

      "addi sp, sp, -4 * 31\n" // reserve space on for stored/restored registers
      // store registers
      "sw ra, 4 * 0(sp)\n"
      "sw gp,  4 * 1(sp)\n"
      "sw tp,  4 * 2(sp)\n"
      "sw t0,  4 * 3(sp)\n"
      "sw t1,  4 * 4(sp)\n"
      "sw t2,  4 * 5(sp)\n"
      "sw t3,  4 * 6(sp)\n"
      "sw t4,  4 * 7(sp)\n"
      "sw t5,  4 * 8(sp)\n"
      "sw t6,  4 * 9(sp)\n"
      "sw a0,  4 * 10(sp)\n"
      "sw a1,  4 * 11(sp)\n"
      "sw a2,  4 * 12(sp)\n"
      "sw a3,  4 * 13(sp)\n"
      "sw a4,  4 * 14(sp)\n"
      "sw a5,  4 * 15(sp)\n"
      "sw a6,  4 * 16(sp)\n"
      "sw a7,  4 * 17(sp)\n"
      "sw s0,  4 * 18(sp)\n"
      "sw s1,  4 * 19(sp)\n"
      "sw s2,  4 * 20(sp)\n"
      "sw s3,  4 * 21(sp)\n"
      "sw s4,  4 * 22(sp)\n"
      "sw s5,  4 * 23(sp)\n"
      "sw s6,  4 * 24(sp)\n"
      "sw s7,  4 * 25(sp)\n"
      "sw s8,  4 * 26(sp)\n"
      "sw s9,  4 * 27(sp)\n"
      "sw s10, 4 * 28(sp)\n"
      "sw s11, 4 * 29(sp)\n"

      // retrieve and save the sp at time of the exception
      "csrr a0, sscratch\n"
      "sw a0, 4 * 30(sp)\n"

      // reset kernel stack
      "addi a0, sp, 4 * 31\n"
      "csrw sscratch, a0\n"

      "csrr a0, sscratch\n"
      "sw a0, 4 * 30(sp)\n"

      "mv a0, sp\n"
      "call handle_trap\n"

      // restore registers
      "lw ra,  4 * 0(sp)\n"
      "lw gp,  4 * 1(sp)\n"
      "lw tp,  4 * 2(sp)\n"
      "lw t0,  4 * 3(sp)\n"
      "lw t1,  4 * 4(sp)\n"
      "lw t2,  4 * 5(sp)\n"
      "lw t3,  4 * 6(sp)\n"
      "lw t4,  4 * 7(sp)\n"
      "lw t5,  4 * 8(sp)\n"
      "lw t6,  4 * 9(sp)\n"
      "lw a0,  4 * 10(sp)\n"
      "lw a1,  4 * 11(sp)\n"
      "lw a2,  4 * 12(sp)\n"
      "lw a3,  4 * 13(sp)\n"
      "lw a4,  4 * 14(sp)\n"
      "lw a5,  4 * 15(sp)\n"
      "lw a6,  4 * 16(sp)\n"
      "lw a7,  4 * 17(sp)\n"
      "lw s0,  4 * 18(sp)\n"
      "lw s1,  4 * 19(sp)\n"
      "lw s2,  4 * 20(sp)\n"
      "lw s3,  4 * 21(sp)\n"
      "lw s4,  4 * 22(sp)\n"
      "lw s5,  4 * 23(sp)\n"
      "lw s6,  4 * 24(sp)\n"
      "lw s7,  4 * 25(sp)\n"
      "lw s8,  4 * 26(sp)\n"
      "lw s9,  4 * 27(sp)\n"
      "lw s10, 4 * 28(sp)\n"
      "lw s11, 4 * 29(sp)\n"
      "lw sp,  4 * 30(sp)\n"
      "sret\n"

  );
}

void handle_trap(struct trap_frame *f) {
  uint32_t scause = READ_CSR(scause);
  uint32_t stval = READ_CSR(stval);
  uint32_t user_pc = READ_CSR(sepc);

  PANIC("unexpected trap: scause=%x, stval=%x, sepc=%x\n", scause, stval,
        user_pc);
}

paddr_t alloc_pages(uint32_t n) {
  static paddr_t next_paddr =
      (paddr_t)__free_ram; // static makes this a global variable??
  paddr_t paddr = next_paddr;
  next_paddr += n * PAGE_SIZE;

  if (next_paddr > (paddr_t)__free_ram_end) {
    PANIC("out of memory");
  }

  memset((void *)paddr, 0, n * PAGE_SIZE);
  return paddr;
}

__attribute__((naked)) void switch_context(uint32_t *prev_sp,
                                           uint32_t *next_sp) {
  __asm__ __volatile__(
      "addi sp, sp, -13 * 4\n" // allocate stack space for 13 4-byte registers

      "sw ra,  0  * 4(sp)\n" // save callee-saved registers only
      "sw s0,  1  * 4(sp)\n"
      "sw s1,  2  * 4(sp)\n"
      "sw s2,  3  * 4(sp)\n"
      "sw s3,  4  * 4(sp)\n"
      "sw s4,  5  * 4(sp)\n"
      "sw s5,  6  * 4(sp)\n"
      "sw s6,  7  * 4(sp)\n"
      "sw s7,  8  * 4(sp)\n"
      "sw s8,  9  * 4(sp)\n"
      "sw s9,  10 * 4(sp)\n"
      "sw s10, 11 * 4(sp)\n"
      "sw s11, 12 * 4(sp)\n"

      "sw sp, (a0)\n" // *prev_sp = sp;
      "lw sp, (a1)\n" // switch stack pointer here (to next_sp)

      "lw ra,  0  * 4(sp)\n"
      "lw s0,  1  * 4(sp)\n"
      "lw s1,  2  * 4(sp)\n"
      "lw s2,  3  * 4(sp)\n"
      "lw s3,  4  * 4(sp)\n"
      "lw s4,  5  * 4(sp)\n"
      "lw s5,  6  * 4(sp)\n"
      "lw s6,  7  * 4(sp)\n"
      "lw s7,  8  * 4(sp)\n"
      "lw s8,  9  * 4(sp)\n"
      "lw s9,  10 * 4(sp)\n"
      "lw s10, 11 * 4(sp)\n"
      "lw s11, 12 * 4(sp)\n"

      "addi sp, sp, 13 * 4\n" // un-allocate stack space
      "ret\n");
}

struct process *create_process(uint32_t pc) {
  // find unused process
  struct process *proc = NULL;
  int i;
  for (i = 0; i < PROCS_MAX; i++) {
    if (procs[i].state == PROC_UNUSED) {
      proc = &procs[i];
      break;
    }
  }

  if (!proc) {
    PANIC("no free process slots");
  }

  // stack callee-saved registers, will be restored by the first context switch
  // in switch_context
  uint32_t *sp = (uint32_t *)&proc->stack[sizeof(proc->stack)];
  *--sp = 0; // s11
  *--sp = 0;
  *--sp = 0;
  *--sp = 0;
  *--sp = 0;
  *--sp = 0;
  *--sp = 0;
  *--sp = 0;
  *--sp = 0;
  *--sp = 0;
  *--sp = 0;
  *--sp = 0;            // s0
  *--sp = (uint32_t)pc; // ra

  // initialize fields
  proc->pid = i + 1;
  proc->state = PROC_RUNNABLE;
  proc->sp = (uint32_t)sp;
  return proc;
}

void yield(void) {
  struct process *next = idle_proc;
  for (int i = 0; i < PROCS_MAX; i++) {
    struct process *proc = &procs[(current_proc->pid + i) % PROCS_MAX];
    if (proc->state == PROC_RUNNABLE && proc->pid > 0) {
      next = proc;
      break;
    }
  }

  // nothing new runnable, continue with current
  if (next == current_proc) {
    return;
  }

  __asm__ __volatile__(
      "csrw sscratch, %[sscratch]\n"
      :
      : [sscratch] "r"((uint32_t)&next->stack[sizeof(next->stack)]));

  // context switch
  struct process *prev = current_proc;
  current_proc = next;
  switch_context(&prev->sp, &next->sp);
}

struct process *proc_a;
struct process *proc_b;

void proc_a_entry(void) {
  printf("starting process A\n");
  while (1) {
    putchar('A');
    yield();

    // waste a bit of time to slow down output
    for (int i = 0; i < 30000000; i++) {
      __asm__ __volatile__("nop");
    }
  }
}

void proc_b_entry(void) {
  printf("starting process B\n");
  while (1) {
    putchar('B');
    yield();

    // waste a bit of time to slow down output
    for (int i = 0; i < 30000000; i++) {
      __asm__ __volatile__("nop");
    }
  }
}

void kernel_main(void) {
  memset(__bss, 0, (size_t)__bss_end - (size_t)__bss);

  WRITE_CSR(stvec, (uint32_t)kernel_entry);

  idle_proc = create_process((uint32_t)NULL);
  idle_proc->pid = -1; // idle
  current_proc = idle_proc;

  proc_a = create_process((uint32_t)proc_a_entry);
  proc_b = create_process((uint32_t)proc_b_entry);

  yield();

  PANIC("switched to idle process");
}

__attribute__((section(".text.boot"))) __attribute__((naked)) void boot(void) {
  __asm__ __volatile__(
      "mv sp, %[stack_top]\n" // set the stack pointer
      "j kernel_main\n"       // jump to kernel main function
      :
      : [stack_top] "r"(
          __stack_top) // pass the stack top address as %[stack_top]
  );
}
