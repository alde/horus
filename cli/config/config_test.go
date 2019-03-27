package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	wd, _ := os.Getwd()
	c := New(fmt.Sprintf("%s/config_test.toml", wd))

	assert.Equal(t, c.Horus.Host, "http://horus.svc:7654/")
}
