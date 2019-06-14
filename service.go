package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type errorAnomalyConfig struct {
	ResponseCode int
	Count        int
}

type slowdownAnomalyConfig struct {
	SlowdownMillis int
	Count          int
}

type crashAnomalyConfig struct {
	Code int
}

type resourceAnomalyConfig struct {
}

type callee struct {
	Adr   string // URL address to call
	Count int    // number of calls per minute
}

type config struct {
	ErrorConfig    errorAnomalyConfig
	SlowdownConfig slowdownAnomalyConfig
	CrashConfig    crashAnomalyConfig
	ResourceConfig resourceAnomalyConfig
	Callees        []callee
}

var conf config

func receiveConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusNoContent)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		//fmt.Printf(string(body))
		err = json.Unmarshal(body, &conf)
		if err != nil {
			fmt.Printf("Config payload wrong")
			panic(err)
		}
	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}
	defer r.Body.Close()
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	// first call all callees we have in the config with the multiplicity given
	for _, element := range conf.Callees {
		fmt.Printf("Call %s %d times\n", element.Adr, element.Count)
		for i := 0; i < element.Count; i++ {
			res, err := http.Get(element.Adr)
			if err == nil {
				defer res.Body.Close()
			}
		}
	}
	// then check if we should crash the process
	if conf.CrashConfig.Code != 0 {
		os.Exit(conf.CrashConfig.Code)
	}
	// then check if we should add a delay
	if conf.SlowdownConfig.SlowdownMillis != 0 && conf.SlowdownConfig.Count > 0 {
		time.Sleep(time.Duration(conf.SlowdownConfig.SlowdownMillis) * time.Millisecond)
		conf.SlowdownConfig.Count = conf.SlowdownConfig.Count - 1
	}
	// then check if we should increase resource consumption

	// then check if the should return an error response code

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
