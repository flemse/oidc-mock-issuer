package issuer

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/flemse/oidc-mock-issuer/internal/tokens"
)

// TokenHandler issues tokens based on request parameters.
func TokenHandler(ks *tokens.KeySet, issuerURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Valid  bool                   `json:"valid"`
			Claims map[string]interface{} `json:"claims"`
		}
		req.Valid = true // default to valid
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.Claims == nil {
			req.Claims = make(map[string]interface{})
		}
		if _, ok := req.Claims["iss"]; !ok {
			req.Claims["iss"] = issuerURL
		}
		tok, err := ks.GenerateToken(req.Valid, req.Claims)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"access_token": tok, "token_type": "Bearer"}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
