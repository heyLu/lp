package main

import "flag"
import "fmt"
import "log"
import "os"
import "os/exec"
import "path"
import "strings"
import "sync"
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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

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

func runAndWatch(file string, cmdLine string) {
	// run command, if file changes (mtime) restart command
	lastMtime := time.Now()
	cmd := ShellCmd(cmdLine)
	cmd.Start()

	for {
		info, err := os.Stat(file)
		if err != nil {
			log.Fatalf("Error: %s disappeared, exiting.", file)
		}

		mtime := info.ModTime()
		if mtime.After(lastMtime) {
			log.Printf("%s changed, rerunning", file)
			cmd.Restart()
		}

		lastMtime = mtime
		time.Sleep(1 * time.Second)
	}
}

type RestartableCommand struct {
	Cmd  *exec.Cmd
	Lock sync.Mutex
	Name string
	Args []string
}

func ShellCmd(cmd string) RestartableCommand {
	return RestartableCommand{nil, sync.Mutex{}, "sh", []string{"-c", cmd}}
}

func (c *RestartableCommand) Start() {
	c.Cmd = exec.Command(c.Name, c.Args...)
	c.Lock.Lock()
	go func() {
		c.Cmd.Run()
		c.Lock.Unlock()
	}()
}

func (c *RestartableCommand) Restart() {
	c.Cmd.Process.Kill()
	c.Start()
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
