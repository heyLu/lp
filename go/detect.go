package main

import (
	"fmt"
	"os"
	"path"

	"./fileutil"
)

/*
	detect - guess the project type from files present
*/

type Detection struct {
	Language string
	Type     string
}

type detectorFunc func(string) bool

type Detector struct {
	Detection Detection
	Detector  detectorFunc
}

var Detectors = []Detector{
	{Detection{"clojure", "leiningen"}, clojureLeiningen},
	{Detection{"go", "default"}, goDefault},
	{Detection{"java", "maven"}, javaMaven},
	{Detection{"make", "default"}, makeDefault},
	{Detection{"python", "default"}, pythonDefault},
	{Detection{"ruby", "rails"}, rubyRails},
	{Detection{"ruby", "rake"}, rubyRake},
	{Detection{"ruby", "default"}, rubyDefault},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

	file := os.Args[1]

	for _, detector := range Detectors {
		fmt.Println(detector.Detection, detector.Detector(file))
	}
}

func matchingFileOrDir(file string, pattern string) bool {
	if fileutil.IsFile(file) {
		_, f := path.Split(file)
		isMatch, _ := path.Match(pattern, f)
		return isMatch
	} else {
		return fileutil.MatchExists(path.Join(path.Dir(file), pattern))
	}
}

func hasFile(fileOrDir string, file string) bool {
	return fileutil.IsFile(fileutil.Join(fileOrDir, file))
}

func clojureLeiningen(file string) bool {
	return hasFile(file, "project.clj")
}

func goDefault(file string) bool {
	return matchingFileOrDir(file, "*.go")
}

func javaMaven(file string) bool {
	return hasFile(file, "pom.xml")
}

func makeDefault(file string) bool {
	return hasFile(file, "Makefile")
}

func pythonDefault(file string) bool {
	return matchingFileOrDir(file, "*.py")
}

func rubyRails(file string) bool {
	return hasFile(file, "bin/rails")
}

func rubyRake(file string) bool {
	return hasFile(file, "Rakefile")
}

func rubyDefault(file string) bool {
	return matchingFileOrDir(file, "*.rb")
}
