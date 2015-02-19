package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler)

	http.Handle("/",
		logHandler(
			jsonHandler(
				tokenHandler(r))))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var IndexHandler = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{}")
}

func logHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	}
}

func jsonHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	}
}

func tokenHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "token" {
			http.Error(w, "Authorization header format must be token {token}", http.StatusBadRequest)
			return
		}

		token := authHeaderParts[1]

		context.Set(r, "token", token)
	}
}
