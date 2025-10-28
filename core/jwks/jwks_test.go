package jwks_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"os"
	"testing"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
)

func TestNewKeyfunc(t *testing.T) {
	k, generateTokenFunc, err := jwks.NewKeyfunc(context.Background())
	if err != nil {
		t.Fatalf("Failed to init keyfunc.\nError: %s", err)
	}
	token, _, err := generateTokenFunc("1234567890", 60, nil)
	if err != nil {
		t.Fatalf("Failed to generate token.\nError: %s", err)
	}
	t.Logf("token: %s", token)
	claims, err := jwks.ValidateToken(token, k)
	if err != nil {
		t.Fatalf("Failed to validate token.\nError: %s", err)
	}
	t.Logf("claims: %+v", claims)
	memberId, err := claims.GetSubject()
	if err != nil {
		t.Fatalf("Failed to get member id.\nError: %s", err)
	}
	assert.Equal(t, memberId, "1234567890")

}

func TestGenerateOTP(t *testing.T) {
	otp, qrCode, err := jwks.GenerateOTP("api-hub", "admin")
	if err != nil {
		t.Fatalf("Failed to generate OTP.\nError: %s", err)
	}
	t.Logf("otp: %s", otp)
	os.WriteFile("../../test/data/qrcode.png", qrCode, 0644)
}

func TestValidateOTP(t *testing.T) {
	ok := totp.Validate("370550", "II5UPLT5LHFTKEIY2Q4NP2VWWHPEWOOV")
	t.Logf("ok: %v", ok)
}

func TestJwt(t *testing.T) {
	pub2, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate given key.\nError: %s", err)
	}
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate given key.\nError: %s", err)
	}
	kid := "2"
	token, _, err := jwks.GenerateToken(kid, "1234567890", string(priv), 60, nil)
	if err != nil {
		t.Fatalf("Failed to generate token.\nError: %s", err)
	}
	t.Logf("token: %s", token)

	jwksets := jwkset.NewMemoryStorage()
	jwks.AddKey(context.Background(), "1", pub, jwksets)
	jwks.AddKey(context.Background(), "2", pub2, jwksets)
	rawJWKS, err := jwksets.JSONPublic(context.Background())
	if err != nil {
		t.Fatalf("Failed to get JWKS.\nError: %s", err)
	}
	t.Logf("rawJWKS: %s", string(rawJWKS))

	k, err := keyfunc.New(keyfunc.Options{
		Ctx:     context.Background(),
		Storage: jwksets,
	})
	if err != nil {
		t.Fatalf("Failed to init keyfunc.\nError: %s", err)
	}
	claims, err := jwks.ValidateToken(token, k)
	if err != nil {
		t.Fatalf("Failed to validate token.\nError: %s", err)
	}
	t.Logf("claims: %+v", claims)
}

func TestRotateKey(t *testing.T) {
	jwksets := jwkset.NewMemoryStorage()
	pub, priv, err := jwks.RotateKey(context.Background(), "api-hub", jwksets)
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

	if err := jwks.AddKey(context.Background(), "api-hub", pub, jwksets); err != nil {
		t.Fatalf("Failed to add key.\nError: %s", err)
	}

	rawJWKS, err = jwksets.JSONPublic(context.Background())
	if err != nil {
		t.Fatalf("Failed to get JWKS.\nError: %s", err)
	}
	t.Logf("rawJWKS: %s", string(rawJWKS))
}
