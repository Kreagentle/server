package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"main/src"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCRUDOperations(t *testing.T) {
	f, err := os.Open("config.yml")
	if err != nil {
		src.Logger.Error("Cant open config yaml", err)
		return
	}
	defer f.Close()

	var cfg src.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		src.Logger.Error("Cant decode config yaml", err)
		return
	}

	router := src.SetupRouter()

	// Test Create
	createData := src.TestData{FirstElement: "abc", SecondElement: 123}
	createJSON, _ := json.Marshal(createData)
	createReq, _ := http.NewRequest("POST", "http://localhost:8000/create", bytes.NewBuffer(createJSON))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, createRec.Code)
	}

	// Test Read
	readReq, _ := http.NewRequest("GET", "http://localhost:8000/read", nil)
	readRec := httptest.NewRecorder()
	router.ServeHTTP(readRec, readReq)

	if readRec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, readRec.Code)
	}

	// Test Update
	updateData := src.TestData{FirstElement: "abc", SecondElement: 234}
	updateJSON, _ := json.Marshal(updateData)
	updateReq, _ := http.NewRequest("PUT", "http://localhost:8000/update/1", bytes.NewBuffer(updateJSON))
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, updateRec.Code)
	}

	// Test Delete
	deleteReq, _ := http.NewRequest("DELETE", "http://localhost:8000/delete/1", nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, deleteRec.Code)
	}
}
