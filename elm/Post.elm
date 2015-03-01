module Post where

import Date (Date)
import Util.Html (viewDate)
import Html (..)

type alias Post = { title:String, content:String, created:Date }

view : Date -> Post -> Html
view ref post = div [] [
                 h3 [] [text post.title],
                 p [] [text post.content],
                 span [] [text "Written ", viewDate ref post.created]
                ]