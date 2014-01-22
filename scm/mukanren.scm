; µKanren - http://webyrd.net/scheme-2013/papers/HemannMuKanren2013.pdf

; to use this in guile: ,import (rnrs lists)

(define empty-state '(() . 0))

(define (var x)  (vector x))
(define (var? x) (vector? x))
(define (var=? x y) (= (vector-ref x 0) (vector-ref y 0)))

(define (walk u s)
  (let ((pr (and (var? u) (assp (lambda (v) (var=? u v)) s))))
    (if pr (walk (cdr pr) s) u))) ; recursion b/c vars can refer to other vars?

(define (ext-s x v s) `((,x . ,v) . ,s))

(define (≡ u v)
  (lambda (s/c)
    (let ((s (unify u v (car s/c))))
      (if s (unit `(,s . ,(cdr s/c))) mzero))))

(define (unit s/c) (cons s/c mzero))
(define mzero '())

(define (unify u v s)
  (let ((u (walk u s)) (v (walk v s)))
    (cond
      ((and (var? u) (var? v) (var=? u v)) s)
      ((var? u) (ext-s u v s))
      ((var? v) (ext-s v u s))
      ((and (pair? u) (pair? v))
       (let ((s (unify (car u) (car v) s)))
         (and s (unify (cdr u) (cdr v) s))))
      (else (and (eqv? u v) s)))))

(define (call/fresh f)
  (lambda (s/c)
    (let ((c (cdr s/c)))
      ((f (var c)) `(,(car s/c) . ,(+ c 1))))))

(define (disj g1 g2) (lambda (s/c) (mplus (g1 s/c) (g2 s/c))))
(define (conj g1 g2) (lambda (s/c) (bind (g1 s/c) g2)))

;; 4.1 finite depth first search

(define (mplus $1 $2)
  (cond
    ((null? $1) $2)
    (else (cons (car $1) (mplus (cdr $1) $2)))))

(define (bind $ g)
  (cond
    ((null? $) mzero)
    (else (mplus (g (car $)) (bind (cdr $) g)))))

; is this a concept or a real example?
(define (fives x) (disj (≡ x 5) (fives x)))
;((call/fresh fives) empty-state) ;=> infinite loop

((call/fresh (lambda (v) (disj (≡ v 0) (≡ v 42)))) empty-state)

;; 4.2 infinite streams

(define (mplus $1 $2)
  (cond
    ((null? $1) $2)
    ((procedure? $1) (lambda () (mplus ($1) $2)))
    (else (cons (car $1) (mplus (cdr $1) $2)))))

(define (bind $ g)
  (cond
    ((null? $) mzero)
    ((procedure? $) (lambda () (bind ($) g)))
    (else (mplus (g (car $)) (bind (cdr $) g)))))

(define (fives x)
  (disj (≡ x 5) (lambda (s/c) (lambda () ((fives x) s/c)))))

((call/fresh fives) empty-state)

;; 4.3 interleaved streams

(define (mplus $1 $2)
  (cond
    ((null? $1) $2)
    ((procedure? $1) (lambda () (mplus $2 ($1))))
    (else (cons (car $1) (mplus (cdr $1) $2)))))

(define (sixes x)
  (disj (≡ x 6) (lambda (s/c) (lambda () ((sixes x) s/c)))))

(define fives-and-sixes
  (call/fresh (lambda (x) (disj (fives x) (sixes x)))))

;; 5 utilities

#|
(define (expand-deferred s)
  (if (procedure? s)
    (expand-deferred (s))
    s))

(define (expand-once s)
  (if (procedure? s)
    (s)
    s))

(define (expand n s)
  (if (and (> n 1) (not (null? s)))
    `(,(expand-deferred (car s)) . ,(expand (- n 1) (expand-once (cdr s))))
    (expand-deferred s)))

(expand 5 (fives-and-sixes empty-state))

(define (run n g)
  (expand n (g empty-state)))
|#

(define-syntax Zzz
  (syntax-rules ()
    ((_ g) (lambda (s/c) (lambda () (g s/c))))))

(define-syntax conj+
  (syntax-rules ()
    ((_ g) (Zzz g))
    ((_ g0 g ...) (conj (Zzz g0) (conj+ g ...)))))

(define-syntax disj+
  (syntax-rules ()
    ((_ g) (Zzz g))
    ((_ g0 g ...) (disj (Zzz g0) (disj+ g ...)))))

(define-syntax conde
  (syntax-rules ()
    ((_ (g0 g ...) ...) (disj+ (conj+ g0 g ...) ...))))

(define-syntax fresh
  (syntax-rules ()
    ((_ () g0 g ...) (conj+ g0 g ...))
    ((_ (x0 x ...) g0 g ...)
     (call/fresh (lambda (x0) (fresh (x ...) g0 g ...))))))

;(run 1 (fresh (x y z) (≡ x 7) (≡ y 8) (disj (≡ x z) (≡ y z))))

;; 5.2 from streams to lists

(define (pull $) (if (procedure? $) (pull ($)) $))

(define (take-all $)
  (let (($ (pull $)))
    (if (null? $) '() (cons (car $) (take-all (cdr $))))))

(define (take n $)
  (if (zero? n) '()
    (let (($ (pull $)))
      (cond
        ((null? $) '())
        (else (cons (car $) (take (- n 1) (cdr $))))))))

(take 10 ((fresh (x y z) (≡ x 7) (≡ y 8) (disj+ (≡ x z) (≡ y z) (≡ x y))) empty-state))

;; 5.3 recovering reification

(define (mK-reify s/c*)
  (map reify-state/1st-var s/c*))

(define (reify-state/1st-var s/c)
  (let ((v (walk* (var 0) (car s/c))))
    (walk* v (reify-s v '()))))

(define (reify-s v s)
  (let ((v (walk v s)))
    (cond
      ((var? v)
       (let ((n (reify-name (length s))))
         (cons `(,v . ,n) s)))
      ((pair? v) (reify-s (cdr v) (reify-s (car v) s)))
      (else s))))

(define (reify-name n)
  (string->symbol
    (string-append "_." (number->string n))))

(define (walk* v s)
  (let ((v (walk v s)))
    (cond
      ((var? v) v)
      ((pair? v) (cons (walk* (car v) s)
                       (walk* (cdr v) s)))
      (else v))))

(mK-reify (take 10 ((fresh (x y z) (≡ x 7) (≡ y 8) (disj+ (≡ x z) (≡ y z) (≡ x y))) empty-state)))

(mK-reify (take 100 (fives-and-sixes empty-state)))

;; 5.4 recovering the interface to scheme

(define (call/empty-state g) (g empty-state))

(define-syntax run
  (syntax-rules ()
    ((_ n (x ...) g0 g ...)
     (mK-reify (take n (call/empty-state
                         (fresh (x ...) g0 g ...)))))))

(define-syntax run*
  (syntax-rules ()
    ((_ (x ...) g0 g ...)
     (mK-reify (take-all (call/empty-state
                           (fresh (x ...) g0 g ...)))))))

(run 10 (y x z) (≡ x 7) (disj+ (≡ y 8) (≡ y 18) (≡ y 28)) (≡ z 9))