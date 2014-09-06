package main

import "errors"
import "flag"
import "fmt"
import "log"
import "os"
import "os/exec"
import "os/signal"
import "path"
import "strings"
import "syscall"
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
	".rb": func(name string) string {
		return fmt.Sprintf("ruby %s", name)
	},
	".py": func(name string) string {
		return fmt.Sprintf("python %s", name)
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
	if !found {
		log.Fatalf("error: no mapping found for `%s'", file)
	}

	runner := MakeRunner(fn(file))
	go runCmd(file, runner)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	log.Printf("got signal: %s, exiting...", s)
	runner.Stop()
}

func runCmd(file string, runner *Runner) {
	runner.Start()
	lastMtime := time.Now()
	for {
		info, err := os.Stat(file)
		if os.IsNotExist(err) {
			log.Fatalf("`%s' disappeared, exiting")
		}

		mtime := info.ModTime()
		if mtime.After(lastMtime) {
			log.Printf("`%s' changed, trying to restart", file)
			runner.Restart()
		}

		lastMtime = mtime
		time.Sleep(1 * time.Second)
	}
}

type Runner struct {
	cmd      *exec.Cmd
	shellCmd string
	started  bool
	restart  bool
}

func MakeRunner(shellCmd string) *Runner {
	return &Runner{nil, shellCmd, false, true}
}

func (r *Runner) Start() error {
	if r.started {
		return errors.New("already started, use Restart()")
	}

	r.started = true
	go func() {
		for {
			log.Printf("running %s", r.shellCmd)
			r.cmd = exec.Command("sh", "-c", r.shellCmd)
			r.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			err := r.cmd.Run()
			log.Printf("%s finished: %s", r.shellCmd, err)

			time.Sleep(1 * time.Second)
			if !r.restart {
				break
			}
		}
	}()

	return nil
}

func (r *Runner) Kill() error {
	pgid, err := syscall.Getpgid(r.cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGTERM)
	}
	return err
}

func (r *Runner) Restart() error {
	return r.Kill()
}

func (r *Runner) Stop() {
	r.restart = false
	r.Kill()
}

func isFile(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}
