package main

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	now := time.Now()

	examples := []struct {
		Input  string
		Result time.Time
	}{
		{"today", now},
		{"tomorrow", now.AddDate(0, 0, 1)},
		{"tomorrow at 3am", truncateHours(now).AddDate(0, 0, 1).Add(3 * time.Hour)},
		{"in 3 days", now.AddDate(0, 0, 3)},
		{"in a month", now.AddDate(0, 1, 0)},
		{"in 3 months", now.AddDate(0, 3, 0)},
		{"next week", now.AddDate(0, 0, 7)},
		{"next month", now.AddDate(0, 1, 0)},
		{"in two weeks", now.AddDate(0, 0, 2*7)},
		{"in 3 weeks", now.AddDate(0, 0, 3*7)},
		//"2016-09-28",
		//"3pm",
		{"in 4 days at 10 pm", truncateHours(now).AddDate(0, 0, 4).Add(22 * time.Hour)},
	}

	for _, example := range examples {
		res, err := parseTimeRelative(example.Input, now)
		if err != nil {
			t.Errorf("parse '%s': unexpected error: %s", err)
		}

		if !res.Equal(example.Result) {
			t.Errorf("parse '%s': %s != %s", example.Input, res, example.Result)
		}
	}
}
