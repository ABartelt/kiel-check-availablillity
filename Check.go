package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type QueryParam struct {
	cal    int
	week   string
	json   bool
	offset bool
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/check", checkAll).Methods("GET")

	port := goDotEnvVariable("PORT")
	fmt.Println("Listening on :" + port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		//Exit
		log.Println("Failed starting server ", err)
		os.Exit(1)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is root handler")
}

// Use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func (param QueryParam) String() string {
	var query = "&"
	query = query + "cal=" + strconv.Itoa(param.cal)
	query = query + "&week=" + param.week
	query = query + "&json=" + "1"
	query = query + "&offset=" + "1"

	return query
}

func checkAll(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	chResponses := make(chan LocationResponse)
	start := time.Now()
	year, week := start.ISOWeek()
	var periodOfweeks int = 8

	for i := 0; i < periodOfweeks; i++ {
		checkLocations(year, week, periodOfweeks, i, chResponses, start, wg)
	}

	var responses []LocationResponse

	for {
		res, ok := <-chResponses
		if !ok {
			fmt.Println("Channel Close ", ok)
			break
		}
		responses = append(responses, res)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responses)
}

func checkLocations(year int, week int, periodOfweeks int, i int, chWeekResponses chan LocationResponse, start time.Time, wg sync.WaitGroup) {
	for _, location := range locations() {
		var actualWeek int
		year, actualWeek = adjustYear(week, periodOfweeks, year, i)
		go makeRequest(location, chWeekResponses, start, year, actualWeek, &wg)
	}
	close(chWeekResponses)
}

func adjustYear(week int, periodOfweeks int, year int, i int) (int, int) {
	if week >= (52 - periodOfweeks) {
		year = year + 1
		week = 52 - (52 - periodOfweeks)
	}

	var actualWeek int = week + i

	return year, actualWeek
}

func makeRequest(location Location, ch chan<- LocationResponse, start time.Time, year int, week int, wg *sync.WaitGroup) {
	baseUrl := "https://terminvergabe-ema-zulassung.kiel.de/tevisema/caldiv"
	url := generateUrl(baseUrl, location, year, week)
	resp, _ := http.Get(url)
	secs := time.Since(start).Seconds()

	// Parse []byte to the go struct pointer
	ch <- parseLocationResponse(resp, location, week, secs)
	defer wg.Done()
}

func generateUrl(baseUrl string, location Location, year int, week int) string {
	url := baseUrl + QueryParam{
		cal:  location.Cal,
		week: strconv.Itoa(year) + strconv.Itoa(week),
	}.String()

	return url
}
