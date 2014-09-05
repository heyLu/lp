package main

import "flag"
import "fmt"
import "log"
import "os"
import "os/exec"
import "path"
import "strings"
import "time"

/*

# the plan

qst: run things quickly

qst - detects the current project type and runs it
qst hello.rb - runs `ruby hello.rb`
qst hello.go - compiles & runs hello.go

qst -watch hello.go - watches changes to hello.go and recompiles and runs it on changes

*/

var mappings = map[string]func(string) string{
	".go": func(name string) string {
		return fmt.Sprintf("go build %s && ./%s", name, strings.TrimSuffix(name, path.Ext(name)))
	},
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <file>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	file := os.Args[1]
	if !isFile(file) {
		fmt.Fprintf(os.Stderr, "Error: %s is not a file.\n", file)
		os.Exit(1)
	}

	ext := path.Ext(file)
	fn, found := mappings[ext]
	if found {
		runAndWatch(file, fn(file))
	} else {
		fmt.Fprintf(os.Stderr, "Error: No command defined for %s files", ext)
		os.Exit(1)
	}
}

func runAndWatch(file string, cmd string) {
	// run command, if file changes (mtime) restart command
	// for now: run command until it exits, wait a bit, run again
	runShellCmd(cmd)
	time.Sleep(1 * time.Second)
	runAndWatch(file, cmd)
}

func runShellCmd(cmd string) {
	log.Printf("Running: `%s'", cmd)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	os.Stderr.Write(output)
	log.Printf("Error running command: %s\n", err.Error())
}

func isFile(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}
