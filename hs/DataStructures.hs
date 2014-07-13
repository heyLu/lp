{-|
    Module: DataStructures
    Description: Having fun with (functional) data structures

    Inspired by:

      * Purely Functional Data Structures, by Chris Okasaki
      * <http://www.infoq.com/presentations/Functional-Data-Structures-in-Scala Functional Data Structures in Scala>, by Daniel Spiewak
-}
module DataStructures where

import Prelude hiding (concat, drop, last, length, reverse, take)

-- examples
sl = fromList [1..10] :: List Integer
bq = fromList [1..10] :: BankersQueue Integer
ft = fromList [1..10] :: FingerTree Integer

-- properties

-- rest nil == nil
-- (reverse . reverse) s == s
-- first [s] == last [s]

class Seq s where
    first :: s a -> Maybe a

    last :: s a -> Maybe a
    last s | isEmpty s = Nothing
    last s | isEmpty (rest s) = 
        case first s of
            Just x -> Just x
    last s = last $ rest s

    cons :: a -> s a -> s a
    nil :: s a

    rest :: s a -> s a

    length :: s a -> Integer
    length s | isEmpty s = 0
    length s = 1 + length (rest s)

    isEmpty  :: s a -> Bool

    -- additional interfaces

    append :: (Seq s) => a -> s a -> s a
    append x s | isEmpty s = cons x nil
    append x s =
        case first s of
             Just x -> cons x $ append x (rest s)

fromList :: (Seq s) => [a] -> s a
fromList [] = nil
fromList (x:xs) = cons x $ fromList xs

toList :: (Seq s) => s a -> [a]
toList s | isEmpty s = []
toList s =
    case first s of
        Just x -> x : toList (rest s)

butLast :: (Seq s) => s a -> s a
butLast s | isEmpty s = nil
butLast s =
    case first s of
        Just x -> cons x $ butLast (rest s)

drop :: (Seq s) => Integer -> s a -> s a
drop 0 s = s
drop _ s | isEmpty s = nil
drop n s = drop (n-1) $ rest s

take :: (Seq s) => Integer -> s a -> s a
take 0 _ = nil
take n s =
    case first s of
        Nothing -> nil
        Just x -> cons x $ take (n-1) (rest s)

concat :: (Seq s) => s a -> s a -> s a
concat l r | isEmpty l = r
concat l r | isEmpty r = l
concat l r =
    case first l of
        Just x -> cons x $ concat (rest l) r

reverse :: (Seq s) => s a -> s a
reverse s = rev s nil
    where rev l r | isEmpty l = r
          rev l r =
              case first l of
                  Nothing -> nil
                  Just x -> rev (rest l) $ cons x r

instance Seq [] where
    first [] = Nothing
    first (x:_) = Just x

    cons x xs = x:xs

    nil = []

    rest [] = []
    rest (_:xs) = xs

    isEmpty [] = True
    isEmpty _ = False

data List a = Nil | Cons a (List a) deriving Show

-- first and cons are O(1), everything else is O(n)
instance Seq List where
    first Nil = Nothing
    first (Cons x _) = Just x

    cons x l = Cons x l
    nil = Nil

    rest Nil = Nil
    rest (Cons _ l) = l

    isEmpty Nil = True
    isEmpty _ = False

class Queue q where
    enqueue :: a -> q a -> q a

    dequeue :: q a -> Maybe (a, q a)

-- fifo, remove from front, insert into rear
data BankersQueue a = BankersQueue {
                          frontSize :: Integer,
                          front :: List a,
                          rearSize :: Integer,
                          rear :: List a
                      } deriving (Show)

instance Queue BankersQueue where
    enqueue x (BankersQueue fs f rs r) = check $ BankersQueue fs f (rs + 1) (cons x r)

    dequeue (BankersQueue fs Nil         rs Nil) = Nothing
    -- not needed because of `check` invariant?
    --dequeue (BankersQueue fs Nil         rs r) = Just (x, check $ BankersQueue fs Nil (rs - 1) r')
    --    where (Just x) = last r
    --          r' = butLast r
    dequeue (BankersQueue fs (Cons x fr) rs r) =
        Just (x, check $ BankersQueue (fs - 1) fr rs r)

dequeueN :: Queue q => Integer -> q a -> Maybe (a, q a)
dequeueN 0 q = Nothing
dequeueN 1 q = dequeue q >>= \(x, q) -> return (x, q)
dequeueN n q = dequeue q >>= \(_, q') -> dequeueN (n-1) q'

check q@(BankersQueue fs f rs r) =
    if rs <= fs
    then q
    else BankersQueue (fs + rs) (f `concat` reverse r) 0 Nil

instance Seq BankersQueue where
    first (BankersQueue _ Nil         _ r) = last r
    first (BankersQueue _ (Cons  x _) _ _) = Just x

    -- O(1) amortized
    last (BankersQueue _ f _ Nil) = last f
    last (BankersQueue _ _ _ r) = first r

    cons x q = enqueue x q
    nil = BankersQueue 0 Nil 0 Nil

    rest q | isEmpty q = nil
    rest q =
        case dequeue q of
            Nothing -> nil
            Just (_, q') -> q'

    length (BankersQueue fs _ rs _) = fs + rs

    isEmpty (BankersQueue _ Nil _ Nil) = True
    isEmpty _ = False

data FingerTree a =
      Empty
    | Single a
    | Deep {
          ftPrefix :: Digit a,
          ftTree :: FingerTree (Node a),
          ftSuffix :: Digit a
      } deriving (Show)

data Digit a = One a | Two a a | Three a a a | Four a a a a deriving Show
data Node a = Node2 a a | Node3 a a a deriving Show

instance Seq Digit where
    first (One x) = Just x
    first (Two x _) = Just x
    first (Three x _ _) = Just x
    first (Four x _ _ _) = Just x

    last (One x) = Just x
    last (Two _ x) = Just x
    last (Three _ _ x) = Just x
    last (Four _ _ _ x) = Just x

    cons x (One a) = Two x a
    cons x (Two a b) = Three x a b
    cons x (Three a b c) = Four x a b c

    nil = error "can't be empty"

    rest (Two _ a) = One a
    rest (Three _ a b) = Two a b
    rest (Four _ a b c) = Three a b c

    isEmpty _ = False

    append x (One a) = Two a x
    append x (Two a b) = Three a b x
    append x (Three a b c) = Four a b c x

instance Seq FingerTree where
    first Empty = Nothing
    first (Single x) = Just x
    first (Deep p _ _) = first p

    last Empty = Nothing
    last (Single x) = Just x
    last (Deep _ _ s) = last s

    cons x Empty = Single x
    cons x (Single y) = Deep (One x) Empty (One y)
    cons x (Deep (Four a b c d) t s) =
        Deep (Two x a) (cons (Node3 b c d) t) s
    cons x (Deep p t s) = Deep (cons x p) t s

    nil = Empty

    rest Empty = Empty
    rest (Single _) = Empty
    rest (Deep (One _) t s) =
        case first t of
            Nothing ->
                case s of
                    One x -> Single x
                    Two x y -> Deep (One x) Empty (One y)
                    Three x y z -> Deep (Two x y) Empty (One z)
                    Four x y z w -> Deep (Three x y z) Empty (One w)
            Just (Node2 x y) -> Deep (Two x y) (rest t) s
            Just (Node3 x y z) -> Deep (Three x y z) (rest t) s
    rest (Deep p t s) = Deep (rest p) t s

    isEmpty Empty = True
    isEmpty _ = False

    append x Empty = Single x
    append x (Single y) = Deep (One x) Empty (One y)
    append x (Deep p t (Four a b c d)) =
        Deep p (append (Node3 a b c) t) (Two d x)
    append x (Deep p t s) = Deep p t (append x s)

instance Queue FingerTree where
    enqueue x ft = append x ft

    dequeue Empty = Nothing
    dequeue (Single x) = Just (x, Empty)
    dequeue ft =
        case last ft of
            Just x -> Just (x, rest ft) -- broken, we'd need a different version of rest
