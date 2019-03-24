package mock

import (
	"context"
	"strings"

	gax "github.com/googleapis/gax-go/v2"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// Encryptor is a mock of the Google Cloud KMS Encryption
type Encryptor struct {
	EncryptFn        func(context.Context, *kmspb.EncryptRequest, ...gax.CallOption) (*kmspb.EncryptResponse, error)
	EncryptFnInvoked bool
}

func (e Encryptor) defaultEncryptFn(ctx context.Context, request *kmspb.EncryptRequest, opts ...gax.CallOption) (*kmspb.EncryptResponse, error) {
	return &kmspb.EncryptResponse{
		Ciphertext: []byte(strings.ToUpper(string(request.Plaintext))),
	}, nil
}

// Encrypt mock
func (e Encryptor) Encrypt(ctx context.Context, request *kmspb.EncryptRequest, opts ...gax.CallOption) (*kmspb.EncryptResponse, error) {
	e.EncryptFnInvoked = true
	if e.EncryptFn == nil {
		return e.defaultEncryptFn(ctx, request, opts...)
	}
	return e.EncryptFn(ctx, request, opts...)
}
