package fsfetchtoken

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cherrysuryp/fatsecret-mcp/internal/config"
)

const (
	requestTokenURL = "https://authentication.fatsecret.com/oauth/request_token"
	authorizeURL    = "https://authentication.fatsecret.com/oauth/authorize"
	accessTokenURL  = "https://authentication.fatsecret.com/oauth/access_token"
)

type oauthFlow struct {
	cfg            *config.Config
	consumerKey    string
	consumerSecret string
	httpClient     *http.Client
}

// percentEncode encodes s per RFC 3986 (unreserved chars pass through, rest → %XX).
// Must not use url.QueryEscape — it encodes space as '+'.
func percentEncode(s string) string {
	var buf strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') ||
			b == '-' || b == '_' || b == '.' || b == '~' {
			buf.WriteByte(b)
		} else {
			fmt.Fprintf(&buf, "%%%02X", b)
		}
	}
	return buf.String()
}

// signedParams builds the full OAuth param map (oauth_ params + requestParams +
// oauth_signature) ready to send as a form body or query string. tokenKey and
// tokenSecret are the user-level OAuth token pair; pass empty strings for
// two-legged requests such as fetching a request token.
func (f *oauthFlow) signedParams(method, rawURL string, requestParams map[string]string, oauthToken, tokenSecret string) (map[string]string, error) {
	nonce, err := nonce()
	if err != nil {
		return nil, err
	}

	all := map[string]string{
		"oauth_consumer_key":     f.consumerKey,
		"oauth_nonce":            nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
		"oauth_version":          "1.0",
	}
	if oauthToken != "" {
		all["oauth_token"] = oauthToken
	}
	for k, v := range requestParams {
		all[k] = v
	}

	// Build normalized param string
	pairs := make([]string, 0, len(all))
	for k, v := range all {
		pairs = append(pairs, percentEncode(k)+"="+percentEncode(v))
	}
	sort.Strings(pairs)
	normalized := strings.Join(pairs, "&")

	// Build signature base string
	baseString := strings.ToUpper(method) + "&" + percentEncode(rawURL) + "&" + percentEncode(normalized)

	// HMAC-SHA1 sign
	signingKey := percentEncode(f.consumerSecret) + "&" + percentEncode(tokenSecret)
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(baseString))
	all["oauth_signature"] = base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return all, nil
}

func nonce() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("fsfetchtoken: generate nonce: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// parseFormResponse decodes a form-encoded response body into a map.
func parseFormResponse(body []byte) (map[string]string, error) {
	vals, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(vals))
	for k, v := range vals {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result, nil
}

// requestToken performs step 1 of the OAuth flow: exchanges consumer credentials
// for a temporary request token using a form-encoded POST body.
func (f *oauthFlow) requestToken() (string, string, error) {
	params, err := f.signedParams("POST", requestTokenURL, map[string]string{"oauth_callback": "oob"}, "", "")
	if err != nil {
		return "", "", err
	}

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	req, err := http.NewRequest(http.MethodPost, requestTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", "", fmt.Errorf("fsfetchtoken: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("fsfetchtoken: request token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("fsfetchtoken: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("fsfetchtoken: request token returned %d: %s", resp.StatusCode, body)
	}

	parsed, err := parseFormResponse(body)
	if err != nil {
		return "", "", fmt.Errorf("fsfetchtoken: parse request token response: %w", err)
	}

	token := parsed["oauth_token"]
	secret := parsed["oauth_token_secret"]
	if token == "" || secret == "" {
		return "", "", fmt.Errorf("fsfetchtoken: missing token or secret in response: %s", body)
	}
	return token, secret, nil
}

// accessToken performs step 3 of the OAuth flow: exchanges the request token
// and verifier for a permanent access token using a query-string GET request.
func (f *oauthFlow) accessToken(requestToken, requestSecret, verifier string) (string, string, string, error) {
	params, err := f.signedParams("GET", accessTokenURL, map[string]string{"oauth_verifier": verifier}, requestToken, requestSecret)
	if err != nil {
		return "", "", "", err
	}

	u, _ := url.Parse(accessTokenURL)
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", "", "", fmt.Errorf("fsfetchtoken: build request: %w", err)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("fsfetchtoken: access token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("fsfetchtoken: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("fsfetchtoken: access token returned %d: %s", resp.StatusCode, body)
	}

	parsed, err := parseFormResponse(body)
	if err != nil {
		return "", "", "", fmt.Errorf("fsfetchtoken: parse access token response: %w", err)
	}

	token := parsed["oauth_token"]
	secret := parsed["oauth_token_secret"]
	userID := parsed["user_id"]
	if token == "" || secret == "" {
		return "", "", "", fmt.Errorf("fsfetchtoken: missing token or secret in response: %s", body)
	}
	return token, secret, userID, nil
}

// saveUserConfig writes the given credentials to the user config file via the
// config package, which owns path resolution and file permissions.
func (f *oauthFlow) saveUserConfig(userID, accessToken, secretToken string) error {
	return f.cfg.SaveUserConfig(config.FSAPIUserConfig{
		UserID:      userID,
		AccessToken: accessToken,
		SecretToken: secretToken,
	})
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

func (f *oauthFlow) run() error {
	fmt.Println("FatSecret OAuth Setup Utility")
	fmt.Println("Consumer ID:", f.consumerKey)
	fmt.Println()

	fmt.Println("Step 1: Getting request token...")
	requestToken, requestSecret, err := f.requestToken()
	if err != nil {
		return err
	}
	fmt.Println("✓ Request token obtained")

	fmt.Println("\nStep 2: User authorization")
	authURL := fmt.Sprintf("%s?oauth_token=%s", authorizeURL, requestToken)
	fmt.Printf("Authorization URL: %s\n\n", authURL)

	if !openBrowser(authURL) {
		fmt.Println("Could not open browser automatically. Please visit the URL above manually.")
	}

	fmt.Println("Instructions:")
	fmt.Println("1. Log in to your FatSecret account (or create one)")
	fmt.Println("2. Click \"Allow\" to authorize this application")
	fmt.Println("3. Copy the verifier code shown on the page")

	verifier, err := prompt("\nEnter the verifier code: ")
	if err != nil {
		return err
	}
	if verifier == "" {
		return fmt.Errorf("verifier code is required")
	}

	fmt.Println("\nStep 3: Getting access token...")
	accessToken, secretToken, userID, err := f.accessToken(requestToken, requestSecret, verifier)
	if err != nil {
		return err
	}
	fmt.Println("✓ Access token obtained")
	if userID != "" {
		fmt.Println("User ID:", userID)
	}

	if err := f.saveUserConfig(userID, accessToken, secretToken); err != nil {
		return err
	}
	fmt.Printf("✓ Credentials saved to: %s\n", f.cfg.FSAPIClientConfig.UserConfigPath)

	fmt.Println("\nSetup complete! You can now use the FatSecret MCP server.")
	return nil
}

// FetchToken runs the interactive OAuth 1.0 three-legged flow, saves the
// resulting access credentials to the user config file, and exits on error.
func FetchToken() {
	cfg := config.MustLoadConfig()

	flow := &oauthFlow{
		cfg:            cfg,
		consumerKey:    cfg.FSAPIClientConfig.ConsumerID,
		consumerSecret: cfg.FSAPIClientConfig.ConsumerSecret,
		httpClient:     http.DefaultClient,
	}

	if err := flow.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}