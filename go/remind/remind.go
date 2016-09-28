package main

import (
	"fmt"
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
