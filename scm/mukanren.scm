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

(define (expand n s)
  (if (> n 1)
    (cond
      ((null? s) s)
      ((procedure? (cdr s)) `(,(car s) . ,(expand (- n 1) ((cdr s)))))
      (else `(,(car s) . ,(expand (- n 1) (cdr s)))))
    s))

(expand 5 (fives-and-sixes empty-state))

(define (run n g)
  (expand n (g empty-state)))