package main

import (
	"encoding/json"
	"errors"
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
	BookingsURL = "http://172.17.0.1:8003/bookings/%s"
	MoviesURL   = "http://172.17.0.1:8001/movies/%s"
	usersPath   = "database/users.json"
)

var (
	users      models.UserCollection
	httpClient http.Client
	port       int
)

func list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, users)
}

func user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, ok := users[vars["username"]]

	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, user)
}

func userBookings(w http.ResponseWriter, r *http.Request) {
	booking, err := getBooking(w, r)
	if err != nil {
		return
	}

	movies, err := getMovies(w, r, booking)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, movies)
}

func getMovies(w http.ResponseWriter, r *http.Request, booking models.Booking) (models.UserBookings, error) {
	movies := make(map[string][]models.Movie)
	for date, movieLists := range booking {
		movies[date] = make([]models.Movie, 0)
		for _, movieID := range movieLists {
			movie, err := getMovie(w, r, movieID)

			if err != nil {
				return nil, err
			}

			if movie != nil {
				movies[date] = append(movies[date], *movie)
			}
		}
	}

	return models.UserBookings(movies), nil
}

func getMovie(w http.ResponseWriter, r *http.Request, movieID string) (*models.Movie, error) {
	resp, err := httpClient.Get(fmt.Sprintf(MoviesURL, movieID))
	if err != nil {
		txt := http.StatusText(http.StatusServiceUnavailable)
		http.Error(w, txt, http.StatusServiceUnavailable)
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		http.NotFound(w, r)
		return nil, nil
	}

	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		txt := http.StatusText(http.StatusServiceUnavailable)
		http.Error(w, txt, http.StatusServiceUnavailable)
		return nil, err
	}

	var movie models.Movie
	err = json.Unmarshal(content, &movie)
	if err != nil {
		txt := http.StatusText(http.StatusInternalServerError)
		http.Error(w, txt, http.StatusInternalServerError)
		return nil, err
	}

	return &movie, nil
}

func getBooking(w http.ResponseWriter, r *http.Request) (models.Booking, error) {
	vars := mux.Vars(r)

	user, ok := users[vars["username"]]

	if !ok {
		http.NotFound(w, r)
		return nil, errors.New("User not found")
	}

	resp, err := httpClient.Get(fmt.Sprintf(BookingsURL, user.Username))
	if err != nil {
		txt := http.StatusText(http.StatusServiceUnavailable)
		http.Error(w, txt, http.StatusServiceUnavailable)
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		http.NotFound(w, r)
		return nil, errors.New("No bookings")
	}

	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		txt := http.StatusText(http.StatusServiceUnavailable)
		http.Error(w, txt, http.StatusServiceUnavailable)
		return nil, err
	}

	var bookings models.Booking
	err = json.Unmarshal(content, &bookings)
	if err != nil {
		txt := http.StatusText(http.StatusInternalServerError)
		http.Error(w, txt, http.StatusInternalServerError)
		return nil, err
	}

	return bookings, nil
}

func init() {
	content, err := ioutil.ReadFile(usersPath)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(content, &users)
	httpClient = http.Client{}
}

func main() {
	flag.IntVar(&port, "port", DefaultPort, "webserver port")

	r := mux.NewRouter()
	r.HandleFunc("/users", list).Methods("GET")
	r.HandleFunc("/users/{username}", user).Methods("GET")
	r.HandleFunc("/users/{username}/bookings", userBookings).Methods("GET")
	fmt.Printf("Listening in port %d", port)
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprint(":", port), nil)
}
