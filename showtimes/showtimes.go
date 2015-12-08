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
	DefaultPort   = 8000
	showtimesPath = "database/showtimes.json"
)

var (
	showtimesByDate  models.Showtimes
	showtimesByMovie models.Showtimes
	port             int
)

func list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, showtimesByDate)
}

func date(w http.ResponseWriter, r *http.Request) {
	getItem(w, r, showtimesByDate)
}

func movie(w http.ResponseWriter, r *http.Request) {
	getItem(w, r, showtimesByMovie)
}

func getItem(w http.ResponseWriter, r *http.Request, database models.Showtimes) {
	vars := mux.Vars(r)

	date, ok := database[vars["key"]]

	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(date)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(b))
}

func init() {
	content, err := ioutil.ReadFile(showtimesPath)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(content, &showtimesByDate)

	showtimesByMovie = make(map[string][]string)

	for date, movies := range showtimesByDate {
		for _, movie := range movies {
			if _, ok := showtimesByMovie[movie]; !ok {
				showtimesByMovie[movie] = make([]string, 0)
			}
			showtimesByMovie[movie] = append(showtimesByMovie[movie], date)
		}
	}
}

func main() {
	flag.IntVar(&port, "port", DefaultPort, "webserver port")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/showtimes", list).Methods("GET")
	r.HandleFunc("/showtimes/{key:[0-9]{8}}", date).Methods("GET")
	r.HandleFunc("/showtimes/{key:[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}}", movie).Methods("GET")
	fmt.Printf("Listening in port %d", port)
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprint(":", port), nil)
}
