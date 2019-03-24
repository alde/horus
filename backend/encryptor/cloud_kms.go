package encryptor

import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"

	"github.com/alde/horus/backend/config"
)

// CloudKMS is a Google Cloud Key Management Service implementation of the Encryptor interface
type CloudKMS struct {
	ctx    context.Context
	config *config.Config
	client kmsEncryptor
}

type kmsEncryptor interface {
	Encrypt(context.Context, *kmspb.EncryptRequest, ...gax.CallOption) (*kmspb.EncryptResponse, error)
}

// NewGoogleCloudKMS sets up a new CloudKMS instance
func NewGoogleCloudKMS(ctx context.Context, cfg *config.Config) (Encryptor, error) {
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	return &CloudKMS{
		ctx:    ctx,
		config: cfg,
		client: client,
	}, nil
}

// Encrypt a byte-array using Google Cloud KMS
func (c *CloudKMS) Encrypt(bytes []byte) ([]byte, error) {
	keyname := fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		c.config.GoogleKMS.Project,
		c.config.GoogleKMS.Location,
		c.config.GoogleKMS.KeyRing,
		c.config.GoogleKMS.KeyName)

	req := &kmspb.EncryptRequest{
		Name:      keyname,
		Plaintext: bytes,
	}

	resp, err := c.client.Encrypt(c.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Ciphertext, nil
}
