package main

// lingua evalia
//
// try it with `curl -i localhost:8000/run --data-binary @hello-world.go`

import (
	"crypto/rand"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Language interface {
	RunFile(f *os.File) (result []byte, err error)
	Name() string
	Extension() string
}

type LanguageGeneral struct {
	name    string
	ext     string
	command string
	args    []string
}

func (l LanguageGeneral) RunFile(f *os.File) ([]byte, error) {
	args := append(l.args, f.Name())
	cmd := exec.Command(l.command, args...)
	return cmd.CombinedOutput()
}

func (l LanguageGeneral) Name() string { return l.name }

func (l LanguageGeneral) Extension() string { return l.ext }

var Go = LanguageGeneral{"Go", "go", "go", []string{"run"}}
var Python = LanguageGeneral{"Python", "py", "python", []string{}}
var Ruby = LanguageGeneral{"Ruby", "rb", "ruby", []string{}}
var JavaScript = LanguageGeneral{"JavaScript", "js", "node", []string{}}
var Haskell = LanguageGeneral{"Haskell", "hs", "runhaskell", []string{}}
var Rust = LanguageGeneral{"Rust", "rs", "./bin/run-rust", []string{}}
var Julia = LanguageGeneral{"Julia", "jl", "julia", []string{}}
var Pixie = LanguageGeneral{"Pixie", "pxi", "pixie-vm", []string{}}
var C = LanguageGeneral{"C", "c", "./bin/run-c", []string{}}
var Bash = LanguageGeneral{"Bash", "bash", "bash", []string{}}
var Lua = LanguageGeneral{"Lua", "lua", "lua", []string{}}
var CPlusPlus = LanguageGeneral{"C++", "cpp", "./bin/run-c++", []string{}}

var languageMappings = map[string]Language{
	"go":         Go,
	"python":     Python,
	"ruby":       Ruby,
	"javascript": JavaScript,
	"haskell":    Haskell,
	"rust":       Rust,
	"julia":      Julia,
	"pixie":      Pixie,
	"c":          C,
	"bash":       Bash,
	"lua":        Lua,
	"cpp":        CPlusPlus,
}

func writeCode(code string, extension string) (*os.File, error) {
	// create tmp file
	f, err := tempFile("/tmp", "linguaevalia", extension)
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

func tempFile(dir, prefix, suffix string) (*os.File, error) {
	rnd, _ := rand.Int(rand.Reader, big.NewInt(999999))
	f, err := os.Create(path.Join(dir, fmt.Sprintf("%s%d.%s", prefix, rnd, suffix)))
	return f, err
}

func Eval(lang Language, code string) ([]byte, error) {
	// write code to temp file
	f, err := writeCode(code, lang.Extension())
	defer f.Close()
	defer os.Remove(f.Name())
	if err != nil {
		return nil, err
	}
	// `go run` it
	res, err := lang.RunFile(f)
	if err != nil {
		return res, err
	}
	// return output
	return res, nil
}

func runCodeHandler(w http.ResponseWriter, r *http.Request) {
	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lang := getLanguage(r)
	res, err := Eval(lang, string(code))
	if err != nil {
		http.Error(w, string(res), http.StatusNotAcceptable)
		return
	}
	w.Write(res)
}

func getLanguage(r *http.Request) Language {
	langName := r.URL.Query().Get("language")
	if langName != "" {
		lang, ok := languageMappings[langName]
		if ok {
			return lang
		}
	}
	return Go
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	bindings := map[string]interface{}{
		"languages": languageMappings,
	}
	homePageTemplate.Execute(w, bindings)
}

var homePageTemplate = template.Must(template.New("homepage").Parse(homePageTemplateStr))

func runServer() {
	address := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("running on %s\n", address)

	http.HandleFunc("/run", runCodeHandler)
	http.HandleFunc("/codemirror.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "lib/codemirror.js")
	})
	http.HandleFunc("/codemirror.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "lib/codemirror.css")
	})
	http.HandleFunc("/", homePageHandler)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func languageForExtension(extension string) *Language {
	var language *Language = nil
	for _, lang := range languageMappings {
		if "."+lang.Extension() == extension {
			return &lang
		}
	}
	return language
}

func runOnce(args []string) {
	var (
		f        *os.File
		err      error
		langName string = *language
	)

	if len(args) > 0 {
		if *language == "" {
			l := languageForExtension(path.Ext(args[0]))
			if l == nil {
				fmt.Printf("Error: Don't know how to handle '%s' files\n", path.Ext(args[0]))
				os.Exit(1)
			}
			langName = (*l).Name()
		}
		f, err = os.Open(args[0])
	} else {
		f, err = os.Stdin, nil
	}
	defer f.Close()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	langName = strings.ToLower(langName)
	lang, ok := languageMappings[langName]
	if !ok {
		fmt.Printf("Error: Unknown language '%s'\n", langName)
		os.Exit(1)
	}

	if f == os.Stdin {
		f, err = tempFile("/tmp", "linguaevalia", lang.Extension())
		_, err = io.Copy(f, os.Stdin)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		defer os.Remove(f.Name())
	}

	res, err := lang.RunFile(f)
	os.Stdout.Write(res)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func runHelp() {
	fmt.Printf(`Usage: %s [cmd] [options]

Availlable commands:

server		- Starts a web server. (Default.)
run		- Runs code from a file or from stdin.
		  (If running from stdin, you must pass the language
		   using the -l flag.)
help		- Display this help message.

`,
		os.Args[0])
}

func parseCommand() (string, []string) {
	if len(os.Args) == 1 {
		return "server", []string{}
	} else {
		return os.Args[1], os.Args[2:]
	}
}

var language = flag.String("l", "", "The language to use for code passed via stdin.")
var host = flag.String("h", "localhost", "The host to listen on.")
var port = flag.Int("p", 8000, "The port to listen on.")

func main() {
	cmd, args := parseCommand()
	flag.CommandLine.Parse(args)

	switch cmd {
	case "server":
		runServer()
	case "run":
		runOnce(flag.Args())
	case "help":
		runHelp()
	default:
		fmt.Println("Error: Unknown command:", cmd)
		os.Exit(1)
	}
}

const homePageTemplateStr = `
<!doctype html>
<html>
  <head>
    <title>lingua evalia</title>
    <meta charset="utf-8" />
    <style type="text/css">
    #codeContainer {
      position: relative;
      display: inline-block;
    }

    #code {
      border: none;
    }

    #language {
      position: absolute;
      right: 0;
      top: 0;
      z-index: 10; /* above codemirror */
    }

    .error { color: red; }
    </style>
    <link rel="stylesheet" type="text/css" href="/codemirror.css" />
    <style type="text/css">
    .CodeMirror {
      min-width: 80ex;
    }
    </style>
  </head>

  <body>
    <div id="codeContainer">
      <textarea id="code" autofocus rows="20" cols="80">package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}
</textarea>
      <select id="language">
        {{ range $short, $name := .languages }}
        <option value="{{ $short }}">{{ $name.Name }}</option>
        {{ end }}
      </select>
    </div>
    <pre id="result"></pre>

    <script>
      var codeEl = document.getElementById("code");
      var languageEl = document.getElementById("language");
      var resultEl = document.getElementById("result");

      codeEl.onkeydown = function(ev) {
        if (ev.ctrlKey && ev.keyCode == 13) {
          resultEl.textContent = "";
          sendCode(codeEl.value, languageEl.value, function(xhr) {
            resultEl.className = xhr.status == 200 ? "success" : "error";
            resultEl.textContent = xhr.response;
          });
        }
      }


      function sendCode(code, language, cb) {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/run?language=" + language);
        xhr.onreadystatechange = function(ev) {
          if (xhr.readyState == XMLHttpRequest.DONE) {
            cb(xhr);
          }
        };
        xhr.send(code);
      }
    </script>

    <script src="/codemirror.js"></script>
    <script>
      var cm = CodeMirror.fromTextArea(codeEl, {mode: languageToMode(languageEl.value)});

      cm.on("changes", function(cm) { codeEl.value = cm.getValue(); });

      cm.setOption("extraKeys", {
        "Ctrl-Enter": function(cm) {
          resultEl.textContent = "";
          sendCode(cm.getValue(), languageEl.value, function(xhr) {
            resultEl.className = xhr.status == 200 ? "success" : "error";
            resultEl.textContent = xhr.response;
          });
        }
      });

      languageEl.onchange = function(ev) {
        cm.setOption("mode", languageToMode(languageEl.value));
      };

      function languageToMode(language) {
        switch(language) {
        case "bash": return "shell";
        case "pixie": return "clojure";
        case "c": return "text/x-csrc";
        case "cpp": return "text/x-c++src";
        default: return language;
        }
      }
    </script>
  </body>
</html>
`
