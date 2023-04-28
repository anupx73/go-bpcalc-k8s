package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/anupx73/go-bpcalc-backend-k8s/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// Get all bpReadings stored
	bpReadings, err := app.bpReadings.All()
	if err != nil {
		app.serverError(w, err)
	}

	// Convert bpReadings list into json encoding
	b, err := json.Marshal(bpReadings)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Patient records have been listed")

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	// Find patient by id
	m, err := app.bpReadings.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Patient not found")
			return
		}
		// Any other error will send an internal server error
		app.serverError(w, err)
	}

	// Convert patient bpReading to json encoding
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Found a patient record")

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Define movie model
	var m models.BPReading

	responseMsg := ""

	// Get request information
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	// Update reading time
	m.ReadingTime = time.Now()

	// Convert
	systolic, err := strconv.Atoi(m.Systolic)
	if err != nil {
		fmt.Println("Error during Systolic value conversion")
		return
	}

	diastolic, err := strconv.Atoi(m.Diastolic)
	if err != nil {
		fmt.Println("Error during Diastolic value conversion")
		return
	}

	// bp reading validation
	if systolic < 70 || systolic > 190 {
		responseMsg = "Invalid systolic value for bp reading!!"
		app.infoLog.Printf(responseMsg)
		app.clientResponse(w, 416, responseMsg)
		return
	}

	if diastolic < 40 || diastolic > 100 {
		responseMsg = "Invalid diastolic value for bp reading!!"
		app.infoLog.Printf(responseMsg)
		app.clientResponse(w, 416, responseMsg)
		return
	}

	if diastolic > systolic {
		responseMsg = "Diastolic value must be lower than systolic!!"
		app.infoLog.Printf(responseMsg)
		app.clientResponse(w, 416, responseMsg)
		return
	}

	// bp category calc
	m.Category = calcCategory(systolic, diastolic)
	if m.Category == "" {
		responseMsg = "Invalid systolic/diastolic value for bp reading!!"
		app.infoLog.Printf(responseMsg)
		app.clientResponse(w, 416, responseMsg)
		return
	}

	// Insert new bp reading
	insertResult, err := app.bpReadings.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

	// Send the response back to client
	app.clientResponse(w, 202, m.Category)

	app.infoLog.Printf("Patient record added, id=%s", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete patient by id
	deleteResult, err := app.bpReadings.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("Patient record deleted %d patient(s)", deleteResult.DeletedCount)
}
