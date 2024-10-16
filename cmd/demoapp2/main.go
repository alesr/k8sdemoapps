// package main for demoapp2 implements a single endpoint that returns a simple message.
package main

import (
	"log"
	"net/http"
)

var addr = ":8080"

func main() {
	http.HandleFunc("/demoapp2", handler)
	log.Println("Starting demoapp2 server on:", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving: demoapp2")
	w.Write([]byte("demoapp2"))
}
