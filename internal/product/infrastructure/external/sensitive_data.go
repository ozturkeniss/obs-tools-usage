package external

import (
	"regexp"
	"strings"
)

// SensitiveDataMasker handles masking of sensitive data in logs
type SensitiveDataMasker struct {
	patterns map[string]*regexp.Regexp
}

// NewSensitiveDataMasker creates a new sensitive data masker
func NewSensitiveDataMasker() *SensitiveDataMasker {
	return &SensitiveDataMasker{
		patterns: map[string]*regexp.Regexp{
			"password":     regexp.MustCompile(`(?i)(password|pass|pwd)\s*[:=]\s*["']?([^"'\s]+)["']?`),
			"api_key":      regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[:=]\s*["']?([^"'\s]+)["']?`),
			"token":        regexp.MustCompile(`(?i)(token|access[_-]?token|bearer)\s*[:=]\s*["']?([^"'\s]+)["']?`),
			"secret":       regexp.MustCompile(`(?i)(secret|secret[_-]?key)\s*[:=]\s*["']?([^"'\s]+)["']?`),
			"credit_card":  regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`),
			"ssn":          regexp.MustCompile(`\b\d{3}-?\d{2}-?\d{4}\b`),
			"email":        regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
			"phone":        regexp.MustCompile(`\b\d{3}[-.\s]?\d{3}[-.\s]?\d{4}\b`),
			"ip_address":   regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`),
			"jwt":          regexp.MustCompile(`eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*`),
		},
	}
}

// MaskSensitiveData masks sensitive data in a string
func (m *SensitiveDataMasker) MaskSensitiveData(data string) string {
	masked := data
	
	for patternName, pattern := range m.patterns {
		masked = m.maskPattern(masked, pattern, patternName)
	}
	
	return masked
}

// maskPattern masks a specific pattern
func (m *SensitiveDataMasker) maskPattern(data string, pattern *regexp.Regexp, patternName string) string {
	return pattern.ReplaceAllStringFunc(data, func(match string) string {
		return m.getMaskedValue(match, patternName)
	})
}

// getMaskedValue returns the appropriate masked value based on pattern type
func (m *SensitiveDataMasker) getMaskedValue(original, patternName string) string {
	switch patternName {
	case "password", "api_key", "token", "secret":
		// For key-value pairs, mask the value part
		parts := strings.SplitN(original, ":", 2)
		if len(parts) == 2 {
			return parts[0] + ": [MASKED]"
		}
		parts = strings.SplitN(original, "=", 2)
		if len(parts) == 2 {
			return parts[0] + "= [MASKED]"
		}
		return "[MASKED]"
	case "credit_card":
		// Show only last 4 digits
		digits := regexp.MustCompile(`\d`).FindAllString(original, -1)
		if len(digits) >= 4 {
			return "****-****-****-" + strings.Join(digits[len(digits)-4:], "")
		}
		return "[MASKED]"
	case "ssn":
		// Show only last 4 digits
		digits := regexp.MustCompile(`\d`).FindAllString(original, -1)
		if len(digits) >= 4 {
			return "***-**-" + strings.Join(digits[len(digits)-4:], "")
		}
		return "[MASKED]"
	case "email":
		// Mask the local part, show domain
		parts := strings.Split(original, "@")
		if len(parts) == 2 {
			localPart := parts[0]
			if len(localPart) > 2 {
				return localPart[:2] + "***@" + parts[1]
			}
			return "***@" + parts[1]
		}
		return "[MASKED]"
	case "phone":
		// Show only last 4 digits
		digits := regexp.MustCompile(`\d`).FindAllString(original, -1)
		if len(digits) >= 4 {
			return "***-***-" + strings.Join(digits[len(digits)-4:], "")
		}
		return "[MASKED]"
	case "ip_address":
		// Mask the last octet
		parts := strings.Split(original, ".")
		if len(parts) == 4 {
			return parts[0] + "." + parts[1] + "." + parts[2] + ".***"
		}
		return "[MASKED]"
	case "jwt":
		// Show only first part of JWT
		parts := strings.Split(original, ".")
		if len(parts) >= 1 {
			return parts[0] + ".[MASKED]"
		}
		return "[MASKED]"
	default:
		return "[MASKED]"
	}
}

// MaskFields masks sensitive fields in a map
func (m *SensitiveDataMasker) MaskFields(fields map[string]interface{}) map[string]interface{} {
	maskedFields := make(map[string]interface{})
	
	for key, value := range fields {
		if m.isSensitiveField(key) {
			maskedFields[key] = "[MASKED]"
		} else if strValue, ok := value.(string); ok {
			maskedFields[key] = m.MaskSensitiveData(strValue)
		} else {
			maskedFields[key] = value
		}
	}
	
	return maskedFields
}

// isSensitiveField checks if a field name indicates sensitive data
func (m *SensitiveDataMasker) isSensitiveField(fieldName string) bool {
	sensitiveFields := []string{
		"password", "pass", "pwd",
		"api_key", "apikey", "api_key",
		"token", "access_token", "bearer",
		"secret", "secret_key",
		"credit_card", "card_number",
		"ssn", "social_security",
		"email", "phone", "phone_number",
		"ip", "ip_address",
		"jwt", "jwt_token",
	}
	
	fieldLower := strings.ToLower(fieldName)
	for _, sensitiveField := range sensitiveFields {
		if strings.Contains(fieldLower, sensitiveField) {
			return true
		}
	}
	
	return false
}

// Global masker instance
var GlobalMasker = NewSensitiveDataMasker()

// MaskSensitiveData is a convenience function for global masking
func MaskSensitiveData(data string) string {
	return GlobalMasker.MaskSensitiveData(data)
}

// MaskFields is a convenience function for global field masking
func MaskFields(fields map[string]interface{}) map[string]interface{} {
	return GlobalMasker.MaskFields(fields)
}
