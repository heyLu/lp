-- An evaluator for the untyped calculus of booleans and numbers, as
-- presented in chapter 3 of Types and Programming Languages.
module Arith where

import Prelude hiding (pred, succ)

data Expr = ATrue
          | AFalse
          | AIfThenElse Expr Expr Expr
          | AZero
          | ASucc Expr
          | APred Expr
          | AIsZero Expr
          deriving (Show)

true, false :: Expr
true  = ATrue
false = AFalse

ifthenelse :: Expr -> Expr -> Expr -> Expr
ifthenelse = AIfThenElse

zero :: Expr
zero = AZero

succ, pred :: Expr -> Expr
succ = ASucc
pred = APred

iszero :: Expr -> Expr
iszero = AIsZero

eval :: Expr -> Expr
eval ATrue = ATrue
eval AFalse = AFalse
eval (AIfThenElse test thenExpr elseExpr) =
    case eval test of
        ATrue -> eval thenExpr
        AFalse -> eval elseExpr
        _ -> error "test must return true or false"
eval AZero = AZero
eval (ASucc e) =
    case eval e of
        AZero -> ASucc AZero
        (ASucc f) -> ASucc f
        (APred f) -> f
        _ -> error "succ must operate on a number"
eval (APred e) =
    case eval e of
        AZero -> APred AZero
        (ASucc f) -> f
        (APred f) -> APred f
        _ -> error "pred must operate on a number"
eval (AIsZero e) =
    case eval e of
        AZero -> ATrue
        (ASucc _) -> AFalse
        (APred _) -> AFalse
        _ -> error "iszero must operate on a number"
