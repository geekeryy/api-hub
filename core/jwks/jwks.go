package jwks

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"log"

	"github.com/MicahParks/jwkset"
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

func RotateKey(ctx context.Context,kid string, jwksets *jwkset.MemoryJWKSet) (ed25519.PublicKey, ed25519.PrivateKey, error) {
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
