package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Language interface {
	Eval(code string) (result []byte, err error)
}

type LanguageGo struct{}

var Go = LanguageGo{}

func (g LanguageGo) Eval(code string) ([]byte, error) {
	// write code to temp file
	f, err := writeCode(code)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	// `go run` it
	res, err := runCode(f)
	if err != nil {
		return nil, err
	}
	// remove the file
	os.Remove(f.Name())
	// return output
	return res, nil
}

func writeCode(code string) (*os.File, error) {
	// create tmp file
	f, err := os.Create("/tmp/linguaevalia-go.go") // FIXME: actually create a tmpfile
	if err != nil {
		return f, err
	}
	// write code to it
	_, err = f.Write([]byte(code))
	if err != nil {
		return f, err
	}
	return f, nil
}

func runCode(f *os.File) ([]byte, error) {
	cmd := exec.Command("go", "run", f.Name())
	return cmd.CombinedOutput()
}

func main() {
	res, err := Go.Eval("package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n")
	if err != nil {
		fmt.Println("Error evaluating: ", err)
	}
	os.Stdout.Write(res)
}
