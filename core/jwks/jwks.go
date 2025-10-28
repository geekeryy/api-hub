package jwks

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"log"
	"maps"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

type GenerateTokenFunc func(string, int64, jwt.MapClaims) (string, time.Time, error)

func GetGenerateTokenFunc(kid string, priv ed25519.PrivateKey) func(string, int64, jwt.MapClaims) (string, time.Time, error) {
	return func(memberId string, accessExpire int64, extraClaims jwt.MapClaims) (string, time.Time, error) {
		now := time.Now()
		exp := now.Add(time.Duration(accessExpire) * time.Second)
		claims := map[string]interface{}{
			"exp": exp.Unix(),
			"sub": memberId,
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
