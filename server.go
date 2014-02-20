package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	listen string
)

func post(w http.ResponseWriter, r *http.Request) {
	var m map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		panic(err)
	}

	m["ip"] = r.RemoteAddr
	m["date"] = time.Now()
	json.NewEncoder(os.Stdout).Encode(m)
}

func main() {
	flag.StringVar(&listen, "listen", ":8000", "Listen address")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/post", post)

	log.Fatal(http.ListenAndServe(listen, nil))
}
