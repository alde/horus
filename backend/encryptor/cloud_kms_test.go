package encryptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alde/horus/backend/config"
	"github.com/alde/horus/backend/mock"
)

func Test_Create(t *testing.T) {
	a := assert.New(t)
	kms, err := NewGoogleCloudKMS(context.Background(), &config.Config{})
	a.Nil(err)
	a.NotNil(kms)
}

func Test_Encrypt(t *testing.T) {
	a := assert.New(t)
	m := mock.Encryptor{}
	kms := &CloudKMS{
		ctx:    context.Background(),
		config: &config.Config{},
		client: m,
	}

	encrypted, e := kms.Encrypt([]byte("a-secret"))
	a.Nil(e)
	a.NotEqual("a-secret", string(encrypted))
	a.Equal("A-SECRET", string(encrypted))
}
