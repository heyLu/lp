;; assembly resources:
;;  - https://en.wikipedia.org/wiki/X86_instruction_listings
(define (emit instr . args)
  (display "\t")
  (display (apply format instr args))
  (display "\n"))

(define label-counter 0)

(define (unique-label)
  (let ((c label-counter))
    (set! label-counter (+ c 1))
    (format "label_~d" c)))

(define wordsize 4)

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

(define (primcall-operand2 x)
  (caddr x))

(define (emit-compare)
  (emit "movl $0,  %eax")                ; zero %eax to put the result of the comparison into
  (emit "sete %al")                      ; set low byte of %eax to 1 if cmp succeeded
  (emit "sall $~a,  %eax" boolean-shift) ; construct correctly tagged boolean value
  (emit "xorl $~a, %eax" #b0011111))

(define (variable? x) (symbol? x))

(define (let? x)
  (and (list? x) (eq? (car x) 'let) (list? (cadr x)) (list? (caadr x))))

(define (if? x)
  (and (list? x) (eq? (car x) 'if)))

(define (emit-expr x si env)
  (cond
    ((immediate? x)
     (emit "movl $~a, %eax" (immediate-rep x)))
    ((variable? x)
     (emit "movl ~a(%rsp), %eax" (lookup x env)))
    ((let? x)
     (emit-let (bindings x) (body x) si env))
    ((if? x)
     (emit-if (test x) (conseq x) (altern x) si env))
    ((primcall? x) (emit-primitive-call x si env))))

(define (lookup x env)
  (cond
    ((null? env) (error 'lookup (format "~a is undefined" x)))
    ((eq? x (caar env)) (cadar env))
    (else (lookup x (cdr env)))))

(define (bindings let-expr) (cadr let-expr))
(define (body let-expr) (caddr let-expr))

(define (extend-env name val env)
  (cons (cons name (cons val '())) env))

(define (emit-let bindings body si env)
  (let f ((b* bindings) (new-env env) (si si))
    (cond
      ; if we're done with the bindings, emit the code for the body
      ((null? b*) (emit-expr body si new-env))
      ; otherwise, continue evaluating bindings in sequence
      (else
        (let ((b (car b*))) ; current binding
          ; emit code for current binding
          (emit-expr (rhs b) si env)
          ; move it onto the stack
          (emit "movl %eax, ~a(%rsp)" si)
          ; store current binding in env, continue generating bindings
          (f (cdr b*)
             (extend-env (lhs b) si new-env)
             (- si wordsize)))))))

(define (lhs binding) (car binding))
(define (rhs binding) (cadr binding))

(define (emit-if test conseq altern si env)
  (let ((L0 (unique-label)) (L1 (unique-label)))
    (emit-expr test si env)
    (emit "cmpl $~a, %eax" (immediate-rep #f))
    (emit "je ~a" L0)
    (emit-expr conseq si env)
    (emit "jmp ~a" L1)
    (emit-label L0)
    (emit-expr altern si env)
    (emit-label L1)))

(define (emit-label L)
  (display (format "~a:\n" L)))

(define (test if-expr)   (cadr if-expr))
(define (conseq if-expr) (caddr if-expr))
(define (altern if-expr) (cadddr if-expr))

(define (emit-primitive-call x si env)
  (case (primcall-op x)
    ((add1)
     (emit-expr (primcall-operand1 x) si env)
     (emit "addl $~a, %eax" (immediate-rep 1)))
    ((integer->char)
     (emit-expr (primcall-operand1 x) si env)
     (emit "shl $6, %eax")
     (emit "xorl $15, %eax"))
    ((char->integer)
     (emit-expr (primcall-operand1 x) si env)
     (emit "shrl $6, %eax"))
    ((zero?)
     (emit-expr (primcall-operand1 x) si env)
     (emit "cmpl $0,  %eax") ; x == 0
     (emit-compare))
    ((null?)
     (emit-expr (primcall-operand1 x) si env)
     (emit "cmpl $~a, %eax" #b00101111)
     (emit-compare))
    ((integer?)
     (emit-expr (primcall-operand1 x) si env)
     (emit "andl $~a, %eax" #b11)
     (emit-compare))
    ((char?)
     (emit-expr (primcall-operand1 x) si env)
     (emit "andl $~a, %eax" #b11111111)
     (emit "cmpl $~a, %eax" #b00001111)
     (emit-compare))
    ((boolean?)
     (emit-expr (primcall-operand1 x) si env)
     (emit "andl $~a, %eax" #b1111111)
     (emit "cmpl $~a, %eax" #b0011111)
     (emit-compare))
    ((+)
     (emit-expr (primcall-operand2 x) si env)
     (emit "movl %eax, ~a(%rsp)" si) ; move second arg on the stack
     (emit-expr (primcall-operand1 x) (- si wordsize) env)
     (emit "addl ~a(%rsp), %eax" si))
    ((cons)
     (emit-expr (primcall-operand1 x) si env)
     (emit "movl %eax, 0(%rsi)") ; set the car
     (emit-expr (primcall-operand2 x) si env)
     (emit "movl %eax, 4(%rsi)") ; set the cdr
     (emit "movq %rsi, %rax") ; rax = rsi | 1  (cons cell/pair tag)
     (emit "orq  $~a, %rax" #b001)
     (emit "addq $8,  %rsi")) ; bump rsi
    ((car)
     (emit-expr (primcall-operand1 x) si env)
     (emit "movl -1(%rax), %eax"))
    ((cdr)
     (emit-expr (primcall-operand1 x) si env)
     (emit "movl 3(%rax), %eax"))
    ((cddr)
     (emit-expr (primcall-operand1 x) si env)
     (emit "movl 3(%rax), %eax")
     (emit "movl 3(%rax), %eax"))
    ((cddar)
     (emit-expr (primcall-operand1 x) si env)
     (emit "movl 3(%rax), %eax")
     (emit "movl 3(%rax), %eax")
     (emit "movl -1(%rax), %eax"))
    ))

(define (compile-program x)
  (display ".globl scheme_entry\n\n")
  (display "scheme_entry:\n")

  (emit "movq %rax, %rsi") ; store pointer to heap memory
  (emit-expr x (- wordsize) '())
  (emit "ret"))
