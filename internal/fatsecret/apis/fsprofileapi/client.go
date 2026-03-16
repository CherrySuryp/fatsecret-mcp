package fsprofileapi

import "github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsauth"

const baseURL = "https://platform.fatsecret.com/rest/"

type Client struct {
	accessToken       string
	accessTokenSecret string
	client            *fsauth.FSOAuth1Client
}

func New(AccessToken, AccessTokenSecret string, Oauth1Client *fsauth.FSOAuth1Client) *Client {
	return &Client{
		accessToken:       AccessToken,
		accessTokenSecret: AccessTokenSecret,
		client:            Oauth1Client,
	}
}

func (c *Client) get(path string, params map[string]string) (map[string]any, error) {
	params["format"] = "json"
	return c.client.MakeRequest("GET", baseURL+path, params, c.accessToken, c.accessTokenSecret)
}

func (c *Client) post(path string, params map[string]string) (map[string]any, error) {
	params["format"] = "json"
	return c.client.MakeRequest("POST", baseURL+path, params, c.accessToken, c.accessTokenSecret)
}

func (c *Client) delete(path string, params map[string]string) (map[string]any, error) {
	params["format"] = "json"
	return c.client.MakeRequest("DELETE", baseURL+path, params, c.accessToken, c.accessTokenSecret)
}

func (c *Client) put(path string, params map[string]string) (map[string]any, error) {
	params["format"] = "json"
	return c.client.MakeRequest("PUT", baseURL+path, params, c.accessToken, c.accessTokenSecret)
}
