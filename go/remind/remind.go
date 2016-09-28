package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

type Reminder struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

type byDate []Reminder

func (r byDate) Len() int           { return len(r) }
func (r byDate) Less(i, j int) bool { return r[i].Date.Before(r[j].Date) }
func (r byDate) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

var flags struct {
	showAll bool
}

func init() {
	flag.BoolVar(&flags.showAll, "all", false, "Show all reminders")
}

func main() {
	flag.Parse()

	cmd := "list"
	if flag.NArg() >= 1 {
		cmd = flag.Arg(0)
	}

	needWrite := false
	f, err := os.OpenFile("remind.json", os.O_RDWR, 0644)
	if err != nil && !os.IsNotExist(err) {
		exit(err)
	}
	defer f.Close()

	var reminders []Reminder
	if err == nil {
		dec := json.NewDecoder(f)
		err = dec.Decode(&reminders)
		if err != nil {
			exit(err)
		}
	}

	switch cmd {
	case "list", "l":
		min := time.Now()
		max := min.AddDate(0, 0, 7)

		if flag.NArg() > 1 {
			switch flag.Arg(1) {
			case "today":
				min = truncateHours(time.Now())
				max = min.Add(24 * time.Hour)
			}
		}

		sort.Sort(byDate(reminders))
		for _, r := range reminders {
			if flags.showAll || (r.Date.After(min) && r.Date.Before(max)) {
				fmt.Printf("%s - %s\n", r.Date, r.Description)
			}
		}
	case "add":
		if flag.NArg() != 3 {
			flag.Usage()
			os.Exit(1)
		}

		date, err := parseTime(flag.Arg(1))
		if err != nil {
			exit(err)
		}
		description := flag.Arg(2)
		reminders = append(reminders, Reminder{Date: date, Description: description})
		needWrite = true
	default:
		fmt.Fprintf(os.Stderr, "unknown command '%s'\n", cmd)
		os.Exit(1)
	}

	if needWrite {
		if f == nil {
			f, err = os.Create("remind.json")
			if err != nil {
				exit(err)
			}
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			exit(err)
		}

		enc := json.NewEncoder(f)
		err = enc.Encode(reminders)
		if err != nil {
			exit(err)
		}
	}
}

func exit(err error) {
	panic(err)
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}
