package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	listen  string
	dir     string
	verbose bool
	out     string
	outMut  sync.Mutex
)

func logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func post(w http.ResponseWriter, r *http.Request) {
	var m map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		panic(err)
	}

	m["ip"] = r.RemoteAddr
	m["date"] = time.Now()
	appendJSON(m)
}

func appendJSON(m map[string]interface{}) {
	var f io.WriteCloser
	if len(out) > 0 {
		outMut.Lock()
		defer outMut.Unlock()

		var err error
		f, err = os.OpenFile(out, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	} else {
		f = os.Stdout
	}

	err := json.NewEncoder(f).Encode(m)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.StringVar(&listen, "listen", ":8000", "Listen address")
	flag.StringVar(&dir, "dir", "static", "Files location")
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.StringVar(&out, "out", "", "Output file")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.HandleFunc("/post", post)

	if verbose {
		log.Fatal(http.ListenAndServe(listen, logging(http.DefaultServeMux)))
	} else {
		log.Fatal(http.ListenAndServe(listen, http.DefaultServeMux))
	}
}
