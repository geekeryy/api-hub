package jwks

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"log"
	"log/slog"
	"maps"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

func AddKey(ctx context.Context, kid string, pub ed25519.PublicKey, jwksets *jwkset.MemoryJWKSet) error {
	// Turn the key into a JWK.
	marshalOptions := jwkset.JWKMarshalOptions{
		Private: true,
	}

	metadata := jwkset.JWKMetadataOptions{
		KID: kid,
	}
	options := jwkset.JWKOptions{
		Marshal:  marshalOptions,
		Metadata: metadata,
	}
	jwk, err := jwkset.NewJWKFromKey(pub, options)
	if err != nil {
		log.Fatalf("Failed to create a JWK from the given key.\nError: %s", err)
		return err
	}

	// Write the JWK to the server's storage.
	err = jwksets.KeyWrite(ctx, jwk)
	if err != nil {
		log.Fatalf("Failed to write the JWK to the server's storage.\nError: %s", err)
		return err
	}
	return nil
}

func RotateKey(ctx context.Context, kid string, jwksets *jwkset.MemoryJWKSet) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate given key.\nError: %s", err)
		return nil, nil, err
	}
	err = AddKey(ctx, kid, pub, jwksets)
	if err != nil {
		log.Fatalf("Failed to add the key to the server's storage.\nError: %s", err)
		return nil, nil, err
	}
	return pub, priv, nil
}

// NewKeyfunc 生成一个keyfunc.Keyfunc和一个GenerateTokenFunc
func NewKeyfunc(ctx context.Context) (keyfunc.Keyfunc, GenerateTokenFunc, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate given key.\nError: %s", err)
		return nil, nil, err
	}

	kid := uuid.New().String()
	jwk, err := jwkset.NewJWKFromKey(pub, jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{
			KID: kid,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create a JWK from the given key.\nError: %s", err)
		return nil, nil, err
	}

	jwksets := jwkset.NewMemoryStorage()
	err = jwksets.KeyWrite(ctx, jwk)
	if err != nil {
		log.Fatalf("Failed to write the JWK to the server's storage.\nError: %s", err)
		return nil, nil, err
	}

	k, err := keyfunc.New(keyfunc.Options{
		Storage: jwksets,
	})
	if err != nil {
		log.Fatalf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
		return nil, nil, err
	}
	return k, GetGenerateTokenFunc(kid, priv), nil
}

type GenerateTokenFunc func(memberUUID string, accessExpire int64, extraClaims jwt.MapClaims) (token string, exp time.Time, err error)

func GetGenerateTokenFunc(kid string, priv ed25519.PrivateKey) func(string, int64, jwt.MapClaims) (string, time.Time, error) {
	return func(memberUUID string, accessExpire int64, extraClaims jwt.MapClaims) (string, time.Time, error) {
		now := time.Now()
		exp := now.Add(time.Duration(accessExpire) * time.Second)
		claims := map[string]interface{}{
			"exp": exp.Unix(),
			"sub": memberUUID,
			"iat": now.Unix(),
			"nbf": now.Unix(),
		}
		maps.Copy(claims, extraClaims)
		tokenOption := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims(claims))
		tokenOption.Header[jwkset.HeaderKID] = kid
		token, err := tokenOption.SignedString(priv)
		if err != nil {
			return "", time.Time{}, err
		}
		return token, exp, nil
	}
}

func NewDefaultOverrideCtx(ctx context.Context, requestJWKSetFunc func(ctx context.Context) (jwkset.JWKSMarshal, error), override keyfunc.Override) (keyfunc.Keyfunc, error) {
	rateLimitWaitMax := time.Minute
	if override.RateLimitWaitMax != 0 {
		rateLimitWaitMax = override.RateLimitWaitMax
	}
	refreshErrorHandler := func(u string) func(ctx context.Context, err error) {
		return func(ctx context.Context, err error) {
			slog.Default().ErrorContext(ctx, "Failed to refresh JWK Set from resource.", "error", err)
		}
	}
	if override.RefreshErrorHandlerFunc != nil {
		refreshErrorHandler = override.RefreshErrorHandlerFunc
	}
	refreshInterval := time.Hour
	if override.RefreshInterval > 0 {
		refreshInterval = override.RefreshInterval
	}
	refreshUnknownKID := rate.NewLimiter(rate.Every(5*time.Minute), 1)
	if override.RefreshUnknownKID != nil {
		refreshUnknownKID = override.RefreshUnknownKID
	}

	storage, err := jwkset.NewCustomStorage(jwkset.CustomStorageOptions{
		Ctx:                   ctx,
		NoErrorReturnFirstReq: true,
		RefreshErrorHandler:   refreshErrorHandler(""),
		RefreshInterval:       refreshInterval,
		ValidateOptions: jwkset.JWKValidateOptions{
			SkipAll: override.ValidationSkipAll,
		},
		RequestJWKSetFunc: requestJWKSetFunc,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create custom storage: %w", err)
	}

	storageOverride, err := jwkset.NewHTTPClient(jwkset.HTTPClientOptions{
		HTTPURLs: map[string]jwkset.Storage{
			"custom_storage": storage,
		},
		RateLimitWaitMax:  rateLimitWaitMax,
		RefreshUnknownKID: refreshUnknownKID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client storage: %w", err)
	}

	options := keyfunc.Options{
		Ctx:          ctx,
		Storage:      storageOverride,
		UseWhitelist: nil,
	}

	return keyfunc.New(options)
}
