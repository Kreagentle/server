package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
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

var logger *slog.Logger
var database map[int]TestData
var counter int

func create(w http.ResponseWriter, r *http.Request) {
	var data TestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logger.Error("Invalid JSON format: %s\n", counter)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	database[counter] = data
	counter += 1

	w.WriteHeader(http.StatusCreated)
	logger.Info("Data created with key: %s\n", counter)
}

func read(w http.ResponseWriter, r *http.Request) {
	logger.Info("Reading is started")
	var responseData []TestData
	for _, data := range database {
		responseData = append(responseData, data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var data TestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logger.Error("Invalid JSON format: %s\n", counter)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	nums, err := strconv.Atoi(key)
	if err == nil {
		logger.Error("Cant convert num: %s\n", counter)
		http.Error(w, "Cant convert num", http.StatusBadRequest)
		return
	}

	if _, ok := database[nums]; !ok {
		logger.Error("Data not found: %s\n", counter)
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	database[counter] = data

	w.WriteHeader(http.StatusOK)
	logger.Info("Data updated for key: %s\n", key)
}

func delete_(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	num, err := strconv.Atoi(key)
	if err == nil {
		logger.Error("Cant convert num: %d\n", counter)
		http.Error(w, "Cant convert num", http.StatusBadRequest)
		return
	}

	if _, ok := database[num]; !ok {
		logger.Error("Data not found %s\n", counter)
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	delete(database, num)

	w.WriteHeader(http.StatusOK)
	logger.Info("Data deleted for key: %s\n", key)
}

func init() {
	// logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	logger = slog.New(handler)
}

func main() {
	// config
	f, err := os.Open("config.yml")
	if err != nil {
		logger.Error("Cant open config yaml", err)
		return
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Error("Cant decode config yaml", err)
		return
	}

	// router
	router := mux.NewRouter()

	router.HandleFunc("/create", create).Methods("POST")
	router.HandleFunc("/read", read).Methods("GET")
	router.HandleFunc("/update", update).Methods("PUT")
	router.HandleFunc("/delete", delete_).Methods("DELETE")

	// server
	err = http.ListenAndServe(":"+cfg.Server.Port, router)
	if err != nil {
		logger.Error("Problems with server: ", err)
	}
	logger.Info("server " + cfg.Server.Port + " is started")
}
