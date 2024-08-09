package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/heyLu/lp/go/things/storage"
)

var settings struct {
	Addr string
}

func main() {
	flag.StringVar(&settings.Addr, "addr", "localhost:5000", "Address to listen on")
	flag.Parse()

	dbStorage, err := storage.NewDBStorage(context.Background(), "file:things.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbStorage.Close()

	things := &Things{
		handlers: []Handler{
			HandleReminders,
			HandleMath,
			HandleHelp,
		},
		storage: dbStorage,
	}

	http.HandleFunc("/", things.HandleIndex)
	http.HandleFunc("/thing", things.HandleThing)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Listening on http://%s", settings.Addr)
	log.Fatal(http.ListenAndServe(settings.Addr, nil))
}

type Things struct {
	handlers []Handler

	storage storage.Storage
}

type Handler func(ctx context.Context, storage storage.Storage, namespace string, w http.ResponseWriter, input string) error

var ErrNotHandled = errors.New("not handled")

func (t *Things) HandleIndex(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, `<!doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<title>things</title>

	<link rel="stylesheet" href="/static/things.css" />
</head>

<body>
	<main>
		<div>
			<input id="tell-me" name="tell-me" type="text" autofocus placeholder="tell me things"
				hx-post="/thing"
				hx-trigger="load, input changed delay:250ms"
				hx-target="#answer"
				hx-indicator="#waiting" />
		    <img id="waiting" class="htmx-indicator" src="/static/three-dots.svg" />
	    </div>

		<section id="answer">
		</section>
	</main>

	<script src="/static/htmx.min.js"></script>
	<script src="/static/things.js"></script>
</body>
</html>`)
}

func (t *Things) HandleThing(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	tellMe := req.Form.Get("tell-me")
	if tellMe == "" {
		fmt.Fprintln(w)
		return
	}

	namespace := "test" // FIXME: get from path/cookie/stuff (like trackl does)

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	for _, handler := range t.handlers {
		err := handler(ctx, t.storage, namespace, w, tellMe)
		if err == ErrNotHandled {
			continue
		}

		if err != nil {
			fmt.Fprintln(w, err)
		}

		break
	}
}

var mathRe = regexp.MustCompile(`([0-9]|eur|usd)`)

func HandleMath(ctx context.Context, _ storage.Storage, _ string, w http.ResponseWriter, input string) error {
	if !mathRe.MatchString(input) {
		return ErrNotHandled
	}

	cmd := exec.CommandContext(ctx, "qalc", "--terse", "--color=0", input)
	// cmd.Stdin = strings.NewReader(input + "\n")

	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(buf.String())
	}

	fmt.Fprintf(w, "<pre>%s</pre>", html.EscapeString(buf.String()))
	return nil
}

type Reminder struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

var reminders = []Reminder{}

func HandleReminders(ctx context.Context, storage storage.Storage, namespace string, w http.ResponseWriter, input string) error {
	if !strings.HasPrefix(input, "remind") {
		return ErrNotHandled
	}

	parts := strings.SplitN(input, " ", 3)
	if len(parts) == 1 {
		rows, err := storage.Query(ctx, namespace, "reminder", 2)
		if err != nil {
			return err
		}
		defer rows.Close()

		fmt.Fprintln(w, "<ul>")
		for rows.Next() {
			var reminder Reminder
			var date string

			_, err := rows.Scan(&reminder.Description, &date)
			if err != nil {
				return err
			}

			reminder.Date, err = time.Parse(time.RFC3339, date)
			if err != nil {
				return err
			}

			fmt.Fprintf(w, "<li><time time=%q>in %s</time> %s</li>\n",
				reminder.Date,
				reminder.Date.Sub(time.Now()).Truncate(time.Minute),
				reminder.Description)
		}
		fmt.Fprintln(w, "</ul>")
		return nil
	}

	var dur time.Duration
	var err error
	if len(parts) > 1 {
		dur, err = time.ParseDuration(parts[1])
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}

	if strings.HasSuffix(input, "!save") {
		if len(parts) != 3 {
			fmt.Fprintln(w, "usage: remind <time> <description>")
			return nil
		}

		if dur == 0 {
			return nil
		}

		date := time.Now().Add(dur)

		reminder := Reminder{
			Date:        date.Truncate(time.Minute).UTC(),
			Description: parts[2][:len(parts[2])-len("!save")],
		}

		rows, err := storage.Query(ctx, namespace, "reminder", 0, reminder.Description)
		if err != nil {
			return err
		}
		defer rows.Close()

		if !rows.Next() {
			_, err = storage.Insert(ctx, namespace, "reminder", reminder.Description, reminder.Date.Format(time.RFC3339))
			if err != nil {
				return err
			}

			fmt.Fprintln(w, "saved!")
		}
	}

	return nil
}

func HandleHelp(ctx context.Context, _ storage.Storage, _ string, w http.ResponseWriter, input string) error {
	if input != "help" {
		fmt.Fprint(w, "don't know that thing, sorry.  ")
	}

	fmt.Fprintln(w, "try math, echo, ...")

	return nil
}
