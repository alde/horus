package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alde/horus/backend/config"
	"github.com/alde/horus/backend/mock"
)

var (
	cfg          = config.DefaultConfig()
	mockDatabase = &mock.DB{}
)

func Test_NewRouter(t *testing.T) {
	h := NewHandler(cfg, mockDatabase, enc)
	nr := NewRouter(cfg, mockDatabase, enc)

	for _, r := range routes(h) {
		assert.NotNil(t, nr.GetRoute(r.Name))
	}
}

func Test_routes(t *testing.T) {
	h := NewHandler(cfg, mockDatabase, enc)
	assert.Len(t, routes(h), 4, "4 routes is the magic number.")
}
