package database

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Filestore is a hacky filesystem "database". To be removed.
type Filestore struct {
	folder string
}

// NewFilestore creates a new Datastore database object
func NewFilestore(folder string) (Database, error) {
	return &Filestore{
		folder: folder,
	}, nil
}

func (db *Filestore) makeFileName(repository, key string) string {
	repo := base64.StdEncoding.EncodeToString([]byte(repository))
	return fmt.Sprintf("%s/out_%s_%s.secret", db.folder, repo, key)
}

// Put writes a file to the filesystem
func (db *Filestore) Put(repository, key string, secret []byte) error {
	file := db.makeFileName(repository, key)
	err := ioutil.WriteFile(file, secret, 0644)
	return err
}

// Get reads a file from the filesystem
func (db *Filestore) Get(repository, key string) (string, error) {
	file := db.makeFileName(repository, key)

	b, err := ioutil.ReadFile(file)
	return string(b), err

}

// Remove is used to delete a file from the filesystem
func (db *Filestore) Remove(repository, key string) error {
	file := db.makeFileName(repository, key)
	return os.Remove(file)
}

// Has is used to check the precence of a a key
func (db *Filestore) Has(repository string, key string) bool {
	file := db.makeFileName(repository, key)
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}

// List returns a list of keys stored for the given repo
func (db *Filestore) List(repository string) []string {
	var secrets []string
	files, err := ioutil.ReadDir(db.folder)
	if err != nil {
		return []string{}
	}
	repo := base64.StdEncoding.EncodeToString([]byte(repository))
	for _, f := range files {
		filePrefix := fmt.Sprintf("out_%s_", repo)
		if strings.HasPrefix(f.Name(), filePrefix) {
			secret := strings.Replace(f.Name(), filePrefix, "", 1)
			secret = strings.Replace(secret, ".secret", "", 1)
			secrets = append(secrets, secret)
		}
	}
	return secrets
}
