package jwks_test

import (
	"context"
	"testing"

	"github.com/MicahParks/jwkset"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"
)

func TestRotateKey(t *testing.T) {
	jwksets := jwkset.NewMemoryStorage()
	pub, priv, err := jwks.RotateKey(context.Background(), jwksets)
	if err != nil {
		t.Fatalf("Failed to rotate key.\nError: %s", err)
	}
	t.Logf("pub: %s", string(pub))
	t.Logf("priv: %s", string(priv))

	rawJWKS, err := jwksets.JSONPublic(context.Background())
	if err != nil {
		t.Fatalf("Failed to get JWKS.\nError: %s", err)
	}
	t.Logf("rawJWKS: %s", string(rawJWKS))

	encryptPub, err := xstrings.AesCbcEncryptBase64(string(pub), "public_key_secre", nil)
	if err != nil {
		t.Fatalf("Failed to encrypt public key.\nError: %s", err)
	}
	encryptPriv, err := xstrings.AesCbcEncryptBase64(string(priv), "private_key_secr", nil)
	if err != nil {
		t.Fatalf("Failed to encrypt private key.\nError: %s", err)
	}
	t.Logf("encryptPub: %s", encryptPub)
	t.Logf("encryptPriv: %s", encryptPriv)

	decryptPub, err := xstrings.AesCbcDecryptBase64(encryptPub, "public_key_secre", nil)
	if err != nil {
		t.Fatalf("Failed to decrypt public key.\nError: %s", err)
	}
	decryptPriv, err := xstrings.AesCbcDecryptBase64(encryptPriv, "private_key_secr", nil)
	if err != nil {
		t.Fatalf("Failed to decrypt private key.\nError: %s", err)
	}
	t.Logf("decryptPub: %s", decryptPub)
	t.Logf("decryptPriv: %s", decryptPriv)

	if err := jwks.AddKey(context.Background(), pub, jwksets); err != nil {
		t.Fatalf("Failed to add key.\nError: %s", err)
	}

	rawJWKS, err = jwksets.JSONPublic(context.Background())
	if err != nil {
		t.Fatalf("Failed to get JWKS.\nError: %s", err)
	}
	t.Logf("rawJWKS: %s", string(rawJWKS))
}
