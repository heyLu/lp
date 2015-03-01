module PrettyDate (prettyDate) where

import Date
import Date (Date)

prettyDate : Date -> Date -> String
prettyDate reference date =
    let referenceSeconds = Date.toTime reference
        dateSeconds = Date.toTime date
        diffSeconds = referenceSeconds - dateSeconds
    in format diffSeconds

format : Float -> String
format diff =
    let minute = 60 * 1000
        hour = 60 * minute
        day = 24 * hour
        week = 7 * day
        month = 31 * day
        year = 365 * day
        roundToString = toString << ceiling
    in if | diff < minute -> "moments ago"
          | diff < hour   -> roundToString (diff / minute) ++ " minutes ago"
          | diff < day    -> roundToString (diff / hour)   ++ " hours ago"
          | diff < week   -> roundToString (diff / day)    ++ " days ago"
          | diff < month  -> roundToString (diff / week)   ++ " weeks ago"
          | diff < year   -> roundToString (diff / month)  ++ " months ago"
          | otherwise     -> roundToString (diff / year)   ++ " years ago"