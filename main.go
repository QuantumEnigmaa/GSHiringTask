package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

////////////////////////////////
// Structs declaration
////////////////////////////////

type generalResponse struct {
	Title   string `json:"Title"`
	Message string `json:"Message"`
}

type checkResponse struct {
	Wings bool   `json:"Wings"`
	Legs  int    `json:"Legs"`
	Color string `json:"Color"`
	Alone bool   `json:"Alone"`
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

////////////////////////////////
// End of structs declaration
////////////////////////////////

////////////////////////////////
// Prometheus metrics declaration
////////////////////////////////

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Number of requests recorded",
	},
	[]string{"path"})

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"})

var responseTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Response time for http requests for webapp",
}, []string{"path"})

////////////////////////////////
// end of Prometheus metrics declaration
////////////////////////////////

func welcomeFunc(w http.ResponseWriter, r *http.Request) {
	// Handler of the "/" endpoint : greet the user to the application
	fmt.Fprint(w, "Welcome to the beeant web application ! You will find here many information about bees and ants !")
}

func beeFunc(w http.ResponseWriter, r *http.Request) {
	// handler of the "/bee" endpoint : describe what a bee is
	rep := generalResponse{
		Title:   "What is a Bee ?",
		Message: "A bee is a wonderful insect, maybe one of the most important creature living in the surface.",
	}
	json.NewEncoder(w).Encode(rep)
}

func antFunc(w http.ResponseWriter, r *http.Request) {
	// Handler of the "/ant" endpoint : describe what an ant is
	rep := generalResponse{
		Title:   "What is an Ant ?",
		Message: "An ant is a fascinating insect, living in a colony.",
	}
	json.NewEncoder(w).Encode(rep)
}

func checkFunc(w http.ResponseWriter, r *http.Request) {
	// Handler of the "/check" endpoint : check if the insect described is a bee, an ant, or else...
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Deserialization of the POST request body (JSON format) into custom struct
	var rep checkResponse
	json.Unmarshal(reqBody, &rep)

	// Declaration of all possible answers
	isBee := generalResponse{
		Title:   "Analysis done",
		Message: "Your insect is a bee",
	}
	isAnt := generalResponse{
		Title:   "Analysis done",
		Message: "Your insect is an ant",
	}
	isBeeant := generalResponse{
		Title:   "Analysis done",
		Message: "It is difficult to recognize your insect. It might be one of the fabled beeants",
	}

	// Basic handling of the analysis of the POST request based on its fields
	if rep.Legs != 6 {
		json.NewEncoder(w).Encode(isBeeant)
	} else if rep.Wings {
		if rep.Color == "yellow" || rep.Color == "yellow and black" || rep.Color == "black and yellow" {
			json.NewEncoder(w).Encode(isBee)
		} else if rep.Color == "black" {
			json.NewEncoder(w).Encode(isAnt)
		} else {
			json.NewEncoder(w).Encode(isBeeant)
		}
	} else if !rep.Alone {
		json.NewEncoder(w).Encode(isAnt)
	} else {
		json.NewEncoder(w).Encode(isBeeant)
	}
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, err := route.GetPathTemplate()
		if err != nil {
			log.Fatal(err)
		}
		// Adding status code to request
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode
		timer := prometheus.NewTimer(responseTime.WithLabelValues(path))

		// Updating each metric
		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()
		timer.ObserveDuration()
	})
}

func handleRequest() {
	router := mux.NewRouter().StrictSlash(true)
	// Each request goes through the promethus middleware before being forwared to the right endpoint
	router.Use(prometheusMiddleware)

	// Appplication endpoints
	router.HandleFunc("/", welcomeFunc)
	router.HandleFunc("/bee", beeFunc)
	router.HandleFunc("/ant", antFunc)
	router.HandleFunc("/check", checkFunc).Methods("POST")

	// Prometheus endpoint for exposing metrics
	router.Path("/metrics").Handler(promhttp.Handler())

	// Application listening on port 8088
	log.Fatal(http.ListenAndServe(":8088", router))
}

func init() {
	// Initialize each prometheus metrics
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(responseTime)
}

func main() {
	fmt.Println("Beeant web application is running and listening on port 8088 !")
	handleRequest()
}
