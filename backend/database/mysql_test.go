package database

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/alde/horus/backend/config"
)

func setupMySQLTestContainer() *config.Config {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mysql:latest",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "horus_password",
			"MYSQL_DATABASE":      "horus_testing",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL"),
	}
	mysqlC, _ := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	cfg := config.DefaultConfig()
	cfg.MySQL.Host, _ = mysqlC.Host(ctx)
	p, _ := mysqlC.MappedPort(ctx, "3306/tcp")
	cfg.MySQL.Port = p.Int()
	cfg.MySQL.Database = "horus_testing"
	cfg.MySQL.Username = "root"
	tmpdir := os.TempDir()
	file, _ := ioutil.TempFile(tmpdir, "sql")
	_, _ = file.WriteString("horus_password")
	cfg.MySQL.PasswordFile = file.Name()

	return cfg
}

func Test_MySQL_FullCycle(t *testing.T) {
	cfg := setupMySQLTestContainer()
	db, err := NewMySQL(cfg)

	if err != nil {
		t.Logf("failed setting up database: %+v", err)
		t.Fail()
		return
	}

	err = db.Put("github.com/alde/horus", "DOCKER_LOGIN", []byte("totally_secret"))
	assert.Nil(t, err, "error should be nil")
	err = db.Put("github.com/alde/horus", "DOCKER_PASSWORD", []byte("a different secret"))
	assert.Nil(t, err, "error should be nil")

	secret, err := db.Get("github.com/alde/horus", "DOCKER_LOGIN")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("totally_secret")), secret, "the stored and retrieved secret should match")

	secrets := db.List("github.com/alde/horus")
	assert.Len(t, secrets, 2, "should have 2 secrets")
	assert.Equal(t, secrets, []string{"DOCKER_LOGIN", "DOCKER_PASSWORD"})

	err = db.Remove("github.com/alde/horus", "DOCKER_LOGIN")
	assert.Nil(t, err, "error should be nil")

	has := db.Has("github.com/alde/horus", "DOCKER_LOGIN")
	assert.False(t, has, "should no longer have DOCKER_LOGIN")
}
