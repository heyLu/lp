package main

import (
	"fmt"
	"time"
)

var _ TasksStore = &memoryStore{}

type memoryStore struct{}

func (ms *memoryStore) Tasks() ([]Task, error) {
	return []Task{
		{Icon: "🧹", State: TaskNotDone, Description: "cleaned a bit"},
		{Icon: "🌪️", State: TaskDone, Description: "vaccumed this week"},
		{Icon: "✨", State: TaskNotDone, Description: "kitchen counter clean"},
		{Icon: "🎶", State: TaskNotDone, Description: "played some music"},
		{Icon: "📚", State: TaskNotDone, Description: "practiced some"},
		{Icon: "👚", State: TaskDone, Description: "washed clothes this week"},
		{Icon: "🐑", State: TaskNotDone, Description: "washed sheets this month"},
	}, nil
}

func (ms *memoryStore) Events() ([]Event, error) {
	return []Event{
		{
			Icon:          "⏳",
			Date:          time.Date(time.Now().Year(), time.December, 31, 23, 59, 59, 0, time.Now().Location()),
			ReferenceDate: time.Date(time.Now().Year(), time.January, 01, 00, 00, 0, 0, time.Now().Location()),
		},
		{Icon: "🌲", Date: timeMustParse(time.DateOnly, "2023-12-24"), ReferenceDate: timeMustParse(time.DateOnly, "2023-12-15")},
		{Icon: "🕊️", Date: timeMustParse(time.DateOnly, "2024-01-08"), ReferenceDate: timeMustParse(time.DateOnly, "2023-12-15")},
		{Icon: "🧬", Date: timeMustParse(time.DateOnly, fmt.Sprintf("%d-01-01", 1990+81)), ReferenceDate: timeMustParse(time.DateOnly, "1990-01-01")},
	}, nil
}

func timeMustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}
