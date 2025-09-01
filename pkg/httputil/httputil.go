package httputil

import (
	"encoding/json"
	"log/slog"
	"maps"
	"net/http"
)

func LogError(logger *slog.Logger, r *http.Request, err error) {
	logger.Error(err.Error(),
		"method", r.Method,
		"uri", r.URL.RequestURI())
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}
