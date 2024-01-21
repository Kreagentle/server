package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCRUDOperations(t *testing.T) {
	router := mux.NewRouter()

	// Test Create
	createData := TestData{FirstElement: "abc", SecondElement: 123}
	createJSON, _ := json.Marshal(createData)
	createReq, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(createJSON))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, createRec.Code)
	}

	// Test Read
	readReq, _ := http.NewRequest("GET", "/read", nil)
	readRec := httptest.NewRecorder()
	router.ServeHTTP(readRec, readReq)

	if readRec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, readRec.Code)
	}

	// Test Update
	updateData := TestData{FirstElement: "abc", SecondElement: 234}
	updateJSON, _ := json.Marshal(updateData)
	updateReq, _ := http.NewRequest("UPDATE", "/update", bytes.NewBuffer(updateJSON))
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, readRec.Code)
	}

	// Test Delete
	deleteReq, _ := http.NewRequest("DELETE", "/delete/0", nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, deleteRec.Code)
	}

	// Example of checking the response body in the Read test
	var responseData []TestData
	err := json.Unmarshal(readRec.Body.Bytes(), &responseData)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Example of asserting expected data in the response
	if len(responseData) != 1 || responseData[0].FirstElement != "abc" || responseData[0].SecondElement != 123 {
		t.Errorf("Unexpected response data: %v", responseData)
	}
}
