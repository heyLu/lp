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
posts = [{title = "Fancy post", content = "Post may have multiple lines now.\nWhat freedom!\n\n\nThat's weird, though...", created = date "2015-03-01T15:18"},
         {title = "Something else", content = "Well, I can say more than \"Hello\", I guess!", created = date "2015-03-01T14:03"},
         {title = "Hello, Other Things?", content = "There are other things?", created = date "2015-03-01T12:53:21"},
         {title = "Hello, World!", content = "This is my very first post!", created = date "2015-03-01T12:27:00"},
         {title = "Blog setup", content = "I guess I should post something now?", created = date "2014-12-24T10:57:03"},
         {title = "Ancient history", content = "Teenage angst!!!!", created = date "2009-06-07T02:54:29"}
        ]

referenceDate = date "2015-03-01T15:20"

main = div [] (List.map (Post.view referenceDate) (List.sortWith Post.compareByDateReverse posts))