package config

import (
	"os"
	"regexp"
	"strings"
)

var (
	DBFile = getEnv("TRAVELPORTAL_DBFILE", "TODO_DBFILE", "scheduler.db")
	Port   = getEnv("TRAVELPORTAL_PORT", "TODO_PORT", "7540")
)

var vaultPassword = regexp.MustCompile(`(travelportalPassword|todoPassword)\s*[:=]\s*([^\]\s]+)`)

// GetPassword reads the configured password, then checks the Vault-injected file.
func GetPassword() string {
	if v := getEnv("TRAVELPORTAL_PASSWORD", "TODO_PASSWORD", ""); v != "" {
		return v
	}
	b, err := os.ReadFile("/vault/secrets/config.txt")
	if err == nil {
		s := strings.TrimSpace(string(b))
		if s != "" {
			if m := vaultPassword.FindStringSubmatch(s); m != nil {
				return m[2]
			}
			// Fallback: if file is just the password
			return s
		}
	}
	return ""
}

func getEnv(key, legacyKey, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if value, exists := os.LookupEnv(legacyKey); exists {
		return value
	}
	return defaultValue
}
