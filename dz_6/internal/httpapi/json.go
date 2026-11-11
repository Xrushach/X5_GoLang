package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

const maxBody = 1024 // Небольшой лимит позволяет легко проверить 413.

type inputError struct {
	status  int
	message string
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) *inputError {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return &inputError{http.StatusUnsupportedMediaType, "Content-Type must be application/json"}
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return decodeError(err)
	}
	// Сколько JSON-значений прочитал первый Decode?
	if err := decoder.Decode(new(any)); err != io.EOF {
		if err != nil {
			return decodeError(err)
		}
		return &inputError{http.StatusBadRequest, "body must contain one JSON value"}
	}
	return nil
}

func decodeError(err error) *inputError {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		return &inputError{http.StatusRequestEntityTooLarge, "body is too large"}
	}
	return &inputError{http.StatusBadRequest, "invalid JSON body"}
}
