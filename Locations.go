package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Location struct {
	Name            string
	NumberOfPersons int
	Cal             int
}

type LocationResponse struct {
	Id       string          `json:"id"`
	Year     int             `json:"year"`
	Week     int             `json:"week"`
	Month    string          `json:"month"`
	Day      string          `json:"day"`
	Distance int             `json:"distance"`
	Offset   string          `json:"offset"`
	Days     []string        `json:"days"`
	Valid    [][]interface{} `json:"valid"`
}

func locationResponse(w http.ResponseWriter, r *http.Request) LocationResponse {
	// Declare a new Person struct.
	var location LocationResponse

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal()
	}

	return location
}

func locations() []Location {
	var locations = []Location{
		{"Rathaus", 1, 97},
		{"Rathaus", 2, 97},
		{"Hasse", 1, 94},
		{"Hasse", 2, 94},
		{"Pries / Friedrichsort", 1, 96},
		{"Pries / Friedrichsort", 2, 96},
		{"Neumühlen-Dietrichsdorf", 1, 92},
		{"Neumühlen-Dietrichsdorf", 2, 92},
		{"Elmschenhagen", 1, 93},
		{"Elmschenhagen", 1, 93},
		{"Mettenhof", 1, 95},
		{"Mettenhof", 1, 95},
		{"Suchsdorf", 1, 98},
		{"Suchsdorf", 1, 98},
	}

	return locations
}
