package encryptor

import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"

	"github.com/alde/horus/backend/config"
)

// CloudKMS is a Google Cloud Key Management Service implementation of the Encryptor interface
type CloudKMS struct {
	ctx    context.Context
	config *config.Config
}

// NewGoogleCloudKMS sets up a new CloudKMS instance
func NewGoogleCloudKMS(ctx context.Context, cfg *config.Config) Encryptor {
	return &CloudKMS{
		ctx:    ctx,
		config: cfg,
	}
}

// Encrypt a byte-array using Google Cloud KMS
func (c *CloudKMS) Encrypt(bytes []byte) ([]byte, error) {
	client, err := cloudkms.NewKeyManagementClient(c.ctx)
	if err != nil {
		return nil, err
	}

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

	resp, err := client.Encrypt(c.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Ciphertext, nil
}
