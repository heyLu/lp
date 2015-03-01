module Stream where

import Date
import Date (Date)
import FormatDate (formatDate)
import Html (..)
import Html.Attributes as Attr
import List
import String
import PrettyDate (prettyDate)

import Post
import Post (Post)

date s = case (Date.fromString s) of
    Ok  d -> d
    Err e -> Date.fromTime 0

posts : List Post
posts = [{title = "Something else", content = "Well, I can say more than \"Hello\", I guess!", created = date "2015-03-01T14:03"},
         {title = "Hello, Other Things?", content = "There are other things?", created = date "2015-03-01T12:53:21"},
         {title = "Hello, World!", content = "This is my very first post!", created = date "2015-03-01T12:27:00"},
         {title = "Blog setup", content = "I guess I should post something now?", created = date "2014-12-24T10:57:03"},
         {title = "Ancient history", content = "Teenage angst!!!!", created = date "2009-06-07T02:54:29"}
        ]

flipOrder : Order -> Order
flipOrder o = case o of
                LT -> GT
                EQ -> EQ
                GT -> LT

flipCompare : (a -> a -> Order) -> a -> a -> Order
flipCompare compare' a b = flipOrder <| compare' a b

compareBy : (a -> comparable) -> a -> a -> Order
compareBy f a b = compare (f a) (f b)

sortByDate = List.sortBy (.created >> Date.toTime)
sortByDateReverse = List.sortWith (flipCompare <| compareBy (.created >> Date.toTime))

referenceDate = date "2015-03-01T14:09"

main = div [] (List.map (Post.view referenceDate) (sortByDateReverse posts))