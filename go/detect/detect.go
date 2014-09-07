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
	&Project{"c/default", Commands{"run": "gcc -o $(basename {file} .c) {file} && ./$(basename {file} .c)"}, cDefault},
	&Project{"clojure/leiningen", Commands{"build": "lein uberjar", "run": "lein run", "test": "lein test"},
		clojureLeiningen},
	&Project{"coffeescript/default", Commands{"run": "coffee {file}"}, coffeescriptDefault},
	&Project{"docker/fig", Commands{"build": "fig build", "run": "fig up"}, dockerFig},
	&Project{"docker/default", Commands{"build": "docker build ."}, dockerDefault},
	&Project{"executable", Commands{"run": "{file}"}, executableDefault},
	&Project{"go/default", Commands{"build": "go build {file}", "run": "go build $(basename {file}) && ./$(basename {file} .go)"},
		goDefault},
	&Project{"haskell/cabal", Commands{"build": "cabal build", "run": "cabal run", "test": "cabal test"}, haskellCabal},
	&Project{"haskell/default", Commands{"run": "runhaskell {file}"}, haskellDefault},
	&Project{"idris/default", Commands{"run": "idris -o $(basename {file} .idr) {file} && ./$(basename {file} .idr)"},
		idrisDefault},
	&Project{"java/maven", Commands{"build": "mvn compile", "test": "mvn compile test"}, javaMaven},
	&Project{"javascript/npm", Commands{"build": "npm install", "test": "npm test"}, javascriptNpm},
	&Project{"javascript/meteor", Commands{"run": "meteor"}, javascriptMeteor},
	&Project{"javascript/default", Commands{"run": "node {file}"}, javascriptDefault},
	&Project{"julia/default", Commands{"run": "julia {file}"}, juliaDefault},
	&Project{"python/django", Commands{"build": "python manage.py syncdb", "run": "python manage.py runserver",
		"test": "python manage.py test"}, pythonDjango},
	&Project{"python/default", Commands{"run": "python {file}"}, pythonDefault},
	&Project{"ruby/rails", Commands{"build": "bundle exec rake db:migrate", "run": "rails server",
		"test": "bundle exec rake test"}, rubyRails},
	&Project{"ruby/rake", Commands{"run": "rake", "test": "rake test"}, rubyRake},
	&Project{"ruby/default", Commands{"run": "ruby {file}"}, rubyDefault},
	&Project{"rust/cargo", Commands{"build": "cargo build", "run": "cargo run", "test": "cargo test"}, rustCargo},
	&Project{"rust/default", Commands{"run": "rustc {file} && ./$(basename {file} .rs)"}, rustDefault},
	&Project{"cmake", Commands{"build": "mkdir .build && cd .build && cmake .. && make"}, cmakeDefault},
	&Project{"make", Commands{"run": "make", "test": "make test"}, makeDefault},
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

func GetById(id string) *Project {
	for _, project := range ProjectTypes {
		if project.Id == id {
			return project
		}
	}
	return nil
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

func cDefault(file string) bool {
	return matchingFileOrDir(file, "*.c")
}

func clojureLeiningen(file string) bool {
	return hasFile(file, "project.clj")
}

func cmakeDefault(file string) bool {
	return hasFile(file, "CMakeLists.txt")
}

func coffeescriptDefault(file string) bool {
	return matchingFileOrDir(file, "*.coffee")
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

func haskellCabal(file string) bool {
	return matchingFileOrDir(file, "*.cabal")
}

func haskellDefault(file string) bool {
	return matchingFileOrDir(file, "*.hs") || matchingFileOrDir(file, "*.lhs")
}

func idrisDefault(file string) bool {
	return matchingFileOrDir(file, "*.idr")
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

func juliaDefault(file string) bool {
	return matchingFileOrDir(file, "*.jl")
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

func rustCargo(file string) bool {
	return hasFile(file, "Cargo.toml")
}

func rustDefault(file string) bool {
	return matchingFileOrDir(file, "*.rs")
}
