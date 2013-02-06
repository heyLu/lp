{-
 - Imperative programming in Haskell.
 -
 - Another monadic exercise.
 -}
module ImperativeProgramming where

import Data.Maybe (fromJust)
import qualified Data.Map as M
import Control.Monad (liftM)
import Control.Applicative ((<$>), (<*>))

example = do
    "x" .= 3
    "x" .+= 4
    "y" .= (-1)
    "x" .+ "y"

data Imperative a b = Imperative (M.Map String a -> (M.Map String a, b))

executeImperative env (Imperative program) = program env

instance Monad (Imperative a) where
    Imperative lastOperation >>= action = Imperative $ \env ->
        let (env', value) = lastOperation env
            (Imperative f) = action value
        in f env'

    return x = Imperative $ \env -> (env, x)

name .= value = Imperative $ \env -> (M.insert name value env, value)
name .+= value = Imperative $ \env ->
    (M.insertWith (+) name value env, value + (fromJust $M.lookup name env))
binaryOp op = \var1 var2 -> Imperative $ \env ->
    (env, fromJust $ op <$> M.lookup var1 env <*> M.lookup var2 env)

(.+) = binaryOp (+)
(.-) = binaryOp (-)
(.*) = binaryOp (*)
(./) = binaryOp (/)
