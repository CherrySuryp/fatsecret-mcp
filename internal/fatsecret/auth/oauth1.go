package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// OAuth1Client signs and executes HTTP requests using OAuth 1.0a HMAC-SHA1.
type OAuth1Client struct {
	config *Config
}

// NewOAuth1Client returns an OAuth1Client configured with the provided credentials.
func NewOAuth1Client(cfg *Config) *OAuth1Client {
	return &OAuth1Client{config: cfg}
}

// MakeRequest signs and executes an OAuth 1.0a request against rawURL.
// extraParams are the business-level query or body parameters for the endpoint
// (e.g. {"format": "json"}); OAuth protocol params are added automatically.
// token and tokenSecret are the user-level OAuth token pair; pass empty strings
// for two-legged (app-only) requests such as fetching a request token.
// The response body is parsed as JSON first, falling back to query-string form,
// and returned as a flat string map.
func (c *OAuth1Client) MakeRequest(
	method string,
	rawURL string,
	extraParams map[string]string,
	token string,
	tokenSecret string,
) (map[string]string, error) {
	oauthParams := c.buildOAuthParams(token)
	if token == "" {
		delete(oauthParams, "oauth_token")
	}

	allParams := make(map[string]string)
	for k, v := range oauthParams {
		allParams[k] = v
	}
	for k, v := range extraParams {
		allParams[k] = v
	}

	sig := generateSignature(method, rawURL, allParams, c.config.ClientSecret, tokenSecret)
	allParams["oauth_signature"] = sig

	var req *http.Request
	var err error

	if method == "GET" {
		// For GET requests all parameters (OAuth + extra) are appended as query
		// string. FatSecret's REST API accepts this form for read-only calls.
		vals := url.Values{}
		for k, v := range allParams {
			vals.Set(k, v)
		}
		req, err = http.NewRequest("GET", rawURL+"?"+vals.Encode(), nil)
	} else {
		// For POST requests parameters are sent as an application/x-www-form-urlencoded
		// body. This is required for OAuth token-exchange endpoints (request token,
		// access token) which do not accept query-string parameters.
		vals := url.Values{}
		for k, v := range allParams {
			vals.Set(k, v)
		}
		body := vals.Encode()
		req, err = http.NewRequest("POST", rawURL, strings.NewReader(body))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal(body, &jsonResult); err == nil {
		result := make(map[string]string)
		for k, v := range jsonResult {
			result[k] = fmt.Sprintf("%v", v)
		}
		return result, nil
	}

	parsed, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %s", string(body))
	}
	result := make(map[string]string)
	for k, vs := range parsed {
		if len(vs) > 0 {
			result[k] = vs[0]
		}
	}
	return result, nil
}

// buildOAuthParams returns the standard OAuth 1.0a protocol parameters for a
// request. When token is empty the oauth_token field is omitted by the caller.
func (c *OAuth1Client) buildOAuthParams(token string) map[string]string {
	return map[string]string{
		"oauth_consumer_key":     c.config.ClientID,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_version":          "1.0",
		"oauth_token":            token,
	}
}

// generateNonce returns a cryptographically random 32-character hex string
// suitable for use as an OAuth nonce.
func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// percentEncode applies RFC 3986 percent-encoding to s, leaving unreserved
// characters (A-Z a-z 0-9 - _ . ~) unencoded as required by OAuth 1.0a.
func percentEncode(s string) string {
	var buf strings.Builder
	for _, b := range []byte(s) {
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
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') || b == '-' || b == '_' || b == '.' || b == '~'
}

// signatureBaseString constructs the OAuth 1.0a signature base string:
// METHOD & percent-encoded URL & percent-encoded normalized parameters.
// Parameters are sorted lexicographically before normalization as required by
// the spec (RFC 5849 §3.4.1).
func signatureBaseString(method, rawURL string, params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, percentEncode(k)+"="+percentEncode(params[k]))
	}
	normalizedParams := strings.Join(parts, "&")

	return strings.Join([]string{
		strings.ToUpper(method),
		percentEncode(rawURL),
		percentEncode(normalizedParams),
	}, "&")
}

// generateSignature computes the HMAC-SHA1 OAuth signature and returns it
// base64-encoded. The signing key is clientSecret&tokenSecret (both
// percent-encoded), where tokenSecret is empty for request-token calls.
func generateSignature(method, rawURL string, params map[string]string, clientSecret, tokenSecret string) string {
	base := signatureBaseString(method, rawURL, params)
	signingKey := percentEncode(clientSecret) + "&" + percentEncode(tokenSecret)
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(base))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

