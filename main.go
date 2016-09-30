package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type input struct {
	Name  string   `json:"name"`
	Names []string `json:"names"`
}

type output struct {
	Result  string   `json:"result,omitempty"`
	Results []string `json:"results,omitempty"`
}

func main() {
	// Read config from environment
	bindAddr := os.Getenv("BIND_ADDR")
	if len(bindAddr) == 0 {
		bindAddr = ":8080"
	}

	fmt.Printf("Starting server on %s\n", bindAddr)

	// Setup some routing
	http.HandleFunc("/", wrapper(handler))

	// Create a HTTPS listener
	if err := http.ListenAndServeTLS(bindAddr, "cert.pem", "key.pem", nil); err != nil {
		log.Fatal(err)
	}
}

func wrapper(f http.HandlerFunc) http.HandlerFunc {
	// Return a HTTP handler which calculates the request processing time
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		f(w, req)
		diff := time.Now().Sub(start)
		log.Printf("duration: %+v\n", diff)
	}
}

func generate(input ...string) (output []string) {
	for _, name := range input {
		output = append(output, fmt.Sprintf("Hello %s", name))
	}
	return
}

func handler(w http.ResponseWriter, req *http.Request) {
	// Read the request body
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	req.Body.Close()

	// Decode the JSON into a struct
	var i input
	err = json.Unmarshal(b, &i)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	// Sanity check the input
	if len(i.Name) == 0 && len(i.Names) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("must specify 'name' or 'names'"))
		return
	}

	// Generate the 'Hello {name}' results
	var o output
	if len(i.Name) > 0 {
		o.Result = generate(i.Name)[0]
	} else {
		o.Results = generate(i.Names...)
	}

	// Encode the struct as JSON
	b, err = json.Marshal(&o)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// Return the encoded JSON to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}
