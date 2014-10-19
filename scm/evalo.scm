(load "/home/lu/t/miniKanren/mk.scm")

; a simple lambda calculus in miniKanren

; x - variables
; (lambda (x) x) - abstraction
; (e1 e2) - application

(define lookupo
  (lambda (x env val)
    (fresh (y v rest)
      (== `((,y . ,v) . ,rest) env)
      (conde
       [(== x y) (== v val)]
       [(=/= x y) (lookupo x rest val)]))))

(define evalmanyo
  (lambda (l env r)
    (conde
     [(== l '()) (== l r)]
     [(fresh (h t h^ t^)
        (== l `(,h  . ,t))
        (== r `(,h^ . ,t^))
        (evalo h env h^)
        (evalmanyo t env t^))])))

(define evalo
  (lambda (expr env val)
    (conde
     [(symbolo expr)
      (lookupo expr env val)]
     [(fresh (e)
        (== `(quote ,e) expr)
        (== e val))]
     [(fresh (els)
        (== `(list . ,els) expr)
        (evalmanyo els env val))]
     [(fresh (x body)
        (== `(lambda (,x) ,body) expr)
        (== `(,expr in ,env) val))]
     [(fresh (e1 e2 x body env^ arg)
        (== `(,e1 ,e2) expr)
        (evalo e1 env `((lambda (,x) ,body) in ,env^))
        (evalo e2 env arg)
        (evalo body `((,x . ,arg) . ,env^) val))])))

(define quineo
  (lambda (q)
    (evalo q '() q)))

(define twineo
  (lambda (p q)
    (=/= p q)
    (evalo p '() q)
    (evalo q '() p)))
