;; assembly resources:
;;  - https://en.wikipedia.org/wiki/X86_instruction_listings

(define (emit instr . args)
  (display "\t")
  (display (apply format instr args))
  (display "\n"))

(define fixnum-shift 2)
(define char-shift 8)
(define boolean-shift 7)

(define (immediate-rep x)
  (cond
    ((integer? x) (bitwise-arithmetic-shift x fixnum-shift))
    ((char? x) (bitwise-xor (bitwise-arithmetic-shift (char->integer x) char-shift) #b00001111))
    ((boolean? x) (bitwise-xor (bitwise-arithmetic-shift (cond
                                                           ((boolean=? x #t) 1)
                                                           ((boolean=? x #f) 0))
                                                         boolean-shift)
                               #b0011111))
    ((eq? x '()) #b00101111)))

(define (immediate? x)
  (or (integer? x) (char? x) (boolean? x) (eq? x '())))

(define (primcall? x)
  (list? x))

(define (primcall-op x)
  (car x))

(define (primcall-operand1 x)
  (cadr x))

(define (emit-expr x)
  (cond
    ((immediate? x)
     (emit "movl $~a, %eax" (immediate-rep x)))
    ((primcall? x)
     (case (primcall-op x)
       ((add1)
        (emit-expr (primcall-operand1 x))
        (emit "addl $~a, %eax" (immediate-rep 1)))
       ((integer->char)
        (emit-expr (primcall-operand1 x))
        (emit "shl $6, %eax")
        (emit "xorl $15, %eax"))
       ((char->integer)
        (emit-expr (primcall-operand1 x))
        (emit "shrl $6, %eax"))
       ((zero?)
        (emit-expr (primcall-operand1 x))
        (emit "cmpl $0,  %eax")                ; x == 0
        (emit "movl $0,  %eax")                ; zero %eax to put the result of the comparison into
        (emit "sete %al")                      ; set low byte of %eax to 1 if cmp succeeded
        (emit "sall $~a,  %eax" boolean-shift) ; construct correctly tagged boolean value
        (emit "xorl $31, %eax"))
       ))))

(define (compile-program x)
  (display ".globl scheme_entry\n\n")
  (display "scheme_entry:\n")

  (emit-expr x)
  (emit "ret"))

(compile-program '(zero? -1))
