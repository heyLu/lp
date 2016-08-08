package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
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
			fmt.Printf("'%s':\n", rs[0])
			cmd = exec.Command("head", "-n3", path.Join(dir, rs[0]))
		}
	} else {
		cmd = exec.Command("head", "-n3", path.Join(dir, f.Name()))
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
