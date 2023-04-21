package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientResponse(w http.ResponseWriter, status int, responseTxt string) {

	data := `{"status":"` + strconv.Itoa(status) + `","message":"` + responseTxt + `"}`
	jsData := []byte(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsData)
}
