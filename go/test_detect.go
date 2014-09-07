package main

import (
	"fmt"
	"os"

	"./detect"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

	file := os.Args[1]

	for _, project := range detect.ProjectTypes {
		runCmd := project.Commands["run"]
		fmt.Printf("%v (%v): %v\n", project.Id, runCmd, project.Detect(file))
	}
}
