package tokens

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// KeySet holds an RSA key pair used to sign and verify tokens.
type KeySet struct {
	PrivateKey *rsa.PrivateKey
	KID        string
}

// NewKeySet generates a new RSA key pair.
func NewKeySet() (*KeySet, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &KeySet{PrivateKey: priv, KID: "key-1"}, nil
}

// GenerateToken issues a signed JWT. If valid is false, an unsigned invalid token string is returned.
func (k *KeySet) GenerateToken(valid bool, claims map[string]interface{}) (string, error) {
	if !valid {
		return "invalid.token.here", nil
	}
	jwtClaims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	for key, val := range claims {
		jwtClaims[key] = val
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)
	token.Header["kid"] = k.KID
	return token.SignedString(k.PrivateKey)
}
