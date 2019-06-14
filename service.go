package main

import (
	"fmt"
	"net/http"
	"strings"
)

func receiveConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusNoContent)

	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")

		//w.Write([]byte("500 - Something bad happened!"))
	}
	defer r.Body.Close()
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
	//fmt.Printf("GET / request processed.")
	defer r.Body.Close()
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/config", receiveConfig)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
