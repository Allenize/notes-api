package auth

import (
	"testing"
	"time"
)

func TestIssueAndVerifyToken(t *testing.T) {
	secret := []byte("test-secret")

	token, err := IssueToken("user-123", "access", secret, time.Minute)
	if err != nil {
		t.Fatalf("IssueToken returned error: %v", err)
	}

	userID, tokenType, err := VerifyToken(token, secret)
	if err != nil {
		t.Fatalf("VerifyToken returned error: %v", err)
	}
	if userID != "user-123" {
		t.Fatalf("expected user-123, got %s", userID)
	}
	if tokenType != "access" {
		t.Fatalf("expected access, got %s", tokenType)
	}
}

func TestVerifyTokenRejectsExpired(t *testing.T) {
	secret := []byte("test-secret")

	token, _ := IssueToken("user-123", "access", secret, -time.Minute)

	_, _, err := VerifyToken(token, secret)
	if err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken for expired token, got %v", err)
	}
}

func TestVerifyTokenRejectsTamperedSignature(t *testing.T) {
	secret := []byte("test-secret")

	token, _ := IssueToken("user-123", "access", secret, time.Minute)
	tampered := token[:len(token)-4] + "abcd"

	_, _, err := VerifyToken(tampered, secret)
	if err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken for tampered token, got %v", err)
	}
}

func TestVerifyTokenRejectsWrongSecret(t *testing.T) {
	token, _ := IssueToken("user-123", "access", []byte("secret-a"), time.Minute)

	_, _, err := VerifyToken(token, []byte("secret-b"))
	if err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken for wrong secret, got %v", err)
	}
}
