package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s <cmd>\n", os.Args[0])
		fmt.Fprintln(os.Stderr)

		fmt.Fprintf(os.Stderr, "Available commands:\n")
		fmt.Fprintf(os.Stderr, "  add  <date> <description>\n")
		fmt.Fprintf(os.Stderr, "  list <when>\n")
		fmt.Fprintf(os.Stderr, "    where `when` is empty or one of: today\n")
		fmt.Fprintf(os.Stderr, "  l    (alias for `list`)\n")

		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
	}
}

func isCommand(s string) bool {
	return s == "list" || s == "l" || s == "add"
}

func main() {
	flag.Parse()

	start := 0
	cmd := "list"
	if flag.NArg() >= 1 && isCommand(flag.Arg(0)) {
		start = 1
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
	case "list", "l":
		fallthrough
	default:
		min := time.Now()
		max := min.AddDate(0, 0, 7)

		if flag.NArg() > start {
			switch flag.Arg(start) {
			case "today":
				min = truncateHours(time.Now())
				max = min.Add(24 * time.Hour)
			case "this":
				if flag.NArg() > start+1 && flag.Arg(start+1) == "week" {
					min = truncateHours(time.Now())
					min = min.AddDate(0, 0, -int(min.Weekday()))
					max = min.AddDate(0, 0, 7)
				} else {
					fmt.Fprintf(os.Stderr, "invalid specifier: '%s'\n", strings.Join(flag.Args()[start:], " "))
					os.Exit(1)
				}
			default:
				fmt.Fprintf(os.Stderr, "unknown command '%s'\n", cmd)
				os.Exit(1)
			}
		}

		sort.Sort(byDate(reminders))
		for _, r := range reminders {
			if flags.showAll || (r.Date.After(min) && r.Date.Before(max)) {
				fmt.Printf("%s - %s\n", r.Date, r.Description)
			}
		}
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
