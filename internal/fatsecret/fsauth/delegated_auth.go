package fsauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

// DelegatedAuth performs OAuth 1.0 signing for signed-and-delegated requests,
// which operate on behalf of a specific user profile.
type DelegatedAuth struct {
	BaseAuth
	oauthToken       string
	oauthTokenSecret string
}

// NewDelegatedAuth creates a new DelegatedAuth with the given OAuth 1.0 credentials.
func NewDelegatedAuth(consumerKey, consumerSecret, oauthToken, oauthTokenSecret string) *DelegatedAuth {
	return &DelegatedAuth{
		BaseAuth: BaseAuth{
			consumerKey:    consumerKey,
			consumerSecret: consumerSecret,
		},
		oauthToken:       oauthToken,
		oauthTokenSecret: oauthTokenSecret,
	}
}

// AuthorizationHeader returns the OAuth 1.0 Authorization header value for a
// delegated request. method is the HTTP verb, rawURL is the base URL without
// query string, and requestParams are any additional parameters (query/body)
// that must be included in the signature. The returned header value is ready
// for use with req.Header.Set("Authorization", header).
func (c *DelegatedAuth) AuthorizationHeader(method, rawURL string, requestParams map[string]string) (string, error) {
	oauthP, err := c.oauthParams()
	if err != nil {
		return "", err
	}

	allParams := mergeParams(oauthP, requestParams)
	normalized := normalizedParamString(allParams)
	baseString := signatureBaseString(method, rawURL, normalized)
	signature := c.sign(baseString)

	oauthP["oauth_signature"] = signature

	return buildAuthHeader(oauthP), nil
}

// oauthParams builds a fresh map of the six OAuth protocol parameters,
// including oauth_token for the delegated user.
func (c *DelegatedAuth) oauthParams() (map[string]string, error) {
	nonce, err := generateNonce()
	if err != nil {
		return nil, fmt.Errorf("fsauth: failed to generate nonce: %w", err)
	}

	return map[string]string{
		"oauth_consumer_key":     c.consumerKey,
		"oauth_nonce":            nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
		"oauth_token":            c.oauthToken,
		"oauth_version":          "1.0",
	}, nil
}

// sign computes the HMAC-SHA1 signature of baseString using both the consumer
// secret and the user's token secret.
func (c *DelegatedAuth) sign(baseString string) string {
	signingKey := percentEncode(c.consumerSecret) + "&" + percentEncode(c.oauthTokenSecret)
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(baseString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
