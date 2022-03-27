package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

	fmt.Println("Listening on :8070")
	if err := http.ListenAndServe(":8070", r); err != nil {
		//Exit
		log.Println("Failed starting server ", err)
		os.Exit(1)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is root handler")
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
	tn := time.Now().UTC()
	year, week := tn.ISOWeek()
	yearWeek := strconv.Itoa(year) + strconv.Itoa(week)
	var baseUrl = "https://terminvergabe-ema-zulassung.kiel.de/tevisema/caldiv"
	param := QueryParam{
		cal:    97,
		week:   yearWeek,
		json:   true,
		offset: true,
	}

	resp, err := http.Get(baseUrl + param.String())

	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

	return
}
