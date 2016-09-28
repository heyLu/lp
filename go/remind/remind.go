package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	dates := []string{
		"today",
		"tomorrow",
		"tomorrow at 3am",
		"in 3 days",
		"in a month",
		"in 3 months",
		"next week",
		"next month",
		"in two weeks",
		"in 3 weeks",
		"2016-09-28",
		"3pm",
		"in 4 days at 10 pm",
	}
	for _, d := range dates {
		fmt.Printf("'%s': ", d)
		t, err := parseTime(d)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("%s\n", t)
		}
	}
}

func parseTime(s string) (time.Time, error) {
	var t time.Time
	now := time.Now().Round(time.Second)

	parts := strings.Fields(s)

	if len(parts) == 0 {
		return t, fmt.Errorf("empty date spec")
	}

	switch parts[0] {
	case "today":
		if len(parts) == 1 {
			return now, nil
		}
	case "tomorrow":
		if len(parts) == 1 {
			return now.AddDate(0, 0, 1), nil
		}
	case "in":
		if len(parts) == 3 {
			n, err := parseNumber(parts[1])
			if err != nil {
				return t, err
			}
			modifier, err := parseModifier(parts[2])
			if err != nil {
				return t, err
			}
			return modifier(n, now), nil
		}
	case "next":
		if len(parts) == 2 {
			modifier, err := parseModifier(parts[1])
			if err != nil {
				return t, err
			}
			return modifier(1, now), nil
		}
	default:
		return t, fmt.Errorf("unknown date spec '%s'", s)
	}

	return t, fmt.Errorf("unknown date spec '%s' (unexpected)", s)
}

func parseNumber(n string) (int, error) {
	switch n {
	case "a", "an", "one":
		return 1, nil
	case "two":
		return 2, nil
	case "three":
		return 3, nil
	case "four":
		return 4, nil
	case "five":
		return 5, nil
	case "six":
		return 6, nil
	case "seven":
		return 7, nil
	case "eight":
		return 8, nil
	case "nine":
		return 9, nil
	case "ten":
		return 10, nil
	default:
		return strconv.Atoi(n)
	}
}

func parseModifier(m string) (func(int, time.Time) time.Time, error) {
	switch m {
	case "second", "seconds":
		return durationModifier(time.Second), nil
	case "minute", "minutes":
		return durationModifier(time.Minute), nil
	case "hour", "hours":
		return durationModifier(time.Hour), nil
	case "day", "days":
		return dateModifier(0, 0, 1), nil
	case "week", "weeks":
		return dateModifier(0, 0, 7), nil
	case "month", "months":
		return dateModifier(0, 1, 0), nil
	case "year", "years":
		return dateModifier(1, 0, 0), nil
	default:
		return nil, fmt.Errorf("unknown modifier '%s'", m)
	}
}

func durationModifier(d time.Duration) func(int, time.Time) time.Time {
	return func(n int, t time.Time) time.Time {
		return t.Add(time.Duration(n) * d)
	}
}

func dateModifier(years, months, days int) func(int, time.Time) time.Time {
	return func(n int, t time.Time) time.Time {
		return t.AddDate(n*years, n*months, n*days)
	}
}
