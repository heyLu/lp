package main

import "fmt"
import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Path[1:]
		if name == "" {
			name = "World"
		}
		w.Write([]byte(fmt.Sprintf("Hello, %s!", name)))
	})

	fmt.Println("Running on :8080")
	http.ListenAndServe(":8080", nil)
}
