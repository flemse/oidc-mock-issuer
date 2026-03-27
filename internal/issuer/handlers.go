package issuer

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/flemse/oidc-mock-issuer/internal/tokens"
)

// DiscoveryHandler serves the OIDC discovery document.
func DiscoveryHandler(issuerURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"issuer":         issuerURL,
			"jwks_uri":       issuerURL + "/.well-known/jwks.json",
			"token_endpoint": issuerURL + "/token",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// JWKSHandler serves the public JWKS derived from the given KeySet.
func JWKSHandler(ks *tokens.KeySet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pub := &ks.PrivateKey.PublicKey
		key := map[string]interface{}{
			"kty": "RSA",
			"use": "sig",
			"alg": "RS256",
			"kid": ks.KID,
			"n":   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
			"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
		}
		jwks := map[string]interface{}{"keys": []interface{}{key}}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jwks); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
