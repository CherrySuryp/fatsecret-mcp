package fsprofileapi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsauth"
)

const baseURL = "https://platform.fatsecret.com/rest/"

// FSProfileClient performs OAuth 1.0 3-legged (delegated) requests against the
// FatSecret profile API — endpoints that operate on behalf of a specific user.
type FSProfileClient struct {
	delegatedAuth *fsauth.DelegatedAuth
	httpClient    *http.Client
}

// NewProfileClient creates a new FSProfileClient using the given delegated auth credentials.
func NewProfileClient(auth *fsauth.DelegatedAuth) *FSProfileClient {
	return &FSProfileClient{
		delegatedAuth: auth,
		httpClient:    http.DefaultClient,
	}
}

// get performs an authenticated GET request to the given API path with the
// provided query parameters. params are included in the OAuth signature and
// appended as query string. path should not include a leading slash
// (e.g. "food/favorites/v2").
func (c *FSProfileClient) get(path string, params map[string]string) ([]byte, error) {
	rawURL := baseURL + path

	authHeader, err := c.delegatedAuth.AuthorizationHeader("GET", rawURL, params)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: build auth header: %w", err)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: parse url: %w", err)
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: create request: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	return c.do(req)
}

// post performs an authenticated POST request to the given API path with the
// provided parameters as a form-encoded body. params are included in the OAuth
// signature per the OAuth 1.0a spec for application/x-www-form-urlencoded bodies.
func (c *FSProfileClient) post(path string, params map[string]string) ([]byte, error) {
	rawURL := baseURL + "/" + path

	authHeader, err := c.delegatedAuth.AuthorizationHeader("POST", rawURL, params)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: build auth header: %w", err)
	}

	formBody := encodeParams(params)
	req, err := http.NewRequest(http.MethodPost, rawURL, strings.NewReader(formBody))
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: create request: %w", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

// put performs an authenticated PUT request to the given API path with the
// provided parameters as a form-encoded body. params are included in the OAuth
// signature per the OAuth 1.0a spec for application/x-www-form-urlencoded bodies.
func (c *FSProfileClient) put(path string, params map[string]string) ([]byte, error) {
	rawURL := baseURL + "/" + path

	authHeader, err := c.delegatedAuth.AuthorizationHeader("PUT", rawURL, params)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: build auth header: %w", err)
	}

	formBody := encodeParams(params)
	req, err := http.NewRequest(http.MethodPut, rawURL, strings.NewReader(formBody))
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: create request: %w", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req)
}

// delete performs an authenticated DELETE request to the given API path with the
// provided parameters as query string. params are included in the OAuth signature.
func (c *FSProfileClient) delete(path string, params map[string]string) ([]byte, error) {
	rawURL := baseURL + "/" + path

	authHeader, err := c.delegatedAuth.AuthorizationHeader("DELETE", rawURL, params)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: build auth header: %w", err)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: parse url: %w", err)
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: create request: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	return c.do(req)
}

// do executes the request and returns the response body, or an error for
// non-200 responses.
func (c *FSProfileClient) do(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fsprofileclient: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fsprofileclient: unexpected status %d: %s", resp.StatusCode, body)
	}

	return body, nil
}

// encodeParams URL-encodes a map of string key/value pairs into a form body.
func encodeParams(params map[string]string) string {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return vals.Encode()
}

