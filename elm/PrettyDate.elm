module PrettyDate (prettyDate) where

import Date
import Date (Date)

prettyDate : Date -> Date -> String
prettyDate reference date =
    let referenceSeconds = Date.toTime reference
        dateSeconds = Date.toTime date
        diffSeconds = referenceSeconds - dateSeconds
    in format << toRelativeTime <| diffSeconds

type RelativeTime = JustNow
                  | Minute Int
                  | Hour Int
                  | Day Int
                  | Week Int
                  | Month Int
                  | Year Int
                  | Future RelativeTime

format : RelativeTime -> String
format relativeTime =
    case relativeTime of
      JustNow -> "moments ago"
      Minute 1 -> "a minute ago"
      Minute n -> toString n ++ " minutes ago"
      Hour 1 -> "hour ago"
      Hour n -> toString n ++ " hours ago"
      Day 1 -> "a day ago"
      Day n -> toString n ++ " years ago"
      Week 1 -> "a week ago"
      Week n -> toString n ++ " weeks ago"
      Month 1 -> "a month ago"
      Month n -> toString n ++ " months ago"
      Year 1 -> "a year ago"
      Year n -> toString n ++ "years ago"
      Future (JustNow) -> "in a moment"
      Future (Minute n) -> "in " ++ toString n ++ " hours"
      Future (Hour n) -> "in " ++ toString n ++ " hours"
      Future (Day n) -> "in " ++ toString n ++ " days"
      Future (Week n) -> "in " ++ toString n ++ " weeks"
      Future (Month n) -> "in " ++ toString n ++ " months"
      Future (Year n) -> "in " ++ toString n ++ " years"
      _ -> "???"

toRelativeTime : Float -> RelativeTime
toRelativeTime origDiff =
    let minute = 60 * 1000
        hour = 60 * minute
        day = 24 * hour
        week = 7 * day
        month = 31 * day
        year = 365 * day
        roundToString = toString << ceiling
        inThePast = origDiff > 0
        diff = abs origDiff
        relativeTime =
            if | diff < minute -> JustNow
               | diff < hour   -> Minute <| round (diff / minute)
               | diff < day    -> Hour <| round (diff / hour)
               | diff < week   -> Day <| round (diff / day)
               | diff < month  -> Week <| round (diff / week)
               | diff < year   -> Month <| round (diff / month)
               | otherwise     -> Year <| round (diff / year)
    in if inThePast
       then relativeTime
       else Future relativeTime