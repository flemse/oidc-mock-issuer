package tokens

import (
	"strings"
	"testing"
)

func TestNewKeySet(t *testing.T) {
	ks, err := NewKeySet()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ks.PrivateKey == nil {
		t.Error("expected non-nil private key")
	}
	if ks.KID == "" {
		t.Error("expected non-empty KID")
	}
}

func TestGenerateToken_Valid(t *testing.T) {
	ks, err := NewKeySet()
	if err != nil {
		t.Fatalf("failed to create key set: %v", err)
	}
	tok, err := ks.GenerateToken(true, map[string]interface{}{"sub": "user1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		t.Errorf("expected JWT with 3 parts, got %d", len(parts))
	}
}

func TestGenerateToken_Invalid(t *testing.T) {
	ks, err := NewKeySet()
	if err != nil {
		t.Fatalf("failed to create key set: %v", err)
	}
	tok, err := ks.GenerateToken(false, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "invalid.token.here" {
		t.Errorf("expected invalid token placeholder, got %q", tok)
	}
}
