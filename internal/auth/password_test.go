package auth

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if !VerifyPassword("correct-horse-battery-staple", hash) {
		t.Fatal("expected correct password to verify successfully")
	}

	if VerifyPassword("wrong-password", hash) {
		t.Fatal("expected wrong password to fail verification")
	}
}

func TestHashPasswordProducesUniqueSalts(t *testing.T) {
	hash1, _ := HashPassword("same-password")
	hash2, _ := HashPassword("same-password")

	if hash1 == hash2 {
		t.Fatal("expected two hashes of the same password to differ due to random salts")
	}
}

func TestVerifyPasswordRejectsMalformedHash(t *testing.T) {
	if VerifyPassword("anything", "not-a-valid-hash") {
		t.Fatal("expected malformed hash to fail verification")
	}
}
