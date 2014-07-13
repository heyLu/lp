module DataStructures where

import Prelude hiding (concat, drop, last, length, reverse, take)

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

fromList :: (Seq s) => [a] -> s a
fromList [] = nil
fromList (x:xs) = cons x $ fromList xs

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

append :: (Seq s) => a -> s a -> s a
append x s | isEmpty s = cons x nil
append x s =
    case first s of
        Just x -> cons x $ append x (rest s)

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

-- fifo, remove from front, insert into rear
data BankersQueue a = BankersQueue {
                          frontSize :: Integer,
                          front :: List a,
                          rearSize :: Integer,
                          rear :: List a
                      } deriving (Show)

enqueue x (BankersQueue fs f rs r) = check $ BankersQueue fs f (rs + 1) (cons x r)

dequeue (BankersQueue fs Nil         rs Nil) = Nothing
dequeue (BankersQueue fs Nil         rs r) = Just (x, check $ BankersQueue fs Nil (rs - 1) r')
    where (Just x) = last r
          r' = butLast r
dequeue (BankersQueue fs (Cons x fr) rs r) =
    Just (x, check $ BankersQueue (fs - 1) fr rs r)

dequeueN 0 q = Nothing
dequeueN 1 q = dequeue q >>= \(x, _) -> return x
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

-- examples
sl = fromList [1..10] :: List Integer
bq = fromList [1..10] :: BankersQueue Integer
