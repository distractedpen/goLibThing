package main

import (
	"fmt"
	"net/http"
)


func handleRoot(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	fmt.Fprintf(w, "Hello, %s", queryParams.Get("name"))
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.ListenAndServe(":8080", nil)
}
