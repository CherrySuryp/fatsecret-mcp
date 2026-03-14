package fsclient

import "github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsauth"

const baseURL = "https://platform.fatsecret.com/rest/"

type Client struct {
	cfg   *fsauth.Config
	oauth *fsauth.FSOAuth1Client
}

func NewClient(cfg *fsauth.Config) *Client {
	return &Client{
		cfg:   cfg,
		oauth: fsauth.NewFSOAuth1Client(cfg),
	}
}

func (c *Client) get(path string, params map[string]string) (map[string]any, error) {
	params["format"] = "json"
	return c.oauth.MakeRequest("GET", baseURL+path, params, c.cfg.AccessToken, c.cfg.AccessTokenSecret)
}

func (c *Client) post(path string, params map[string]string) (map[string]any, error) {
	params["format"] = "json"
	return c.oauth.MakeRequest("POST", baseURL+path, params, c.cfg.AccessToken, c.cfg.AccessTokenSecret)
}

func (c *Client) delete(path string, params map[string]string) (map[string]any, error) {
	params["format"] = "json"
	return c.oauth.MakeRequest("DELETE", baseURL+path, params, c.cfg.AccessToken, c.cfg.AccessTokenSecret)
}

func (c *Client) put(path string, params map[string]string) (map[string]any, error) {
	params["format"] = "json"
	return c.oauth.MakeRequest("PUT", baseURL+path, params, c.cfg.AccessToken, c.cfg.AccessTokenSecret)
}
