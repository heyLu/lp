module Util.Html where

import Date (Date)
import Html (..)
import Html.Attributes as Attr
import FormatDate (formatDate)
import PrettyDate (prettyDate)

viewDate : Date -> Date -> Html
viewDate ref d = let dateString = prettyDate ref d
                     isoDate = formatDate "%Y-%m-%dT%H:%M:%SZ" d
                 in time [Attr.title isoDate, Attr.datetime isoDate] [(text dateString)]