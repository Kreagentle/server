package src

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
}

type TestData struct {
	FirstElement  string `json:"firstel"`
	SecondElement int    `json:"secondel"`
}

var (
	database = make(map[int]TestData)
	counter  = 1
)

var Logger *slog.Logger

func Create(w http.ResponseWriter, r *http.Request) {
	var data TestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		Logger.Error("Invalid JSON format: %s\n", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	database[counter] = data
	counter += 1

	w.WriteHeader(http.StatusCreated)
}

func Read(w http.ResponseWriter, r *http.Request) {
	var responseData []TestData
	for _, data := range database {
		responseData = append(responseData, data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var data TestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		Logger.Error("Invalid JSON format: %s\n", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	nums, err := strconv.Atoi(key)
	if err != nil {
		Logger.Error("Cant convert num: %s\n", err)
		http.Error(w, "Cant convert num", http.StatusBadRequest)
		return
	}

	if _, ok := database[nums]; !ok {
		Logger.Error("Data not found\n")
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	database[nums] = data

	w.WriteHeader(http.StatusOK)
}

func Delete_(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	num, err := strconv.Atoi(key)
	if err != nil {
		Logger.Error("Cant convert num: %s\n", err)
		http.Error(w, "Cant convert num", http.StatusBadRequest)
		return
	}

	if _, ok := database[num]; !ok {
		Logger.Error("Data not found %s\n", err)
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	delete(database, num)

	w.WriteHeader(http.StatusOK)
}

func SetupRouter() *mux.Router {
	// router
	router := mux.NewRouter()

	router.HandleFunc("/create", Create).Methods("POST")
	router.HandleFunc("/read", Read).Methods("GET")
	router.HandleFunc("/update/{key}", Update).Methods("PUT")
	router.HandleFunc("/delete/{key}", Delete_).Methods("DELETE")

	return router
}
