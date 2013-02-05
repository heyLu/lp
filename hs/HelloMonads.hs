{-
 - Hello, Monads!
 -
 - Exploring what Monads mean (in Haskell).
 -}
import Prelude hiding (Maybe(..), print, putStr, putStrLn, getLine, getContent)

-- define Maybe here so we can implement our own monad
data Maybe a = Just a | Nothing deriving (Show, Eq, Ord)

instance Monad Maybe where
    (Just x) >>= action = action x
    Nothing  >>= action = Nothing

    -- apparently optional:
    -- maybeA   >>  maybeB = maybeA >>= \_ -> maybeB

    return x            = Just x

    -- FIXME: fail?
    -- fail _ = Nothing

maybeAdd m1 m2 = m1 >>= \x -> m2 >>= \y -> return (x + y)

maybeAdd' m1 m2 = do
    x <- m1
    y <- m2
    return $ x + y

type InputState = String
type OutputState = String
-- * Model IO as a function of the previous state to a new state and a result.
--
-- If you look closely (i.e. execute `:i IO` in ghci), you'll notice
-- that this is a simplification of the real IO type.
-- Of course it is missing file handles, actual output and other
-- gimmicks, but it gives an idea of what IO does (and why).
data WeirdIO a = WeirdIO ((InputState, OutputState) -> ((InputState, OutputState), a))

executeWeirdIO :: InputState -> WeirdIO a -> ((InputState, OutputState), a)
executeWeirdIO input (WeirdIO f) = {-snd . fst $-} f (input, "")

instance Monad WeirdIO where
    (WeirdIO changeState) >>= action = WeirdIO $ \state ->
        let ((i, o), x) = changeState state
            (WeirdIO f) = action x
        in f ((i, o))

    return x = WeirdIO $ \state -> (state, x)

print :: Show a => a -> WeirdIO ()
print x = putStrLn $ show x
putStr, putStrLn :: String -> WeirdIO ()
putStr s = WeirdIO $ \(i, o) -> ((i, o ++ s), ())
putStrLn s = putStr $ s ++ "\n"

getChar :: WeirdIO Char
getChar = WeirdIO $ \(c:cs, o) -> ((cs, o), c)
getLine, getContent :: WeirdIO String
getLine = WeirdIO $ \(i, o) -> (((dropUntil (== '\n') i), o), takeWhile (/= '\n') i)
getContent = WeirdIO $ \(i, o) -> (("", o), i)

takeUntil :: (a -> Bool) -> [a] -> [a]
takeUntil p [] = []
takeUntil p (x:xs) = if (p x) then [x] else x : takeUntil p xs

dropUntil p [] = []
dropUntil p (x:xs) = if (p x) then xs else dropUntil p xs
