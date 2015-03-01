module Post where

import Date
import Date (Date)
import Html (..)
import Util.Compare (compareBy, flipCompare)
import Util.Html (viewDate)

type alias Post = { title:String, content:String, created:Date }

compareByDate = compareBy (.created >> Date.toTime)
compareByDateReverse = flipCompare <| compareBy (.created >> Date.toTime)

view : Date -> Post -> Html
view ref post = div [] [
                 h3 [] [text post.title],
                 p [] [text post.content],
                 span [] [text "Written ", viewDate ref post.created]
                ]