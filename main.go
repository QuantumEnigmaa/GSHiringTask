package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

func welcomeFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the beeant web application ! You will find here many information about bees and ants !")
}

func beeFunc(w http.ResponseWriter, r *http.Request) {
	rep := generalResponse{
		Title:   "What is a Bee ?",
		Message: "A bee is a wonderful insect, maybe one of the most important creature living in the surface.",
	}
	json.NewEncoder(w).Encode(rep)
}

func antFunc(w http.ResponseWriter, r *http.Request) {
	rep := generalResponse{
		Title:   "What is an Ant ?",
		Message: "An ant is a fascinating insect, living in a colony.",
	}
	json.NewEncoder(w).Encode(rep)
}

func checkFunc(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rep checkResponse
	json.Unmarshal(reqBody, &rep)

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

func handleRequest() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", welcomeFunc)
	router.HandleFunc("/bee", beeFunc)
	router.HandleFunc("/ant", antFunc)
	router.HandleFunc("/check", checkFunc).Methods("POST")

	log.Fatal(http.ListenAndServe(":8088", router))
}

func main() {
	fmt.Println("Beeant web application is running and listening on port 8088 !")
	handleRequest()
}
