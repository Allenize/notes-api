package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
	pbkdf2Iterations = 100000
	saltBytes        = 16
	keyBytes         = 32
)

func pbkdf2(password, salt []byte, iterations, keyLen int) []byte {
	prf := func(key, data []byte) []byte {
		mac := hmac.New(sha256.New, key)
		mac.Write(data)
		return mac.Sum(nil)
	}

	hashLen := sha256.Size
	numBlocks := (keyLen + hashLen - 1) / hashLen
	derived := make([]byte, 0, numBlocks*hashLen)

	for block := 1; block <= numBlocks; block++ {
		blockIndex := []byte{
			byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block),
		}
		u := prf(password, append(append([]byte{}, salt...), blockIndex...))
		t := append([]byte{}, u...)
		for i := 1; i < iterations; i++ {
			u = prf(password, u)
			for j := range t {
				t[j] ^= u[j]
			}
		}
		derived = append(derived, t...)
	}

	return derived[:keyLen]
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := pbkdf2([]byte(password), salt, pbkdf2Iterations, keyBytes)
	return fmt.Sprintf("%d$%s$%s", pbkdf2Iterations, hex.EncodeToString(salt), hex.EncodeToString(hash)), nil
}

func VerifyPassword(password, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 3 {
		return false
	}
	iterations, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	salt, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	expected, err := hex.DecodeString(parts[2])
	if err != nil {
		return false
	}
	actual := pbkdf2([]byte(password), salt, iterations, len(expected))
	return subtle.ConstantTimeCompare(actual, expected) == 1
}
