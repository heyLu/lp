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

func runCodeHandler(w http.ResponseWriter, r *http.Request) {
	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println(string(code))
	res, err := Go.Eval(string(code))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	w.Write(res)
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
          sendCode(codeEl.value, function(result) {
            resultEl.textContent = result;
          });
        }
      }

      function sendCode(code, cb) {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/run");
        xhr.onreadystatechange = function(ev) {
          if (xhr.readyState == XMLHttpRequest.DONE) {
            cb(xhr.response);
          }
        };
        xhr.send(code);
      }
    </script>
  </body>
</html>
`
