(define (emit1 instr arg)
  (display "\t")
  (display (format instr arg))
  (display "\n"))

(define (emit0 instr)
  (display "\t")
  (display instr)
  (display "\n"))

(define (compile-program x)
  (display "\t.globl scheme_entry\n\n")
  (display "scheme_entry:\n")
  (emit1 "movl $~a, %eax" x)
  (emit0 "ret"))

(compile-program 42)
