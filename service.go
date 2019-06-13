package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var slowrequests = 0
var errorrequests = 0

func sayHello(w http.ResponseWriter, r *http.Request) {
	today := time.Now()
	fmt.Println("Hour: ", today.Hour())

	//
	slowdown := r.URL.Query().Get("slowdown")
	if slowdown != "" {
		slowrequests = 200
	}

	//
	error := r.URL.Query().Get("error")
	if error != "" {
		errorrequests = 200
	}

	// artificial slowdown in case of a slowdown problem
	if slowrequests > 0 {
		time.Sleep(200 * time.Millisecond)
		// start some CPU/mem intensive operation
		m1 := [1000][1000]int{}
		for i := 0; i < 1000; i++ {
			for j := 0; j < 1000; j++ {
				m1[i][j] = rand.Int()
			}
		}

		// End of CPU intensive operation
		slowrequests = slowrequests - 1
		fmt.Println("Slowdown pattern active: ", slowrequests)
	}

	if today.Hour() == 10 { // trigger error pattern between 11 and 12 every day
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	} else if errorrequests > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		errorrequests = errorrequests - 1
		fmt.Println("Error pattern active: ", errorrequests)
	} else {
		message := r.URL.Path
		message = strings.TrimPrefix(message, "/")
		message = "Hello " + message
		w.Write([]byte(message))
	}
	defer r.Body.Close()
}

func main() {
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
