package fsauth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"maps"
	"sort"
	"strings"
)

// generateNonce produces a cryptographically random 16-byte hex string.
func generateNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// percentEncode encodes s per RFC 3986: unreserved characters pass through,
// all other bytes are encoded as %XX (uppercase hex).
// Must not use url.QueryEscape — it encodes space as '+' which violates the OAuth spec.
func percentEncode(s string) string {
	var buf strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if isUnreserved(b) {
			buf.WriteByte(b)
		} else {
			fmt.Fprintf(&buf, "%%%02X", b)
		}
	}
	return buf.String()
}

// isUnreserved reports whether b is an RFC 3986 unreserved character.
func isUnreserved(b byte) bool {
	return (b >= 'A' && b <= 'Z') ||
		(b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') ||
		b == '-' || b == '_' || b == '.' || b == '~'
}

// normalizedParamString percent-encodes all keys and values in params, sorts the
// resulting key=value pairs lexicographically, and joins them with '&'.
func normalizedParamString(params map[string]string) string {
	pairs := make([]string, 0, len(params))
	for k, v := range params {
		pairs = append(pairs, percentEncode(k)+"="+percentEncode(v))
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "&")
}

// signatureBaseString constructs the OAuth 1.0 signature base string.
func signatureBaseString(method, rawURL, normalizedParams string) string {
	return strings.ToUpper(method) + "&" + percentEncode(rawURL) + "&" + percentEncode(normalizedParams)
}

// buildAuthHeader formats the OAuth Authorization header value from the given
// OAuth parameters (which must include oauth_signature).
func buildAuthHeader(oauthParams map[string]string) string {
	keys := make([]string, 0, len(oauthParams))
	for k := range oauthParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+`="`+percentEncode(oauthParams[k])+`"`)
	}
	return "OAuth " + strings.Join(parts, ", ")
}

// mergeParams returns a new map containing all key-value pairs from both a and b.
// Neither input map is mutated.
func mergeParams(a, b map[string]string) map[string]string {
	merged := make(map[string]string, len(a)+len(b))
	maps.Copy(merged, a)
	maps.Copy(merged, b)
	return merged
}
