(define (emit instr . args)
  (display "\t")
  (display (apply format instr args))
  (display "\n"))

(define fixnum-shift 2)
(define char-shift 8)
(define boolean-shift 7)

(define (compile-program x)
  (define (immediate-rep x)
    (cond
      ((integer? x) (bitwise-arithmetic-shift x fixnum-shift))
      ((char? x) (bitwise-xor (bitwise-arithmetic-shift (char->integer x) 8) char-shift))
      ((boolean? x) (bitwise-xor (bitwise-arithmetic-shift (cond
                                                             ((boolean=? x #t) 1)
                                                             ((boolean=? x #f) 0))
                                                           boolean-shift)
                                 #b0011111))
      ((eq? x '()) #b00101111)))

  (display ".globl scheme_entry\n\n")
  (display "scheme_entry:\n")

  (emit "movl $~a, %eax" (immediate-rep x))
  (emit "ret"))

(compile-program '())
