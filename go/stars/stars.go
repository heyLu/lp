// `stars` fetches your GitHub stars and updates them as necessary.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

var config struct {
	userName    string
	directory   string
	concurrency int
}

func init() {
	flag.StringVar(&config.userName, "user", "heyLu", "The GitHub user to fetch stars for")
	flag.StringVar(&config.directory, "directory", "github-stars", "The directory to store the repos in")
	flag.IntVar(&config.concurrency, "concurrency", 10, "The number of repos to update concurrently")
}

func main() {
	flag.Parse()

	var stars []starInfo
	decoder := json.NewDecoder(os.Stdin)
	err := decoder.Decode(&stars)
	if err != nil {
		panic(err)
	}

	sem := make(chan bool, config.concurrency)
	var wg sync.WaitGroup

	wg.Add(len(stars))
	for _, info := range stars {
		info := info

		sem <- true
		go func() {
			fmt.Printf("% 48s - %s\n", info.RepoName, info.Description)
			err := updateRepo(info)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			wg.Done()
			<-sem
		}()
	}

	wg.Wait()
}

func updateRepo(info starInfo) error {
	f, err := os.Open(path.Join(config.directory, info.RepoName))
	if err != nil {
		if os.IsNotExist(err) {
			return gitClone(info)
		}

		return err
	}
	f.Close()

	lastCommit, err := gitLastCommit(info)
	if err != nil {
		return err
	}

	if lastCommit.Before(info.PushedAt) {
		return gitPull(info)
	}

	return nil
}

func gitClone(info starInfo) error {
	cmd := exec.Command("git", "clone", info.CloneUrl,
		path.Join(config.directory, info.RepoName))
	return cmd.Run()
}

func gitPull(info starInfo) error {
	cmd := exec.Command("git", "-C", path.Join(config.directory, info.RepoName), "pull")
	return cmd.Run()
}

func gitLastCommit(info starInfo) (time.Time, error) {
	cmd := exec.Command("git", "-C", path.Join(config.directory, info.RepoName),
		"log", "-n", "1", "--format=%cd", "--date=iso8601-strict")
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	outStr := strings.TrimSpace(string(out))
	return time.Parse(time.RFC3339, outStr)
}

type starInfo struct {
	RepoName    string    `json:"full_name"`
	Description string    `json:"description"`
	CloneUrl    string    `json:"git_url"`
	PushedAt    time.Time `json:"pushed_at"`
}
