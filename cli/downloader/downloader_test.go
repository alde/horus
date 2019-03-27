package downloader

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alde/horus/cli/config"
	"github.com/stretchr/testify/assert"
)

func Test_Download(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), `/api/v1/secret?repo=github.com%2Falde%2Fhorus&key=SECRET`)
		rw.Write([]byte(`a-secret`))
	}))
	// Close the server when test finishes
	defer server.Close()
	cfg := &config.Config{}
	cfg.Horus.Host = server.URL

	d := New(server.Client(), cfg)

	secret := d.Download("github.com/alde/horus", "SECRET")
	assert.Equal(t, "a-secret", secret)
}
