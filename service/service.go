package service

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3" //sqlite3 driver
)

//constants
const (
	PROJECTIDLEN = 8 //length of the project id in the token
	USERIDLEN    = 8 //length of the user id in the token
	MASKLEN      = 4 //total legnth of config mask
	ISVALID      = 1 //bit to indicate if token is valid
	ISADMIN      = 2 //bit to indicate if token belongs to admin
)

var (
	dbi       db
	appConfig config
)

type db struct {
	database *sql.DB
}

type config struct {
	LazyAuth bool
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)

	v2 := r.PathPrefix("/v2").Subrouter()
	v2.HandleFunc("/token", v2getTokenHandler).Methods("POST")
	v2.HandleFunc("/token/{tokenID}", v2validateTokenHandler).Methods("POST")

	v3 := r.PathPrefix("/v3").Subrouter()
	v3.HandleFunc("/tokens", v3ValidateTokenHandler).Methods("GET")

	http.Handle("/", r)
}

//Serve - run service
func Serve() {
	http.ListenAndServe(":8080", nil)
}

//V2 Service Handlers:
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "cool")
}
