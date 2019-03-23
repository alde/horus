package database

import (
	"context"

	"github.com/kardianos/osext"

	"github.com/sirupsen/logrus"

	"github.com/alde/horus/backend/config"
)

// Database interface providing the contract that we expect
type Database interface {
	Put(repository, key string, secret []byte) error
	Get(repository, key string) (string, error)
	Remove(repository, key string) error
	Has(repository string, key string) bool
	List(repository string) []string
}

// Init is used to initialize the database settings
func Init(ctx context.Context, cfg *config.Config) (Database, error) {
	logrus.Info("setting up filesystem pretend database")
	folder, _ := osext.ExecutableFolder()
	database, err := NewFilestore(folder)
	if err != nil {
		logrus.WithError(err).Error("unable to create fake database")
	}
	return database, err
}
