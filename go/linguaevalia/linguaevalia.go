package main

// lingua evalia
//
// try it with `curl -i localhost:8000/run --data-binary @hello-world.go`

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Language interface {
	RunFile(f *os.File) (result []byte, err error)
	Extension() string
}

type LanguageGeneral struct {
	Name    string
	Ext     string
	Command string
	Args    []string
}

func (l LanguageGeneral) RunFile(f *os.File) ([]byte, error) {
	args := append(l.Args, f.Name())
	cmd := exec.Command(l.Command, args...)
	return cmd.CombinedOutput()
}

func (l LanguageGeneral) Extension() string { return l.Ext }

var Go = LanguageGeneral{"Go", "go", "go", []string{"run"}}
var Python = LanguageGeneral{"Python", "py", "python", []string{}}
var Ruby = LanguageGeneral{"Ruby", "rb", "ruby", []string{}}
var JavaScript = LanguageGeneral{"JavaScript", "js", "node", []string{}}
var Haskell = LanguageGeneral{"Haskell", "hs", "runhaskell", []string{}}

var languageMappings = map[string]Language{
	"go":         Go,
	"python":     Python,
	"ruby":       Ruby,
	"javascript": JavaScript,
	"haskell":    Haskell,
}

func writeCode(code string, extension string) (*os.File, error) {
	// create tmp file
	f, err := os.Create(fmt.Sprintf("/tmp/linguaevalia-%s.%s", extension, extension)) // FIXME: actually create a tmpfile
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

func Eval(lang Language, code string) ([]byte, error) {
	// write code to temp file
	f, err := writeCode(code, lang.Extension())
	defer f.Close()
	if err != nil {
		return nil, err
	}
	// `go run` it
	res, err := lang.RunFile(f)
	if err != nil {
		return res, err
	}
	// remove the file
	os.Remove(f.Name())
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
	homePageTemplate.Execute(w, nil)
}

var homePageTemplate = template.Must(template.New("homepage").Parse(homePageTemplateStr))

func main() {
	addr, port := "localhost", 8000
	fmt.Printf("running on %s:%d\n", addr, port)

	http.HandleFunc("/run", runCodeHandler)
	http.HandleFunc("/", homePageHandler)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
	if err != nil {
		log.Fatal(err)
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
    }

    .error { color: red; }
    </style>
  </head>

  <body>
    <div id="codeContainer">
      <textarea id="code" rows="20" cols="80">package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}
</textarea>
      <select id="language">
        <option value="go">Go</option>
        <option value="python">Python</option>
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
  </body>
</html>
`
