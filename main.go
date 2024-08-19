package main

import (
	"encoding/json"
	"log"
	"net/http"
//	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"

	. "github.com/AnastasiyaGapochkina01/go-api/model"
	. "github.com/AnastasiyaGapochkina01/go-api/mongo"
)

var dbo = MongoConfig{}

// GET list of wods
func AllWodsEndPoint(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
	wods, err := dbo.getWods(params["status"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, wods)
}

// GET a wod by its ID
func FindWodEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	wod, err := dbo.getWod(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Wod ID")
		return
	}
	respondWithJson(w, http.StatusOK, wod)
}

// POST a new wod
func CreateWodEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var wod Wod
	if err := json.NewDecoder(r.Body).Decode(&wod); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	//wod.Id = bson.NewObjectId()
	if err := dbo.putWod(wod); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, wod)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/wods", AllWodsEndPoint).Methods("GET")
	
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
