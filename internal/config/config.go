package config

import (
"errors"
"os"
"strconv"

"github.com/joho/godotenv"
)

type Config struct {
Addr      string
JWTSecret string
TokenTTL  int
DBDSN     string
LogLevel  string
}

func Load() (*Config, error) {
// Load .env file if it exists (non-blocking if not found)
_ = godotenv.Load()

addr := os.Getenv("SSO_ADDR")
if addr == "" {
addr = ":8080"
}

secret := os.Getenv("SSO_JWT_SECRET")
if secret == "" {
return nil, errors.New("SSO_JWT_SECRET is required")
}

dbdsn := os.Getenv("SSO_DB_DSN")
if dbdsn == "" {
dbdsn = "postgres://postgres:postgres@postgres:5432/sso?sslmode=disable"
}

ttlStr := os.Getenv("SSO_TOKEN_TTL")
ttl := 3600 // default 1 hour
if ttlStr != "" {
if parsed, err := strconv.Atoi(ttlStr); err == nil {
ttl = parsed
}
}

logLevel := os.Getenv("SSO_LOG_LEVEL")
if logLevel == "" {
logLevel = "info"
}

return &Config{
Addr:      addr,
JWTSecret: secret,
TokenTTL:  ttl,
DBDSN:     dbdsn,
LogLevel:  logLevel,
}, nil
}
