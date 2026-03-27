package issuer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flemse/oidc-mock-issuer/internal/tokens"
	"github.com/golang-jwt/jwt/v5"
)

func TestDiscoveryHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/.well-known/openid-configuration", nil)
	rw := httptest.NewRecorder()
	h := DiscoveryHandler("http://localhost:8080")
	h(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rw.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["issuer"] != "http://localhost:8080" {
		t.Errorf("unexpected issuer: %v", resp["issuer"])
	}
	if resp["jwks_uri"] != "http://localhost:8080/.well-known/jwks.json" {
		t.Errorf("unexpected jwks_uri: %v", resp["jwks_uri"])
	}
}

func TestJWKSHandler(t *testing.T) {
	ks, err := tokens.NewKeySet()
	if err != nil {
		t.Fatalf("failed to create key set: %v", err)
	}
	req := httptest.NewRequest("GET", "/.well-known/jwks.json", nil)
	rw := httptest.NewRecorder()
	JWKSHandler(ks)(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rw.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	keys, ok := resp["keys"].([]interface{})
	if !ok || len(keys) != 1 {
		t.Fatalf("expected 1 key, got %v", resp["keys"])
	}
}

func TestTokenHandler_Valid(t *testing.T) {
	ks, err := tokens.NewKeySet()
	if err != nil {
		t.Fatalf("failed to create key set: %v", err)
	}
	body := `{"valid":true,"claims":{"sub":"user1"}}`
	req := httptest.NewRequest("POST", "/token", strings.NewReader(body))
	rw := httptest.NewRecorder()
	TokenHandler(ks, "http://localhost:8080")(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rw.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["access_token"] == "" {
		t.Error("expected non-empty access_token")
	}
	// verify iss claim is set automatically
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	claims := jwt.MapClaims{}
	if _, _, err := parser.ParseUnverified(resp["access_token"], claims); err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if claims["iss"] != "http://localhost:8080" {
		t.Errorf("expected iss=http://localhost:8080, got %v", claims["iss"])
	}
}

func TestTokenHandler_Invalid(t *testing.T) {
	ks, err := tokens.NewKeySet()
	if err != nil {
		t.Fatalf("failed to create key set: %v", err)
	}
	body := `{"valid":false,"claims":{}}`
	req := httptest.NewRequest("POST", "/token", strings.NewReader(body))
	rw := httptest.NewRecorder()
	TokenHandler(ks, "http://localhost:8080")(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rw.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["access_token"] != "invalid.token.here" {
		t.Errorf("expected invalid token, got %q", resp["access_token"])
	}
}

func TestTokenHandler_NoBody(t *testing.T) {
	ks, err := tokens.NewKeySet()
	if err != nil {
		t.Fatalf("failed to create key set: %v", err)
	}
	req := httptest.NewRequest("POST", "/token", nil)
	rw := httptest.NewRecorder()
	TokenHandler(ks, "http://localhost:8080")(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rw.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["access_token"] == "" {
		t.Error("expected non-empty access_token")
	}
}
