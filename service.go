package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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
	Severity int
	Count    int
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
	Balanced       bool
}

var conf config
var reqcount int

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
			log.Println("config payload is wrong")
			panic(err)
		}
		log.Println("received a new service config")
	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}
	defer r.Body.Close()
}

func handleIcon(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	log.Println("received a new service call")
	reqcount++
	fmt.Fprintf(w, "What I did %d:\n", reqcount)
	// first call all callees we have in the config with the multiplicity given
	failures := false

	for ci, element := range conf.Callees {
		if !conf.Balanced || reqcount%len(conf.Callees) == ci {
			for i := 0; i < element.Count; i++ {
				res, err := http.Get(element.Adr)
				if err == nil {
					defer res.Body.Close()
					if res.StatusCode != 200 {
						failures = true
					}
				}
			}
			fmt.Fprintf(w, "Called %s %d times\n", element.Adr, element.Count)
		}
	}
	// then check if we should crash the process
	if conf.CrashConfig.Code != 0 {
		//log.Fatalf("Exiting")
		panic("a problem")
		//os.Exit(conf.CrashConfig.Code)
	}
	// then check if we should add a delay
	if conf.SlowdownConfig.SlowdownMillis != 0 && conf.SlowdownConfig.Count > 0 {
		time.Sleep(time.Duration(conf.SlowdownConfig.SlowdownMillis) * time.Millisecond)
		conf.SlowdownConfig.Count = conf.SlowdownConfig.Count - 1
		fmt.Fprintf(w, "Sleeped for %d millis\n", conf.SlowdownConfig.SlowdownMillis)
		log.Println("slow service call")
	}
	// then check if we should increase resource consumption
	if conf.ResourceConfig.Severity != 0 && conf.ResourceConfig.Count > 0 {
		for c := 0; c <= conf.ResourceConfig.Severity; c++ {
			m1 := [100][100]int{}
			for i := 0; i < 100; i++ {
				for j := 0; j < 100; j++ {
					m1[i][j] = rand.Int()
				}
			}
		}
		fmt.Fprintf(w, "Allocated %d 100x100 matrices with random values\n", conf.ResourceConfig.Severity)
		conf.ResourceConfig.Count = conf.ResourceConfig.Count - 1
		log.Println("high resource consumption service call")
	}
	// then check if the should return an error response code
	if failures || (conf.ErrorConfig.ResponseCode != 0 && conf.ErrorConfig.Count > 0) {
		if conf.ErrorConfig.ResponseCode == 400 {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("400 - Forbidden!"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
		}
		fmt.Fprintf(w, "Returned an error response code\n")
		conf.ErrorConfig.Count = conf.ErrorConfig.Count - 1
		log.Println("failed service call")
	} else {
		message := r.URL.Path
		message = strings.TrimPrefix(message, "/")
		message = "Finally returned " + message
		w.Write([]byte(message))
	}
	defer r.Body.Close()
}

func main() {
	port := 8080
	if len(os.Args) > 1 {
		arg := os.Args[1]
		fmt.Printf("Start demo service at port: %s\n", arg)
		i1, err := strconv.Atoi(arg)
		if err == nil {
			port = i1
		}
	} else {
		fmt.Printf("Start demo service at default port: %d\n", port)
	}

	http.HandleFunc("/", sayHello)
	http.HandleFunc("/favicon.ico", handleIcon)
	http.HandleFunc("/config", receiveConfig)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		panic(err)
	}
}
