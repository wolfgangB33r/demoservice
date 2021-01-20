package main

import (
	"bytes"
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
	Proxy          bool
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
			fmt.Printf("config payload wrong")
			log.Println("config payload is wrong")
			panic(err)
		}
		log.Println("received a new service config deployment")
	default:
		fmt.Fprintf(w, "sorry, only POST method is supported.")
	}
	defer r.Body.Close()
}

func handleIcon(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	reqcount++
	fmt.Fprintf(&buf, "it's the %d call\n", reqcount)
	fmt.Fprintf(&buf, "what I did:\n")
	// first call all callees we have in the config with the multiplicity given
	failures := false

	for ci, element := range conf.Callees {
		if !conf.Balanced || reqcount%len(conf.Callees) == ci {
			for i := 0; i < element.Count; i++ {
				req, err := http.NewRequest("GET", element.Adr, nil)
				if err != nil {
					log.Fatal("error reading request. ", err)
				}
				if conf.Proxy {
					log.Printf("dt header: %s ", r.Header.Get("X-Dynatrace"))
					log.Printf("RemoteAddr: %s ", r.RemoteAddr)
					req.Header.Set("X-Dynatrace", r.Header.Get("X-Dynatrace"))
					req.Header.Set("x-forwarded-for", r.RemoteAddr)
					req.Header.Set("forwarded", r.RemoteAddr)
				}
				req.Header.Set("Cache-Control", "no-cache")

				client := &http.Client{Timeout: time.Second * 10}

				resp, err := client.Do(req)
				if err != nil {
					log.Fatal("error reading response. ", err)
				} else {
					if resp.StatusCode != 200 {
						log.Println("got a bad return")
						failures = true
					}
				}
				defer resp.Body.Close()
			}
			fmt.Fprintf(&buf, "called %s %d times\n", element.Adr, element.Count)
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
		fmt.Fprintf(&buf, "sleeped for %d millis\n", conf.SlowdownConfig.SlowdownMillis)
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
		fmt.Fprintf(&buf, "allocated %d 100x100 matrices with random values\n", conf.ResourceConfig.Severity)
		conf.ResourceConfig.Count = conf.ResourceConfig.Count - 1
		log.Println("high resource consumption service call")
	}
	// then check if the should return an error response code
	if failures || (conf.ErrorConfig.ResponseCode != 0 && conf.ErrorConfig.Count > 0) {
		if conf.ErrorConfig.ResponseCode == 400 {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		conf.ErrorConfig.Count = conf.ErrorConfig.Count - 1
	} else {
		message := r.URL.Path
		message = strings.TrimPrefix(message, "/")
		message = "finally returned " + message
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
