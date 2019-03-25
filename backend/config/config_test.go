package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DefaultConfig(t *testing.T) {
	c := DefaultConfig()
	a := assert.New(t)
	a.Equal("0.0.0.0", c.Server.Address)
	a.Equal(7654, c.Server.Port)
	a.Equal("text", c.Logging.Format)
	a.Equal("DEBUG", c.Logging.Level)
	a.Equal(os.Getenv("USER"), c.Metadata["owner"])

}

func Test_ReadConfigFile(t *testing.T) {
	c := DefaultConfig()
	wd, _ := os.Getwd()
	a := assert.New(t)

	ReadConfigFile(c, fmt.Sprintf("%s/config_test.toml", wd))
	a.Equal("127.0.0.1", c.Server.Address)
	a.Equal(8080, c.Server.Port)

	a.Equal("json", c.Logging.Format)
	a.Equal("INFO", c.Logging.Level)

	a.Equal("a-gcp-project", c.GoogleKMS.Project)
	a.Equal("global", c.GoogleKMS.Location)
	a.Equal("horus-keyring", c.GoogleKMS.KeyRing)
	a.Equal("horus-key", c.GoogleKMS.KeyName)

	a.Equal("alde", c.Metadata["owner"])

	a.Equal("127.0.0.1", c.MySQL.Host)
	a.Equal(3306, c.MySQL.Port)
	a.Equal("horus_user", c.MySQL.Username)
	a.Equal("/etc/credentials/database_password", c.MySQL.PasswordFile)
	a.Equal("horus_database", c.MySQL.Database)
}

func Test_ReadConfigFile_Error(t *testing.T) {
	c := DefaultConfig()
	d := DefaultConfig()

	ReadConfigFile(c, getConfigFilePath(""))

	assert.Equal(t, c, d)
}

func Test_getConfigFilePath(t *testing.T) {
	fp := getConfigFilePath("")
	assert.Empty(t, fp)
}

func Test_Initialize(t *testing.T) {
	c := Initialize("")
	d := DefaultConfig()

	assert.Equal(t, c, d)
}
