package jwks

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(kid string, memberId string, accessSecret string, accessExpire int64, extraClaims jwt.MapClaims) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(time.Duration(accessExpire) * time.Second)
	claims := map[string]interface{}{
		"exp": exp.Unix(),
		"sub": memberId,
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"iss": "api-hub",
		"aud": "member",
	}
	maps.Copy(claims, extraClaims)
	tokenOption := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims(claims))
	tokenOption.Header[jwkset.HeaderKID] = kid
	token, err := tokenOption.SignedString(ed25519.PrivateKey(accessSecret))
	if err != nil {
		return "", time.Time{}, err
	}
	return token, exp, nil
}

func ValidateToken(tokenStr string, k keyfunc.Keyfunc) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, k.Keyfunc)
	if err != nil {
		return nil, err
	}
	claims[jwkset.HeaderKID] = token.Header[jwkset.HeaderKID]
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func MapClaimsParseString(claims map[string]interface{}, key string) (string, error) {
	var (
		ok  bool
		raw interface{}
		s   string
	)
	raw, ok = claims[key]
	if !ok {
		return "", nil
	}

	s, ok = raw.(string)
	if !ok {
		return "", fmt.Errorf("%s is invalid", key)
	}

	return s, nil
}
