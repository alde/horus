package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Filestore_FullCycle(t *testing.T) {
	tmpdir := os.TempDir()
	db, _ := NewFilestore(tmpdir)

	err := db.Put("github.com/alde/horus", "DOCKER_LOGIN", []byte("totally_secret"))
	assert.Nil(t, err, "error should be nil")
	err = db.Put("github.com/alde/horus", "DOCKER_PASSWORD", []byte("a different secret"))
	assert.Nil(t, err, "error should be nil")

	secret, err := db.Get("github.com/alde/horus", "DOCKER_LOGIN")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, "totally_secret", secret, "the stored and retrieved secret should match")

	secrets := db.List("github.com/alde/horus")
	assert.Len(t, secrets, 2, "should have 2 secrets")
	assert.Equal(t, secrets, []string{"DOCKER_LOGIN", "DOCKER_PASSWORD"})

	err = db.Remove("github.com/alde/horus", "DOCKER_LOGIN")
	assert.Nil(t, err, "error should be nil")

	has := db.Has("github.com/alde/horus", "DOCKER_LOGIN")
	assert.False(t, has, "should no longer exist")
}
