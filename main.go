package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var IndexHandler = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{}")
}
