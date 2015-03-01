module PrettyDate (prettyDate) where

{-| Relative time formatting.

Supports both dates in the past and in the future.
 -}

import Date
import Date (Date)

{-| Format the second date relative to the first one. -}
prettyDate : Date -> Date -> String
prettyDate reference date =
    let referenceSeconds = Date.toTime reference
        dateSeconds = Date.toTime date
        diffSeconds = referenceSeconds - dateSeconds
    in format << toRelativeTime <| diffSeconds

type RelativeUnit = JustNow
                  | Minute
                  | Hour
                  | Day
                  | Week
                  | Month
                  | Year

type RelativeTime = Past RelativeUnit Int
                  | Future RelativeUnit Int

stringForm relativeUnit =
    case relativeUnit of
      Minute  -> ("a minute", "minutes")
      Hour    -> ("an hour", "hours")
      Day     -> ("a day", "days")
      Week    -> ("a week", "weeks")
      Month   -> ("a month", "months")
      Year    -> ("a year", "years")

format : RelativeTime -> String
format relativeTime =
    case relativeTime of
      Past   JustNow _ -> "moments ago"
      Future JustNow _ -> "in a moment"
      Past   unit n    -> format' ""    " ago" unit n
      Future unit n    -> format' "in " ""     unit n

format' prefix suffix unit n =
    let (singular, plural) = stringForm unit
    in if n == 1
       then prefix ++ singular
       else prefix ++ toString n ++ " " ++ plural ++ suffix

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
        (relativeUnit, n) =
            if | diff < minute -> (JustNow, 0)
               | diff < hour   -> (Minute, round (diff / minute))
               | diff < day    -> (Hour, round (diff / hour))
               | diff < week   -> (Day, round (diff / day))
               | diff < month  -> (Week, round (diff / week))
               | diff < year   -> (Month, round (diff / month))
               | otherwise     -> (Year, round (diff / year))
    in if inThePast
       then Past relativeUnit n
       else Future relativeUnit n