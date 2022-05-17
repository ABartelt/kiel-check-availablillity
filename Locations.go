package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func parseLocationResponse(resp *http.Response, location Location, week int, secs float64) LocationResponse {
	// Declare a new LocationResponse struct.
	var requestResult LocationResponse

	// Checking http status code.
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Printf("fetched %s for week %d in %.2f successfully", location.Name, week, secs)
	} else {
		fmt.Printf("%d: Location %s could not be fetched for week %d", resp.StatusCode, location.Name, week)
	}

	// Error is not catched here.
	body, _ := ioutil.ReadAll(resp.Body)

	// Map the body byte [] into requestResult.
	if err := json.Unmarshal(body, &requestResult); err != nil {
		fmt.Printf("Can not unmarshal JSON")
	}

	return requestResult
}
