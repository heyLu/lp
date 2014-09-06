package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

/*

# the plan

qst: run things quickly

qst - detects the current project type and runs it
qst hello.rb - runs `ruby hello.rb`
qst hello.go - compiles & runs hello.go

*/

var mappings = map[string]func(string) string{
	".go": func(name string) string {
		outpath := strings.TrimSuffix(name, path.Ext(name))
		if !path.IsAbs(outpath) {
			outpath = fmt.Sprintf("./%s", outpath)
		}
		return fmt.Sprintf("go build -o %s %s && %s", outpath, name, outpath)
	},
	".rb": func(name string) string {
		return fmt.Sprintf("ruby %s", name)
	},
	".py": func(name string) string {
		return fmt.Sprintf("python %s", name)
	},
}

var delay = flag.Duration("delay", 1*time.Second, "time to wait until restart")
var autoRestart = flag.Bool("autorestart", true, "automatically restart after command exists")
var command = flag.String("command", "", "command to run ({file} will be substituted)")

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

	file := args[0]
	if !isFile(file) {
		fmt.Fprintf(os.Stderr, "Error: %s is not a file.\n", file)
		os.Exit(1)
	}

	var cmd string
	if command != nil && strings.TrimSpace(*command) != "" {
		cmd = strings.Replace(*command, "{file}", file, -1)
	} else {
		ext := path.Ext(file)
		fn, found := mappings[ext]
		if !found {
			log.Fatalf("error: no mapping found for `%s'", file)
		}
		cmd = fn(file)
	}
	log.Printf("command to run: `%s'", cmd)

	runner := MakeRunner(cmd, *autoRestart)
	go runCmd(file, runner)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
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
	cmd         *exec.Cmd
	shellCmd    string
	started     bool
	autoRestart bool
	restarting  bool
}

func MakeRunner(shellCmd string, autoRestart bool) *Runner {
	return &Runner{nil, shellCmd, false, autoRestart, false}
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
			r.cmd.Stderr = os.Stderr
			r.cmd.Stdout = os.Stdout
			r.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			err := r.cmd.Run()
			var result interface{}
			if err != nil {
				result = err
			} else {
				result = r.cmd.ProcessState
			}
			log.Printf("%s finished: %s", r.shellCmd, result)

			time.Sleep(*delay)
			if !r.restarting && !r.autoRestart {
				r.started = false
				break
			}

			r.restarting = false
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
	r.restarting = true
	if r.started {
		return r.Kill()
	} else {
		return r.Start()
	}
}

func (r *Runner) Stop() {
	r.autoRestart = false
	r.Kill()
}

func isFile(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}
