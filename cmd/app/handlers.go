package main

import (
	"encoding/json"
	"net/http"
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

	// Convert movie list into json encoding
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
	// Get request information
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	// Update reading time
	m.ReadingTime = time.Now()

	// bp reading validation
	if m.Systolic < 70 || m.Systolic > 190 {
		app.infoLog.Printf("Invalid systolic value for bp reading!!")
		app.clientError(w, 416)
		return
	}

	if m.Diastolic < 40 || m.Diastolic > 100 {
		app.infoLog.Printf("Invalid diastolic value for bp reading!!")
		app.clientError(w, 416)
		return
	}

	if m.Diastolic > m.Systolic {
		app.infoLog.Printf("Diastolic value must be lower than systolic!!")
		app.clientError(w, 416)
		return
	}

	// bp category calc
	if m.Systolic < 90 && m.Diastolic < 60 {
		m.Category = "Low"
	} else if m.Systolic < 120 && m.Diastolic < 80 {
		m.Category = "Ideal"
	} else if m.Systolic < 140 && m.Diastolic < 90 {
		m.Category = "Pre High"
	} else if m.Systolic <= 190 && m.Diastolic <= 100 {
		m.Category = "High"
	} else {
		app.infoLog.Printf("Invalid systolic/diastolic value for bp reading!!")
		app.clientError(w, 416)
		return
	}

	// Insert new bp reading
	insertResult, err := app.bpReadings.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

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
