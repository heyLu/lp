module Stream where

import Date
import Date (Date)
import FormatDate (formatDate)
import Html (..)
import Html.Attributes as Attr
import List
import String

type alias Post = { title:String, content:String, created:Date }

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

referenceDate = date "2015-03-01T14:09"

viewDate : Date -> Html
viewDate d = let isoDate = formatDate "%Y-%m-%dT%H:%M:%SZ" d
             in time [Attr.title isoDate, Attr.datetime isoDate] [(text isoDate)]

viewPost : Post -> Html
viewPost post = div [] [
                 h3 [] [text post.title],
                 p [] [text post.content],
                 span [] [text "Written ", viewDate post.created]
                ]

main = div [] (List.map viewPost posts)