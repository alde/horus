package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

const (
	contentTypeJSON = "application/json; charset=UTF-8"
)

func writeJSON(status int, data interface{}, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func notFound(w http.ResponseWriter) error {
	return writeError(http.StatusNotFound, "Not Found", w)
}

func writeError(status int, message string, w http.ResponseWriter) error {
	data := make(map[string]string)
	data["error"] = message
	return writeJSON(status, data, w)
}

func handleError(err error, w http.ResponseWriter, message string) {
	if err == nil {
		return
	}

	errorMessage := struct {
		Error   error  `json:"error"`
		Message string `json:"message"`
	}{
		err, message,
	}

	if err = writeJSON(422, errorMessage, w); err != nil {
		logrus.WithError(err).WithField("message", message).Panic("Unable to respond")
	}
}

func absURL(r *http.Request, path string) string {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	uri := url.URL{
		Scheme: scheme,
		Host:   r.Host,
		Path:   path,
	}
	return uri.String()
}
