package fsnativeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsauth"
)

const baseURL = "https://platform.fatsecret.com/rest"

type FSNativeAPIClient struct {
	fsSignedAuth fsauth.SignedAuth
	httpClient   *http.Client
}

func NewNativeAPIClient(auth fsauth.SignedAuth) *FSNativeAPIClient {
	return &FSNativeAPIClient{
		fsSignedAuth: auth,
		httpClient:   http.DefaultClient,
	}
}

// get performs an authenticated GET request to the given API path with the
// provided query parameters. All target endpoints for this client use GET
// with query params and return JSON. path should not include a leading slash
// (e.g. "foods/search/v1").
func (c *FSNativeAPIClient) get(path string, params map[string]string) ([]byte, error) {
	rawURL := baseURL + "/" + path

	authHeader, err := c.fsSignedAuth.AuthorizationHeader("GET", rawURL, params)
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: build auth header: %w", err)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: parse url: %w", err)
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: create request: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fsnativeapi: unexpected status %d: %s", resp.StatusCode, body)
	}

	return body, nil
}

// post performs an authenticated POST request to the given API path, serializing
// body as JSON. OAuth 1.0a signs over the URL only — JSON bodies are not included
// in the signature base string per the OAuth spec.
func (c *FSNativeAPIClient) post(path string, body any) ([]byte, error) {
	rawURL := baseURL + "/" + path

	authHeader, err := c.fsSignedAuth.AuthorizationHeader("POST", rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: build auth header: %w", err)
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: create request: %w", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fsnativeapi: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fsnativeapi: unexpected status %d: %s", resp.StatusCode, respBody)
	}

	return respBody, nil
}