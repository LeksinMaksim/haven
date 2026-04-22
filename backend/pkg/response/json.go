package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type envelope map[string]any

func JSON(w http.ResponseWriter, status int, data any) {
	js, err := json.Marshal(data)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	_, _ = w.Write(js)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, envelope{"error": message})
}

func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, envelope{"data": data})
}

func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, envelope{"data": data})
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, "the server encountered a problem")
}
