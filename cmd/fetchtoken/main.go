package main

import (
	"bufio"
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
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

const (
	apiURL          = "https://platform.fatsecret.com/rest/server.api"
	requestTokenURL = "https://authentication.fatsecret.com/oauth/request_token"
	authorizeURL    = "https://authentication.fatsecret.com/oauth/authorize"
	accessTokenURL  = "https://authentication.fatsecret.com/oauth/access_token"
)

type Config struct {
	ClientID          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AccessToken       string `json:"accessToken,omitempty"`
	AccessTokenSecret string `json:"accessTokenSecret,omitempty"`
	UserID            string `json:"userId,omitempty"`
}

type OAuthClient struct {
	config     Config
	configPath string
}

func NewOAuthClient() *OAuthClient {
	home, _ := os.UserHomeDir()
	return &OAuthClient{
		configPath: filepath.Join(home, ".fatsecret-mcp-config.json"),
		config: Config{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		},
	}
}

func (c *OAuthClient) loadConfig() error {
	data, err := os.ReadFile(c.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}
	// Merge: env vars take precedence for credentials if set
	if c.config.ClientID == "" {
		c.config.ClientID = loaded.ClientID
	}
	if c.config.ClientSecret == "" {
		c.config.ClientSecret = loaded.ClientSecret
	}
	c.config.AccessToken = loaded.AccessToken
	c.config.AccessTokenSecret = loaded.AccessTokenSecret
	c.config.UserID = loaded.UserID
	return nil
}

func (c *OAuthClient) saveConfig() error {
	data, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.configPath, data, 0600)
}

func prompt(question string) (string, error) {
	fmt.Print(question)
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(answer), nil
}

func openBrowser(rawURL string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	return cmd.Start() == nil
}

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// percentEncode implements RFC 3986 percent encoding for OAuth 1.0a.
// Only unreserved characters (A-Z a-z 0-9 - _ . ~) are left unencoded.
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

func isUnreserved(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') || b == '-' || b == '_' || b == '.' || b == '~'
}

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

func generateSignature(method, rawURL string, params map[string]string, clientSecret, tokenSecret string) string {
	base := signatureBaseString(method, rawURL, params)
	signingKey := percentEncode(clientSecret) + "&" + percentEncode(tokenSecret)

	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(base))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *OAuthClient) buildOAuthParams(token string) map[string]string {
	return map[string]string{
		"oauth_consumer_key":     c.config.ClientID,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_version":          "1.0",
		"oauth_token":            token,
	}
}

func (c *OAuthClient) makeOAuthRequest(method, rawURL string, extraParams map[string]string, token, tokenSecret string) (map[string]string, error) {
	oauthParams := c.buildOAuthParams(token)
	if token == "" {
		delete(oauthParams, "oauth_token")
	}

	// Merge all params for signature
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
		vals := url.Values{}
		for k, v := range allParams {
			vals.Set(k, v)
		}
		req, err = http.NewRequest("GET", rawURL+"?"+vals.Encode(), nil)
	} else {
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

	// Try JSON first, fall back to query string
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

func (c *OAuthClient) setupCredentials() error {
	fmt.Println("=== FatSecret API Credentials Setup ===")

	if c.config.ClientID != "" && c.config.ClientSecret != "" {
		fmt.Println("Existing credentials found:")
		fmt.Printf("Client ID: %s\n", c.config.ClientID)
		fmt.Printf("Client Secret: %s...\n", c.config.ClientSecret[:min(8, len(c.config.ClientSecret))])

		answer, err := prompt("\nUse existing credentials? (y/n): ")
		if err != nil {
			return err
		}
		if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
			return nil
		}
	}

	fmt.Println("Please enter your FatSecret API credentials.")
	fmt.Println("You can get these from: https://platform.fatsecret.com/")

	clientID, err := prompt("Client ID: ")
	if err != nil {
		return err
	}
	clientSecret, err := prompt("Client Secret: ")
	if err != nil {
		return err
	}

	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("Client ID and Client Secret are required")
	}

	c.config.ClientID = clientID
	c.config.ClientSecret = clientSecret

	if err := c.saveConfig(); err != nil {
		return err
	}
	fmt.Println("✓ Credentials saved successfully")
	return nil
}

func (c *OAuthClient) runOAuthFlow() error {
	fmt.Println("=== Starting OAuth Flow ===")

	// Step 1: Get request token
	fmt.Println("Step 1: Getting request token...")

	response, err := c.makeOAuthRequest(
		"POST",
		requestTokenURL,
		map[string]string{"oauth_callback": "oob"},
		"",
		"",
	)
	if err != nil {
		return fmt.Errorf("failed to get request token: %w", err)
	}

	requestToken := response["oauth_token"]
	requestTokenSecret := response["oauth_token_secret"]

	if requestToken == "" || requestTokenSecret == "" {
		return fmt.Errorf("invalid response: missing token or token secret")
	}
	fmt.Println("✓ Request token obtained")

	// Step 2: User authorization
	fmt.Println("Step 2: User authorization")
	authURL := fmt.Sprintf("%s?oauth_token=%s", authorizeURL, requestToken)

	fmt.Println("Opening authorization URL in your browser...")
	fmt.Printf("URL: %s\n\n", authURL)

	if !openBrowser(authURL) {
		fmt.Println("Could not open browser automatically. Please visit the URL above manually.")
	}

	fmt.Println("Instructions:")
	fmt.Println("1. Log in to your FatSecret account (or create one)")
	fmt.Println("2. Click \"Allow\" to authorize this application")
	fmt.Println("3. Copy the verifier code from the page")
	fmt.Println("4. Paste it below")

	verifier, err := prompt("Enter the verifier code: ")
	if err != nil {
		return err
	}
	if verifier == "" {
		return fmt.Errorf("verifier code is required")
	}

	// Step 3: Get access token
	fmt.Println("\nStep 3: Getting access token...")

	accessResponse, err := c.makeOAuthRequest("GET", accessTokenURL, map[string]string{
		"oauth_verifier": verifier,
	}, requestToken, requestTokenSecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	if accessResponse["oauth_token"] == "" || accessResponse["oauth_token_secret"] == "" {
		return fmt.Errorf("invalid response from access token endpoint. Please try again")
	}

	c.config.AccessToken = accessResponse["oauth_token"]
	c.config.AccessTokenSecret = accessResponse["oauth_token_secret"]
	c.config.UserID = accessResponse["user_id"]

	if err := c.saveConfig(); err != nil {
		return err
	}

	fmt.Println("✓ Access token obtained")
	fmt.Println("✓ OAuth flow completed successfully!")
	fmt.Printf("User ID: %s\n", c.config.UserID)
	fmt.Printf("Authentication details saved to: %s\n", c.configPath)

	// Step 4: Verify credentials with a test API call
	fmt.Println("\nStep 4: Verifying credentials...")
	if err := c.verifyCredentials(); err != nil {
		return fmt.Errorf("credential verification failed: %w", err)
	}
	return nil
}

func (c *OAuthClient) verifyCredentials() error {
	result, err := c.makeOAuthRequest("GET", apiURL, map[string]string{
		"method": "profile.get",
		"format": "json",
	}, c.config.AccessToken, c.config.AccessTokenSecret)
	if err != nil {
		return err
	}
	if errMsg, ok := result["error"]; ok {
		return fmt.Errorf("API error: %s", errMsg)
	}
	fmt.Println("✓ Credentials verified successfully")
	return nil
}

func (c *OAuthClient) checkStatus() {
	fmt.Println("=== Authentication Status ===")

	hasCredentials := c.config.ClientID != "" && c.config.ClientSecret != ""
	hasAccessToken := c.config.AccessToken != "" && c.config.AccessTokenSecret != ""

	check := func(ok bool) string {
		if ok {
			return "✓"
		}
		return "✗"
	}

	fmt.Printf("Credentials configured: %s\n", check(hasCredentials))
	fmt.Printf("User authenticated: %s\n", check(hasAccessToken))

	if hasCredentials {
		fmt.Printf("Client ID: %s\n", c.config.ClientID)
	}
	if hasAccessToken {
		userID := c.config.UserID
		if userID == "" {
			userID = "N/A"
		}
		fmt.Printf("User ID: %s\n", userID)
	}
	fmt.Printf("Config file: %s\n", c.configPath)
}

func (c *OAuthClient) run() error {
	if err := c.loadConfig(); err != nil {
		return err
	}

	fmt.Println("FatSecret OAuth Console Utility")
	fmt.Println("This utility will help you authenticate with the FatSecret API.")

	if err := c.setupCredentials(); err != nil {
		return err
	}

	answer, err := prompt("Do you want to authenticate a user now? (y/n): ")
	if err != nil {
		return err
	}
	if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
		if err := c.runOAuthFlow(); err != nil {
			return err
		}
	}

	fmt.Println()
	c.checkStatus()
	fmt.Println("\nSetup complete! You can now use the FatSecret MCP server.")
	return nil
}

func main() {
	client := NewOAuthClient()
	if err := client.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
