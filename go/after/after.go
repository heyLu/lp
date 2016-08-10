package main

// after - Run a command when the input is finished.
//
// Usage: <input> | after <cmd>

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

const (
	Cr  = 13
	Esc = 27
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: <input> | %s <cmd>", os.Args[0])
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	r := io.TeeReader(os.Stdin, buf)

	// Copy output from command to stdout
	io.Copy(os.Stdout, r)

	// Clear output (after stdin was closed)
	numLines := 0
	for _, ch := range buf.Bytes() {
		if ch == '\n' {
			numLines += 1
		}
	}

	fmt.Fprintf(os.Stdout, "\r")
	for i := 0; i < numLines; i++ {
		fmt.Fprintf(os.Stdout, "%c[K", Esc)
		fmt.Fprintf(os.Stdout, "%c[1A", Esc)
	}
	time.Sleep(2 * time.Second)

	// Run command
	cmd := exec.Command(os.Args[1])
	cmd.Args = os.Args[2:]
	cmd.Stdin = buf
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
