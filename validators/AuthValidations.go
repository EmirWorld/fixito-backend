package validators

import "strings"

func ValidatePassword(password string) (string, bool) {
	if len(password) < 8 {
		return "Password must be at least 8 characters", false
	}
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return "Password must contain at least one uppercase letter", false
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return "Password must contain at least one lowercase letter", false
	}
	if !strings.ContainsAny(password, "0123456789") {
		return "Password must contain at least one number", false
	}
	return "", true
}
