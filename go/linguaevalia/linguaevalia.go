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

type LanguageGo struct{}
type LanguagePython struct{}

var Go = LanguageGo{}
var Python = LanguagePython{}

var languageMappings = map[string]Language{
	"go":     Go,
	"python": Python,
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

func (l LanguageGo) RunFile(f *os.File) ([]byte, error) {
	cmd := exec.Command("go", "run", f.Name())
	return cmd.CombinedOutput()
}

func (l LanguageGo) Extension() string { return "go" }

func (l LanguagePython) RunFile(f *os.File) ([]byte, error) {
	cmd := exec.Command("python", f.Name())
	return cmd.CombinedOutput()
}

func (l LanguagePython) Extension() string { return "py" }

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
    #code {
      border: none;
    }

    .error { color: red; }
    </style>
  </head>

  <body>
    <textarea id="code" rows="20" cols="80">package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}
</textarea>
    <pre id="result"></pre>

    <script>
      var codeEl = document.getElementById("code");
      var resultEl = document.getElementById("result");

      codeEl.onkeydown = function(ev) {
        if (ev.ctrlKey && ev.keyCode == 13) {
          resultEl.textContent = "";
          sendCode(codeEl.value, function(xhr) {
            resultEl.className = xhr.status == 200 ? "success" : "error";
            resultEl.textContent = xhr.response;
          });
        }
      }

      function sendCode(code, cb) {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/run");
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
