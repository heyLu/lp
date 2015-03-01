module Post where

import Date
import Date (Date)
import List
import Html (..)
import Util.Compare (compareBy, flipCompare)
import Util.Html (viewDate)
import String

type alias Post = { title:String, content:String, created:Date }

compareByDate = compareBy (.created >> Date.toTime)
compareByDateReverse = flipCompare <| compareBy (.created >> Date.toTime)

view : Date -> Post -> Html
view ref post =
    let paragraphs = List.filter (not << String.isEmpty) <| String.split "\n" post.content
    in div [] [
            h3 [] [text post.title],
            section [] (List.map (\par -> p [] [text par]) paragraphs),
            span [] [text "Written ", viewDate ref post.created]
           ]