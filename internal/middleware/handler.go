package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"micro-blog/internal/exception"
)

var (
	errNotFound            *exception.NotFoundError
	errUnauthorized        *exception.UnauthorizedError
	errBadRequest          *exception.BadRequestError
	errConflict            *exception.ConflictError
	errUnprocessableEntity *exception.UnprocessableEntityError
)

type HandleWrapper func(w http.ResponseWriter, r *http.Request) error

func (h HandleWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if rec := recover(); rec != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, map[string]string{"msg": fmt.Sprintf("%v", rec)})
		}
	}()

	err := h(w, r)
	if err == nil {
		return
	}

	w.WriteHeader(statusCode(err))

	if _, ok := err.(json.Marshaler); ok {
		writeJSON(w, err)
		return
	}
	writeJSON(w, map[string]string{"msg": err.Error()})
}

func writeJSON(w http.ResponseWriter, data any) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write json response: %v", err)
	}
}

func statusCode(err error) int {
	switch {
	case errors.As(err, &errBadRequest):
		return http.StatusBadRequest
	case errors.As(err, &errUnauthorized):
		return http.StatusUnauthorized
	case errors.As(err, &errNotFound):
		return http.StatusNotFound
	case errors.As(err, &errConflict):
		return http.StatusConflict
	case errors.As(err, &errUnprocessableEntity):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
