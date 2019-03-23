package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PutAndGet(t *testing.T) {
	a := assert.New(t)
	tmpdir := os.TempDir()
	db, _ := NewFilestore(tmpdir)

	err := db.Put("github.com/alde/horus", "DOCKER_LOGIN", []byte("totally_secret"))
	a.Nil(err, "error should be nil")
	err = db.Put("github.com/alde/horus", "DOCKER_PASSWORD", []byte("a different secret"))
	a.Nil(err, "error should be nil")

	secret, err := db.Get("github.com/alde/horus", "DOCKER_LOGIN")
	a.Nil(err, "error should be nil")
	a.Equal("totally_secret", secret, "the stored and retrieved secret should match")

	secrets := db.List("github.com/alde/horus")
	a.Len(secrets, 2, "should have 2 secrets")
	a.Equal(secrets, []string{"DOCKER_LOGIN", "DOCKER_PASSWORD"})

	err = db.Remove("github.com/alde/horus", "DOCKER_LOGIN")
	a.Nil(err, "error should be nil")

	has := db.Has("github.com/alde/horus", "DOCKER_LOGIN")
	a.False(has, "should no longer exist")
}
