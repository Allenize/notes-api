package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type claims struct {
	Sub  string `json:"sub"`
	Type string `json:"type"`
	Iat  int64  `json:"iat"`
	Exp  int64  `json:"exp"`
}

func b64(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func b64Decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

func sign(data string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(data))
	return b64(mac.Sum(nil))
}

func IssueToken(userID, tokenType string, secret []byte, ttl time.Duration) (string, error) {
	header := b64([]byte(`{"alg":"HS256","typ":"JWT"}`))

	now := time.Now()
	c := claims{
		Sub:  userID,
		Type: tokenType,
		Iat:  now.Unix(),
		Exp:  now.Add(ttl).Unix(),
	}
	payloadBytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	payload := b64(payloadBytes)

	unsigned := header + "." + payload
	signature := sign(unsigned, secret)

	return unsigned + "." + signature, nil
}

func VerifyToken(token string, secret []byte) (userID string, tokenType string, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", "", ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSig := sign(unsigned, secret)
	if subtle.ConstantTimeCompare([]byte(expectedSig), []byte(parts[2])) != 1 {
		return "", "", ErrInvalidToken
	}

	payloadBytes, err := b64Decode(parts[1])
	if err != nil {
		return "", "", ErrInvalidToken
	}

	var c claims
	if err := json.Unmarshal(payloadBytes, &c); err != nil {
		return "", "", ErrInvalidToken
	}

	if time.Now().Unix() > c.Exp {
		return "", "", ErrInvalidToken
	}

	return c.Sub, c.Type, nil
}
