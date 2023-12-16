package main

import (
	"fmt"
	"time"
)

var _ TasksStore = &memoryStore{}

type memoryStore struct {
	tasks  []Task
	events []Event
}

func (ms *memoryStore) Tasks() ([]Task, error) {
	return ms.tasks, nil
}

func (ms *memoryStore) FindTask(id string) (*Task, error) {
	for _, task := range ms.tasks {
		if task.ID != id {
			continue
		}

		return &task, nil
	}

	return nil, fmt.Errorf("could not find task %s", id)

}

func (ms *memoryStore) ChangeTaskState(id string, state TaskState) error {
	found := false
	for i, task := range ms.tasks {
		if task.ID != id {
			continue
		}

		ms.tasks[i].State = state

		found = true
		break
	}

	if !found {
		return fmt.Errorf("could not find task %s", id)
	}

	return nil
}

func (ms *memoryStore) Events() ([]Event, error) {
	return ms.events, nil
}

var exampleTasks = generateIDs([]Task{
	{Icon: "🧹", State: TaskNotDone, Description: "cleaned a bit"},
	{Icon: "🌪️", State: TaskDone, Description: "vaccumed this week"},
	{Icon: "✨", State: TaskNotDone, Description: "kitchen counter clean"},
	{Icon: "🎶", State: TaskStarted, Description: "played some music"},
	{Icon: "📚", State: TaskNotDone, Description: "practiced some"},
	{Icon: "🍎", State: TaskNotDone, Description: "ate some fruit"},
	{Icon: "🍵", State: TaskDone, Description: "got hydrated"},
	{Icon: "👚", State: TaskDone, Description: "washed clothes this week"},
	{Icon: "🐑", State: TaskNotDone, Description: "washed sheets this month"},
})

func generateIDs(tasks []Task) []Task {
	for i := range tasks {
		tasks[i].ID = fmt.Sprintf("%d", i+1)
	}
	return tasks
}

var exampleEvents = []Event{
	{
		Icon:          "⏳",
		Date:          time.Date(time.Now().Year(), time.December, 31, 23, 59, 59, 0, time.Now().Location()),
		ReferenceDate: time.Date(time.Now().Year(), time.January, 01, 00, 00, 0, 0, time.Now().Location()),
	},
	{Icon: "🌲", Date: timeMustParse(time.DateOnly, "2023-12-24"), ReferenceDate: timeMustParse(time.DateOnly, "2023-12-15")},
	{Icon: "🕊️", Date: timeMustParse(time.DateOnly, "2024-01-08"), ReferenceDate: timeMustParse(time.DateOnly, "2023-12-15")},
	{Icon: "🧬", Date: timeMustParse(time.DateOnly, fmt.Sprintf("%d-01-01", 1990+81)), ReferenceDate: timeMustParse(time.DateOnly, "1990-01-01")},
}

func timeMustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}
