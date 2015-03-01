module Util.Compare where

flipOrder : Order -> Order
flipOrder o = case o of
                LT -> GT
                EQ -> EQ
                GT -> LT

flipCompare : (a -> a -> Order) -> a -> a -> Order
flipCompare compare' a b = flipOrder <| compare' a b

compareBy : (a -> comparable) -> a -> a -> Order
compareBy f a b = compare (f a) (f b)
