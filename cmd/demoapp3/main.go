// package main for demoapp3 implements a single endpoint that returns a simple message.
package main

import (
	"log"
	"net/http"
)

var addr = ":8080"

func main() {
	http.HandleFunc("/demoapp3", handler)
	log.Println("Starting demoapp3 server on:", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving: demoapp3")
	w.Write([]byte("demoapp3"))
}
