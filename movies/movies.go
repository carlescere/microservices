package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/carlescere/microservices/models"
	"github.com/gorilla/mux"
)

const (
	DefaultPort = 8000
	moviesPath  = "database/movies.json"
)

var (
	movies models.Catalog
	port   int
)

func list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, movies)
}

func movie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	movie, ok := movies[vars["id"]]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, movie)
}

func init() {
	content, err := ioutil.ReadFile(moviesPath)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(content, &movies)
}

func main() {
	flag.IntVar(&port, "port", DefaultPort, "webserver port")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/movies", list).Methods("GET")
	r.HandleFunc("/movies/{id}", movie).Methods("GET")
	fmt.Printf("Listening in port %d", port)
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprint(":", port), nil)
}
