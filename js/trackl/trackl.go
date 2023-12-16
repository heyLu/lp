package main

import (
	"html/template"
	"log"
	"net/http"
	"slices"
	"time"
)

type Task struct {
	Icon        string
	Description string
	State       TaskState
}

type TaskState string

var (
	TaskNotDone TaskState = "not-done"
	TaskStarted TaskState = "started"
	TaskDone    TaskState = "done"
)

type Event struct {
	Icon          string
	Date          time.Time
	ReferenceDate time.Time
}

func (e Event) PercentDone() float64 {
	available := float64(e.Date.Sub(e.ReferenceDate) / (24 * time.Hour))
	left := float64(e.Date.Sub(time.Now()) / (24 * time.Hour))
	return (float64(available-left) / float64(available)) * 100
}

func (e Event) DaysLeft() int {
	return int(e.Date.Sub(time.Now()) / (24 * time.Hour))
}

type TasksStore interface {
	Tasks() ([]Task, error)

	Events() ([]Event, error)
}

var config struct {
	Addr string
}

func main() {
	config.Addr = "0.0.0.0:5000"

	srv := &server{
		store: &memoryStore{},
	}

	http.HandleFunc("/", srv.handleHome)

	http.Handle("/js/htmx.min.js", http.StripPrefix("/js", http.FileServer(http.Dir("."))))

	log.Printf("Listening on http://%s", config.Addr)
	log.Fatal(http.ListenAndServe(config.Addr, nil))
}

type server struct {
	store TasksStore
}

func (s *server) handleHome(w http.ResponseWriter, req *http.Request) {
	tasks, err := s.store.Tasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := s.store.Events()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slices.SortFunc(events, func(a, b Event) int {
		return a.DaysLeft() - b.DaysLeft()
	})

	err = homeTmpl.Execute(w, map[string]any{
		"Tasks":  tasks,
		"Events": events,
	})
	if err != nil {
		log.Println("Error:", err)
	}
}

var homeTmpl = template.Must(template.New("").Parse(`<!doctype html>
<html>
<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1" />
		<title>.trackl</title>

		<style>
		body {
			background-color: white;
		}

		.tasks {
			display: flex;
		}

		.box {
			border: 0.1ch solid black;
			width: 2em;
			height: 2em;
			font-size: medium;

			margin: 0.3ch;

			display: flex;
			align-items: center;
			justify-content: center;
		}

		.box.done {
			background-color: rgba(0, 200, 0, 0.7);
		}

		.events {
			display: flex;
			flex-direction: column;
		}

		.event {
			display: flex;
		}

		.event progress {
			flex-grow: 10;
		}
		</style>
</head>

<body>
		<section id="occasionals" class="tasks">
		{{ range $task := .Tasks }}
			<div class="box {{ $task.State }}" title="{{ $task.Description }}">{{ $task.Icon }}</div>
		{{ end }}
		</section>

		<hr />

		<section class="events">
		{{ range $event := .Events }}
		<div class="event">
		{{ $event.Icon }}<progress max="100" value={{ $event.PercentDone }} title="{{ $event.DaysLeft }} days left"></progress>
		</div>
		{{ end }}
		</section>
	

		<script src="/js/htmx.min.js"></script>
</body>
</html>`))
