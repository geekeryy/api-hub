package jwks_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/geekeryy/api-hub/core/jwks"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/rpc/auth/client/authservice"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"golang.org/x/time/rate"
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
	os.WriteFile("../../test/qrcode.png", qrCode, 0644)

	// assert.True(t, totp.Validate("927432", otp))
	// assert.False(t, totp.Validate("976604", otp))
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

func TestGrpcJwks(t *testing.T) {
	logger := logx.WithContext(context.Background())
	authClient := zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{"localhost:8880"},
	})
	kfunc, err := jwks.NewDefaultOverrideCtx(context.Background(), requestJWKSetFromGrpc(authClient), keyfunc.Override{
		RefreshInterval:   time.Duration(100) * time.Second,
		RateLimitWaitMax:  3 * time.Second,
		RefreshUnknownKID: rate.NewLimiter(rate.Every(1*time.Minute), 2),
		RefreshErrorHandlerFunc: func(u string) func(ctx context.Context, err error) {
			return func(ctx context.Context, err error) {
				logger.Errorf("Failed to refresh JWK Set from resource. Error: %v", err)
			}
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	token := `eyJhbGciOiJFZERTQSIsImtpZCI6IjQ0MjBiYTU5LWY4ZjItNDQyMC1iMjg5LTNjYzRlMjY1ZmY0NCIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJtZW1iZXIiLCJleHAiOjE3NjI4ODAzMDgsImlhdCI6MTc2Mjg3NjcwOCwiaXNzIjoiYXBpLWh1YiIsIm5iZiI6MTc2Mjg3NjcwOCwic3ViIjoiYTJiYWE4MzgtNzc1OC00YTAyLWEyOTktOWZjZTBhNTczOGJkIn0.Vekm17WjnUaw2Pi0A83rAEpq5bCXJTHsGVBrEHof8CI4XSfktvCyGhU1kMnSoVHmvp1lKXmWl1SCYdSnQLF7Dg`

	claims, err := jwks.ValidateToken(token, kfunc)
	if err != nil {
		t.Fatalf("Failed to validate token.\nError: %s", err)
	}
	t.Logf("claims: %+v", claims)
}
func requestJWKSetFromGrpc(conn zrpc.Client) func(ctx context.Context) (jwkset.JWKSMarshal, error) {
	return func(ctx context.Context) (jwkset.JWKSMarshal, error) {
		response, err := authservice.NewAuthService(conn).GetJwks(ctx, &authservice.GetJwksReq{})
		if err != nil {
			return jwkset.JWKSMarshal{}, err
		}
		var jwks jwkset.JWKSMarshal
		err = json.Unmarshal([]byte(response.Data), &jwks)
		if err != nil {
			return jwkset.JWKSMarshal{}, err
		}
		fmt.Printf("jwks: %+v \n", jwks)
		return jwks, nil
	}
}
