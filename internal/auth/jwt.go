package auth

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "fmt"
    "strconv"
    "strings"
    "time"
)

type JWTManager struct {
    Secret string
    TTL    time.Duration
}

func NewJWTManager(secret string, ttlSeconds int) *JWTManager {
    return &JWTManager{Secret: secret, TTL: time.Duration(ttlSeconds) * time.Second}
}

func (m *JWTManager) GenerateToken(userID string) (string, error) {
    expires := time.Now().Add(m.TTL).Unix()
    payload := fmt.Sprintf("%s|%d", userID, expires)
    payloadEnc := base64.RawURLEncoding.EncodeToString([]byte(payload))

    mac := hmac.New(sha256.New, []byte(m.Secret))
    _, _ = mac.Write([]byte(payloadEnc))
    sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

    return payloadEnc + "." + sig, nil
}

func (m *JWTManager) VerifyToken(tokenString string) (string, error) {
    parts := strings.Split(tokenString, ".")
    if len(parts) != 2 {
        return "", errors.New("invalid token format")
    }

    payloadEnc, sig := parts[0], parts[1]

    mac := hmac.New(sha256.New, []byte(m.Secret))
    _, _ = mac.Write([]byte(payloadEnc))
    expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

    if !hmac.Equal([]byte(sig), []byte(expected)) {
        return "", errors.New("invalid token signature")
    }

    raw, err := base64.RawURLEncoding.DecodeString(payloadEnc)
    if err != nil {
        return "", err
    }

    parts2 := strings.Split(string(raw), "|")
    if len(parts2) != 2 {
        return "", errors.New("invalid token payload")
    }

    userID := parts2[0]
    exp, err := strconv.ParseInt(parts2[1], 10, 64)
    if err != nil {
        return "", err
    }

    if time.Now().Unix() > exp {
        return "", errors.New("token expired")
    }

    return userID, nil
}
