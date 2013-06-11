package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Oi pessoal! :)")
	})
	http.ListenAndServe("127.0.0.1:7070", nil)
}
