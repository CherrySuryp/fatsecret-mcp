package fsfetchtoken

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsauth"
)

const (
	profileURL      = "https://platform.fatsecret.com/rest/profile/v1"
	requestTokenURL = "https://authentication.fatsecret.com/oauth/request_token"
	authorizeURL    = "https://authentication.fatsecret.com/oauth/authorize"
	accessTokenURL  = "https://authentication.fatsecret.com/oauth/access_token"
)

type OAuthClient struct {
	config *fsauth.Config
	oauth1 *fsauth.OAuth1Client
}

func NewOAuthClient() *OAuthClient {
	cfg := &fsauth.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
	return &OAuthClient{
		config: cfg,
		oauth1: fsauth.NewOAuth1Client(cfg),
	}
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

	if err := fsauth.SaveConfig(c.config); err != nil {
		return err
	}
	fmt.Println("✓ Credentials saved successfully")
	return nil
}

func (c *OAuthClient) runOAuthFlow() error {
	fmt.Println("=== Starting OAuth Flow ===")

	fmt.Println("Step 1: Getting request token...")
	response, err := c.oauth1.MakeRequest("POST", requestTokenURL, map[string]string{"oauth_callback": "oob"}, "", "")
	if err != nil {
		return fmt.Errorf("failed to get request token: %w", err)
	}

	requestToken := response["oauth_token"]
	requestTokenSecret := response["oauth_token_secret"]

	if requestToken == "" || requestTokenSecret == "" {
		return fmt.Errorf("invalid response: missing token or token secret")
	}
	fmt.Println("✓ Request token obtained")

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

	fmt.Println("\nStep 3: Getting access token...")
	accessResponse, err := c.oauth1.MakeRequest("GET", accessTokenURL, map[string]string{
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

	if err := fsauth.SaveConfig(c.config); err != nil {
		return err
	}

	fmt.Println("✓ Access token obtained")
	fmt.Println("✓ OAuth flow completed successfully!")
	fmt.Printf("User ID: %s\n", c.config.UserID)
	fmt.Printf("Authentication details saved to: %s\n", fsauth.ConfigPath())

	fmt.Println("\nStep 4: Verifying credentials...")
	if err := c.verifyCredentialsGet(); err != nil {
		return fmt.Errorf("credential verification failed: %w", err)
	}
	return nil
}

// verifyCredentialsGet confirms that the stored OAuth credentials are valid by
// calling the REST profile endpoint. It is used after the OAuth flow completes
// to ensure the obtained access token and secret are accepted by the API.
func (c *OAuthClient) verifyCredentialsGet() error {
	result, err := c.oauth1.MakeRequest("GET", profileURL, map[string]string{
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
	fmt.Printf("Config file: %s\n", fsauth.ConfigPath())
}

func (c *OAuthClient) run() error {
	cfg, err := fsauth.LoadConfig()
	if err != nil {
		return err
	}
	c.config = cfg
	c.oauth1 = fsauth.NewOAuth1Client(cfg)

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

func FetchToken() {
	client := NewOAuthClient()
	if err := client.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
