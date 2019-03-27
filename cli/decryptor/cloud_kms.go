package decryptor

import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"

	"github.com/alde/horus/cli/config"
)

// CloudKMS is a Google Cloud Key Management Service implementation of the Decryptor interface
type CloudKMS struct {
	ctx    context.Context
	config *config.Config
	client kmsDecryptor
}

type kmsDecryptor interface {
	Decrypt(context.Context, *kmspb.DecryptRequest, ...gax.CallOption) (*kmspb.DecryptResponse, error)
}

// NewGoogleCloudKMS sets up a new CloudKMS instance
// If kmsClient is nil one will be created
func NewGoogleCloudKMS(ctx context.Context, cfg *config.Config, kmsClient kmsDecryptor) (Decryptor, error) {
	if kmsClient == nil {
		var err error
		kmsClient, err = cloudkms.NewKeyManagementClient(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &CloudKMS{
		ctx:    ctx,
		config: cfg,
		client: kmsClient,
	}, nil
}

// Decrypt a byte-array using Google Cloud KMS
func (c *CloudKMS) Decrypt(secret []byte) ([]byte, error) {
	keyname := fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		c.config.GoogleKMS.Project,
		c.config.GoogleKMS.Location,
		c.config.GoogleKMS.KeyRing,
		c.config.GoogleKMS.KeyName)

	req := &kmspb.DecryptRequest{
		Name:       keyname,
		Ciphertext: []byte(secret),
	}

	resp, err := c.client.Decrypt(c.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}
