(define (emit instr . args)
  (display "\t")
  (display (apply format instr args))
  (display "\n"))

(define (compile-program x)
  (display ".globl scheme_entry\n\n")
  (display "scheme_entry:\n")

  (emit "movl $~a, %eax" x)
  (emit "ret"))

(compile-program 42)
