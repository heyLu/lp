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
augmentedAssignment op = \name value -> Imperative $ \env ->
    (M.insertWith (flip op) name value env, (fromJust $ M.lookup name env) `op` value)
(.+=) = augmentedAssignment (+)
(.-=) = augmentedAssignment (-)
(.*=) = augmentedAssignment (*)
(./=) = augmentedAssignment (/)

binaryOp op = \var1 var2 -> Imperative $ \env ->
    (env, fromJust $ op <$> M.lookup var1 env <*> M.lookup var2 env)

(.+) = binaryOp (+)
(.-) = binaryOp (-)
(.*) = binaryOp (*)
(./) = binaryOp (/)
