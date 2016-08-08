package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var options struct {
	listFiles bool
}

func init() {
	flag.BoolVar(&options.listFiles, "l", false, "Always list files")
}

func main() {
	flag.Parse()

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	f, err := os.Open(dir)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if options.listFiles {
		cmd := exec.Command("ls", dir)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	var cmd *exec.Cmd
	if fi.IsDir() {
		fis, err := f.Readdir(-1)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		var rs []string
		for _, fi := range fis {

			if fi.IsDir() {
				continue
			}

			if strings.HasPrefix(strings.ToLower(fi.Name()), "readme") {
				rs = append(rs, fi.Name())
			}
		}

		if len(rs) == 0 {
			cmd = exec.Command("ls", dir)
		} else {
			if options.listFiles {
				fmt.Println()
			}

			fmt.Printf("'%s':\n", rs[0])
			cmd = exec.Command("head", "-n3", path.Join(dir, rs[0]))
		}
	} else {
		cmd = exec.Command("head", "-n3", path.Join(dir, f.Name()))
	}

	if options.listFiles && cmd.Path == "ls" {
		os.Exit(0)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
