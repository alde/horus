package server

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/alde/horus/backend/config"
	"github.com/alde/horus/backend/database"
	"github.com/alde/horus/backend/encryptor"
	"github.com/alde/horus/backend/version"
)

// Handler holds the server context
type Handler struct {
	config *config.Config
	db     database.Database
	enc    encryptor.Encryptor
}

// NewHandler createss a new HTTP handler
func NewHandler(cfg *config.Config, db database.Database, enc encryptor.Encryptor) *Handler {
	return &Handler{config: cfg, db: db, enc: enc}
}

// ServiceMetadata displays hopefully useful information about the service
func (h *Handler) ServiceMetadata() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := h.config.Metadata
		data["description"] = "horus stores encrypted secrets"
		data["service_name"] = "horus"
		data["service_version"] = version.Version
		data["build_date"] = version.BuildDate

		_ = writeJSON(http.StatusOK, data, w)
	}
}

type storageRequest struct {
	Repo   string `json:"repo"`
	Secret string `json:"secret"`
	Key    string `json:"key"`
}

// Store saves a secret
func (h *Handler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request storageRequest
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			handleError(err, w, "unable to read payload")
			return
		}
		err = json.Unmarshal(body, &request)
		if err != nil {
			handleError(err, w, "unable to deserialize payload")
			return
		}
		logrus.WithFields(logrus.Fields{
			"repo": request.Repo,
			"key":  request.Key,
		}).Info("adding secret")

		encrypted, err := h.enc.Encrypt([]byte(request.Secret))
		b64 := base64.StdEncoding.EncodeToString(encrypted)
		if err != nil {
			logrus.WithError(err).Error("unable to encrypt secret")
			handleError(err, w, "unable to encrypt secret")
			return
		}
		err = h.db.Put(request.Repo, request.Key, []byte(b64))
		if err != nil {
			logrus.WithError(err).Error("unable to store secret")
			handleError(err, w, "unable to store secret")
			return
		}
		logrus.WithFields(logrus.Fields{
			"repo": request.Repo,
			"key":  request.Key,
		}).Info("secret added")
		_ = writeJSON(http.StatusCreated, nil, w)
	}
}

// Fetch returns an encrypted secret
func (h *Handler) Fetch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		repo, err := url.QueryUnescape(query.Get("repo"))
		if err != nil {
			logrus.WithError(err).Debug("unable to parse repo query parameter")
			handleError(err, w, "unable to parse repo query parameter")
			return
		}
		key, err := url.QueryUnescape(query.Get("key"))
		if err != nil {
			logrus.WithError(err).Debug("unable to parse key query parameter")
			handleError(err, w, "unable to parse key query parameter")
			return
		}
		logrus.WithFields(logrus.Fields{
			"key":  key,
			"repo": repo,
		}).Info("attempting to retrieve secret")

		if exists := h.db.Has(repo, key); !exists {
			logrus.WithFields(logrus.Fields{
				"key":  key,
				"repo": repo,
			}).Debug("no such secret")
			_ = writeJSON(
				http.StatusNotFound,
				struct{ Message string }{"no such secret exists"},
				w)
			return
		}

		secret, err := h.db.Get(repo, key)
		if err != nil {
			logrus.WithError(err).Error("unable to retrieve secret")
			handleError(err, w, "unable to retrieve error")
		}
		logrus.WithFields(logrus.Fields{
			"key":  key,
			"repo": repo,
		}).Info("fetched secret")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(secret)

	}
}

// List returns a list of keys available for a certain repo
func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		repo, err := url.QueryUnescape(query.Get("repo"))
		if err != nil {
			logrus.WithError(err).Debug("unable to parse repo query parameter")
			handleError(err, w, "unable to parse repo query parameter")
			return
		}
		secrets := h.db.List(repo)
		if len(secrets) == 0 {
			_ = writeJSON(
				http.StatusNotFound,
				struct {
					Message string `json:"message"`
				}{"no secrets found"},
				w)
			return
		}
		logrus.WithFields(logrus.Fields{
			"repo": repo,
		}).Info("listing secrets (keys only)")

		_ = writeJSON(http.StatusOK, secrets, w)
	}
}
