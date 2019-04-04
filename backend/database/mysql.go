package database

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// Import mysql into the scope of this package (required)
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/alde/horus/backend/config"
)

// MySQL struct implementing the Database interface with MySQL backend
type MySQL struct {
	db *sql.DB
}

func readPasswordFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}

	password, _ := ioutil.ReadAll(file)
	return strings.Trim(string(password), "")
}

// NewMySQL is used to create a new MySQL database connection
func NewMySQL(cfg *config.Config) (Database, error) {
	password := readPasswordFile(cfg.MySQL.PasswordFile)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.MySQL.Username, password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)

	if cfg.MySQL.DisableSSL {
		connectionString += "?tls=skip-verify"
	}

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		logrus.WithField("connectionString", connectionString).WithError(err).Error("failed to initialize driver")
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logrus.WithField("connectionString", connectionString).WithError(err).Error("failed to ping database")
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS horus_secrets ( \n" +
		" `repository` VARCHAR(128) NOT NULL, \n" +
		" `secret_key` VARCHAR(128) NOT NULL, \n" +
		" `secret_value` VARCHAR(128) NOT NULL, \n" +
		" PRIMARY KEY (`repository`, `secret_key`) \n" +
		")")
	if err != nil {
		logrus.WithError(err).Error("unable to create database table")
		return nil, err
	}

	return &MySQL{
		db: db,
	}, nil
}

// Put will insert a secret into the database, overwriting any existing entry
func (m MySQL) Put(repository, key string, secret []byte) error {
	b64 := base64.StdEncoding.EncodeToString(secret)
	query := "INSERT INTO horus_secrets (`repository`, `secret_key`, `secret_value`) " +
		"VALUES(?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE `secret_value` = ?"
	_, err := m.db.Exec(
		query,
		repository, key, b64, b64)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"repo": repository,
			"key":  key,
		}).WithError(err).Error("error inserting secret into database")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"repo": repository,
		"key":  key,
	}).Info("inserted secret into database")

	return nil
}

// Get a secret from the database
func (m MySQL) Get(repository, key string) (string, error) {
	query := "SELECT `secret_value` FROM horus_secrets WHERE `repository` = ? AND `secret_key` = ?"
	row := m.db.QueryRow(query, repository, key)
	var secret string
	err := row.Scan(&secret)
	return secret, err
}

// Remove a secret from the database
func (m MySQL) Remove(repository, key string) error {
	query := "DELETE FROM horus_secrets WHERE `repository` = ? AND `secret_key` = ?"
	_, err := m.db.Exec(query, repository, key)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"repo": repository,
			"key":  key,
		}).WithError(err).Error("error deleting secret from database")
	} else {
		logrus.WithFields(logrus.Fields{
			"repo": repository,
			"key":  key,
		}).Info("deleted secret from database")
	}
	return err
}

// Has checks the existance of a secret in the database
func (m MySQL) Has(repository string, key string) bool {
	query := "SELECT EXISTS(SELECT 1 FROM horus_secrets WHERE `repository` = ? AND `secret_key` = ?)"
	rows, err := m.db.Query(query, repository, key)
	if err != nil {
		return false
	}
	for rows.Next() {
		var count int
		_ = rows.Scan(&count)
		return count == 1
	}
	return false
}

// List all secrets for a repository
func (m MySQL) List(repository string) []string {
	var secrets = []string{}
	query := "SELECT `secret_key` FROM horus_secrets WHERE `repository` = ?"
	rows, err := m.db.Query(query, repository)
	if err != nil {
		logrus.WithField("repo", repository).WithError(err).Error("unable to list secrets")
		return secrets
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			logrus.WithError(err).Error("unable to extract column")
		}
		secrets = append(secrets, key)
	}
	return secrets
}
