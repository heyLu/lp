package detect

import (
	"errors"
	"path"

	"../fileutil"
)

/*
	detect - guess the project type from files present
*/

type Project struct {
	Id       string
	Commands Commands
	Detect   func(string) bool
}

type Commands map[string]string

var ProjectTypes = []*Project{
	&Project{"clojure/leiningen", Commands{"run": "lein run", "test": "lein test"}, clojureLeiningen},
	&Project{"docker/fig", Commands{"run": "fig up"}, dockerFig},
	&Project{"docker/default", Commands{}, dockerDefault},
	&Project{"executable", Commands{"run": "{file}"}, executableDefault},
	&Project{"go/default", Commands{"build": "go build {file}", "run": "go build {file} && $(basename {file} .go)"},
		goDefault},
	&Project{"java/maven", Commands{"build": "mvn compile", "test": "mvn compile test"}, javaMaven},
	&Project{"javascript/npm", Commands{}, javascriptNpm},
	&Project{"javascript/meteor", Commands{"run": "meteor"}, javascriptMeteor},
	&Project{"javascript/default", Commands{"run": "node {file}"}, javascriptDefault},
	&Project{"python/django", Commands{}, pythonDjango},
	&Project{"python/default", Commands{"run": "python {file}"}, pythonDefault},
	&Project{"ruby/rails", Commands{"run": "rails server", "test": "bundle exec rake test"}, rubyRails},
	&Project{"ruby/rake", Commands{"run": "rake"}, rubyRake},
	&Project{"ruby/default", Commands{"run": "ruby {file}"}, rubyDefault},
	&Project{"make", Commands{"run": "make"}, makeDefault},
	&Project{"procfile", Commands{}, procfileDefault},
}

func Detect(file string) (*Project, error) {
	for _, project := range ProjectTypes {
		if project.Detect(file) {
			return project, nil
		}
	}

	return nil, errors.New("no project matches")
}

func DetectAll(file string) []*Project {
	projects := make([]*Project, 0, len(ProjectTypes))

	for _, project := range ProjectTypes {
		if project.Detect(file) {
			n := len(projects)
			projects = projects[0 : n+1]
			projects[n] = project
		}
	}

	return projects
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

func dockerFig(file string) bool {
	return hasFile(file, "fig.yml")
}

func dockerDefault(file string) bool {
	return hasFile(file, "Dockerfile")
}

func executableDefault(file string) bool {
	return fileutil.IsExecutable(file)
}

func goDefault(file string) bool {
	return matchingFileOrDir(file, "*.go")
}

func javaMaven(file string) bool {
	return hasFile(file, "pom.xml")
}

func javascriptNpm(file string) bool {
	return hasFile(file, "package.json")
}

func javascriptMeteor(file string) bool {
	return hasFile(file, ".meteor/.id")
}

func javascriptDefault(file string) bool {
	return matchingFileOrDir(file, "*.js")
}

func makeDefault(file string) bool {
	return hasFile(file, "Makefile")
}

func procfileDefault(file string) bool {
	return hasFile(file, "Procfile")
}

func pythonDjango(file string) bool {
	return hasFile(file, "manage.py")
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
