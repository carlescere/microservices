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
	DefaultPort  = 8000
	bookingsPath = "database/bookings.json"
)

var (
	bookings models.BookingCollection
	port     int
)

func list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, bookings)
}

func user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	booking, ok := bookings[vars["username"]]
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, booking)
}

func init() {
	content, err := ioutil.ReadFile(bookingsPath)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(content, &bookings)
}

func main() {
	flag.IntVar(&port, "port", DefaultPort, "webserver port")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/bookings", list).Methods("GET")
	r.HandleFunc("/bookings/{username}", user).Methods("GET")
	fmt.Printf("Listening in port %d", port)
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprint(":", port), nil)
}
