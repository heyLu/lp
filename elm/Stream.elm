module Stream where

import Date
import Date (Date)
import Html (..)
import Html.Attributes as Attr
import Html.Events (onClick)
import List
import String
import Time
import Signal
import LocalChannel

import Post
import Post (Post)

type alias Model =
    { isEditing     : Bool
    , newPost       : Post.Model
    , posts         : List Post
    , referenceDate : Date
    }

type Action = NoOp
            | CreatePost
            | AddPost ()
            | CancelPost
            | EditPost Post.Action
            | UpdateDate Date

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

-- view

view model = div [] [
              if model.isEditing
              then Post.editingView (Post.Context
                                             (LocalChannel.create EditPost updates)
                                             (LocalChannel.create AddPost updates)) Post.empty
              else button [onClick (Signal.send updates CreatePost)] [text "create new post"],
              div [] (List.map (Post.view model.referenceDate) (List.sortWith Post.compareByDateReverse model.posts))
             ]

-- wiring it all up

initialModel = { isEditing = False, newPost = Post.empty, posts = posts, referenceDate = date "0" }

update : Action -> Model -> Model
update action model =
    case action of
      NoOp -> model

      CreatePost -> { model | isEditing <- True }

      AddPost _ -> { model |
                      isEditing <- False,
                      newPost <- Post.empty,
                      posts <- let p = model.newPost
                               in { p | created <- model.referenceDate } :: model.posts }

      CancelPost -> { model | isEditing <- False, newPost <- Post.empty }

      UpdateDate d -> { model | referenceDate <- d }

      EditPost action' -> { model | newPost <- Post.update action' model.newPost }

model : Signal Model
model = Signal.foldp update initialModel <| Signal.merge (Signal.subscribe updates) (Signal.map UpdateDate currentDate)

updates : Signal.Channel Action
updates = Signal.channel NoOp

currentDate : Signal Date
currentDate = Signal.map Date.fromTime <| Time.every Time.second

main = Signal.map view model