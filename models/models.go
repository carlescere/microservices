package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type BookingCollection map[string]Booking

func (b BookingCollection) String() string {
	c, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	return string(c)
}

type Booking map[string][]string

func (b Booking) String() string {
	c, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	return string(c)
}

type jsonMovie Movie

type Movie struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Rating   float32 `json:"rating"`
	Director string  `json:"director"`
}

func (m Movie) URI() string {
	return fmt.Sprintf("/movies/%s", m.ID)
}

func (m Movie) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		jsonMovie
		URI string `json:"uri"`
	}{
		jsonMovie: jsonMovie(m),
		URI:       m.URI(),
	})
}

func (m Movie) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

type Catalog map[string]Movie

func (c Catalog) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

type Showtimes map[string][]string

func (s Showtimes) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	timestamp := time.Unix(int64(ts), 0)

	*t = Timestamp(timestamp)
	return nil
}

func (t *Timestamp) String() string {
	return time.Time(*t).String()
}

type UserCollection map[string]User

func (u UserCollection) String() string {
	c, err := json.Marshal(&u)
	if err != nil {
		log.Fatal(err)
	}
	return string(c)
}

type User struct {
	Username   string     `json:"id"`
	Name       string     `json:"name"`
	LastActive *Timestamp `json:"last_active"`
}

func (u User) String() string {
	c, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	return string(c)
}

type UserBookings map[string][]Movie

func (u UserBookings) String() string {
	c, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	return string(c)
}
