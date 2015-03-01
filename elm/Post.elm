module Post where

import Date
import Date (Date)
import List
import Html (..)
import Html.Attributes as Attr
import Html.Events (on, targetValue, onClick)
import Util.Compare (compareBy, flipCompare)
import Util.Html (viewDate)
import String
import Signal
import LocalChannel as LC
import LocalChannel (LocalChannel)

type alias Post = { title:String, content:String, created:Date }

empty : Post
empty = {title = "", content = "", created = Date.fromTime 0}

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

-- an editing interface ...

type alias Model = Post

type Action = NoOp
            | UpdateTitle String
            | UpdateContent String

update : Action -> Model -> Model
update action model =
    case action of
      NoOp -> model

      UpdateTitle title -> { model | title <- title }

      UpdateContent content -> { model | content <- content }

type alias Context = { actions : LocalChannel Action
                     , publish : LocalChannel () }

editingView : Context -> Model -> Html
editingView context model = div [] [
                             input [Attr.type' "text", Attr.placeholder "title",
                                    on "input" targetValue (LC.send context.actions << UpdateTitle)] [text model.title],
                             textarea [Attr.placeholder "write something!",
                                       on "input" targetValue (LC.send context.actions << UpdateContent)] [text model.content],
                             button [onClick (LC.send context.publish ())] [text "publish"]
                            ]