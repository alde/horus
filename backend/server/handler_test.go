package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/alde/horus/backend/config"
	"github.com/alde/horus/backend/encryptor"
	"github.com/alde/horus/backend/mock"
)

var cfg0 = config.DefaultConfig()
var enc, _ = encryptor.NewGoogleCloudKMS(context.Background(), cfg0, mock.Encryptor{})

func Test_ServiceMetadata(t *testing.T) {
	m := mux.NewRouter()

	h := NewHandler(cfg0, mock.DB{}, enc)
	m.HandleFunc("/service-metadata", h.ServiceMetadata())
	wr := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/service-metadata", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, wr.Code, http.StatusOK)

	var actual map[string]interface{}
	err := json.Unmarshal(wr.Body.Bytes(), &actual)
	assert.Nil(t, err)

	expectedKeys := []string{
		"service_name", "service_version", "description", "owner",
	}

	for _, k := range expectedKeys {
		_, ok := actual[k]
		assert.True(t, ok, "expected key %s", k)
	}
}

func TestHandler_Fetch(t *testing.T) {
	t.Skip("not implemented")
	m := mux.NewRouter()
	h := NewHandler(cfg0, mock.DB{}, enc)
	m.HandleFunc("/api/v1/secret", h.Fetch())
	wr := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/api/v1/secret?repo=github.com%2Falde%2Fhorus&key=SECRET", nil)
	m.ServeHTTP(wr, r)

	assert.Equal(t, http.StatusOK, wr.Code)
}
