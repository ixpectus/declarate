package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

var pollCounter atomic.Int32

// The `json:"whatever"` bit is a way to tell the JSON
// encoder and decoder to use those names instead of the
// capitalised names
type person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Items []int  `json:"items"`
}

var tom *person = &person{
	Name: "Tom",
	Age:  28,
	Items: []int{
		1, 2, 3, 4,
	},
}

var pollPerson *person = &person{
	Name: "Tommy",
	Age:  31,
	Items: []int{
		1, 2, 4,
	},
}

func tomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Just send out the JSON version of 'tom'
		j, _ := json.Marshal(tom)
		w.Write(j)
	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		p := &person{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tom = p
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func pollHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Just send out the JSON version of 'tom'
		v := pollCounter.Add(1)
		if v == 5 {
			j, _ := json.Marshal(pollPerson)
			w.Write(j)
			pollCounter.Store(0)
		} else {
			j, _ := json.Marshal(tom)
			w.Write(j)
		}
	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		p := &person{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tom = p
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func Handle() {
	http.HandleFunc("/tom", tomHandler)
	http.HandleFunc("/poll", pollHandler)
	log.Println("Go!")
	http.ListenAndServe("127.0.0.1:8181", nil)
}
