// `stars` fetches your GitHub stars and updates them as necessary.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

var userName = "heyLu"
var directory = "github-stars"

func main() {
	var stars []starInfo
	decoder := json.NewDecoder(os.Stdin)
	err := decoder.Decode(&stars)
	if err != nil {
		panic(err)
	}

	for _, info := range stars {
		fmt.Printf("% 48s - %s\n", info.RepoName, info.Description)
		err := updateRepo(info)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func updateRepo(info starInfo) error {
	f, err := os.Open(path.Join(directory, info.RepoName))
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
		path.Join(directory, info.RepoName))
	return cmd.Run()
}

func gitPull(info starInfo) error {
	cmd := exec.Command("git", "-C", path.Join(directory, info.RepoName), "pull")
	return cmd.Run()
}

func gitLastCommit(info starInfo) (time.Time, error) {
	cmd := exec.Command("git", "-C", path.Join(directory, info.RepoName),
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
