package infra

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type HTTPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	enc := json.NewEncoder(w)
	enc.Encode(v)
}

func InternalServerError(w http.ResponseWriter, err error, source string) {
	log.Error().Stack().Err(err).Str("source", source).Msg("")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	writeJSON(w, HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
	})
}

func BadRequestError(w http.ResponseWriter, details interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	writeJSON(w, HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Bad Request",
		Details: details,
	})
}

func Created(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, data)
}

func OK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	writeJSON(w, data)
}

func NoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

func NotFoundError(w http.ResponseWriter, id string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	writeJSON(w, HTTPError{
		Code:    http.StatusNotFound,
		Message: "Resource Not Found",
		Details: map[string]string{"id": id},
	})
}
