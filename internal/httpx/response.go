package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go-api/internal/gen/openapi"
)

func DecodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(new(struct{})); !errors.Is(err, io.EOF) {
		return http.ErrBodyNotAllowed
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, statusCode int, code string, message string, details []openapi.ErrorDetail) {
	var responseDetails *[]openapi.ErrorDetail
	if len(details) > 0 {
		responseDetails = &details
	}

	WriteJSON(w, statusCode, openapi.ErrorResponse{
		Code:    code,
		Message: message,
		Details: responseDetails,
	})
}
