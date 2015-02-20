package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/asm-products/landline-api/models"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Unauthentic routes
	r := mux.NewRouter().Host("{subdomain:[a-z-]+}.landline.io").Subrouter()
	r.HandleFunc("/", IndexHandler)

	// authenticated routes
	// user := mux.NewRouter()
	r.Handle("/user", authenticated(UserHandler))
	r.Handle("/sessions", authenticated(SessionsHandler)).Methods("POST")
	// user.HandleFunc("/user", SessionsHandler)
	// user.HandleFunc("/sessions", CreateSessionHandler).Methods("POST")
	//
	// r.PathPrefix("/user").Handler(negroni.New(
	// 	negroni.NewRecovery(),
	// 	negroni.NewLogger(),
	// 	negroni.HandlerFunc(TokenMiddleware),
	// 	negroni.Wrap(user),
	// ))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)

	n.UseHandler(r)

	n.Run(":" + port)
}

var IndexHandler = func(w http.ResponseWriter, r *http.Request) {
	subdomain := mux.Vars(r)["subdomain"]

	var endpoints models.OAuthEndpoints
	err := Db.SelectOne(&endpoints, "select oauth_authorize_url, oauth_token_url from teams where slug = $1", subdomain)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(endpoints)
}

var SessionsHandler = func(w http.ResponseWriter, r *http.Request) {
	subdomain := mux.Vars(r)["subdomain"]
	fmt.Fprintf(w, "session "+subdomain)
}

var UserHandler = func(w http.ResponseWriter, r *http.Request) {
	subdomain := mux.Vars(r)["subdomain"]
	fmt.Fprintf(w, "user "+subdomain)
}

var CreateSessionHandler = func(w http.ResponseWriter, r *http.Request) {
	subdomain := mux.Vars(r)["subdomain"]
	fmt.Fprintf(w, "create session "+subdomain)
}

func authenticated(h http.HandlerFunc) http.Handler {
	return negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(TokenMiddleware),
		negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			h(rw, r)
		}),
	)
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

func TokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
	next(w, r)
}
