package fsclient

import "github.com/cherrysuryp/fatsecret-mcp/internal/fatsecret/fsauth"

const baseURL = "https://platform.fatsecret.com/rest/"

type FSClient struct {
	cfg   *fsauth.Config
	oauth *fsauth.FSOAuth1Client
}

func NewFSClient(cfg *fsauth.Config) *FSClient {
	return &FSClient{
		cfg:   cfg,
		oauth: fsauth.NewFSOAuth1Client(cfg),
	}
}

func (c *FSClient) get(path string, params map[string]string) (map[string]interface{}, error) {
	params["format"] = "json"
	return c.oauth.MakeRequest("GET", baseURL+path, params, c.cfg.AccessToken, c.cfg.AccessTokenSecret)
}

func (c *FSClient) post(path string, params map[string]string) (map[string]interface{}, error) {
	params["format"] = "json"
	return c.oauth.MakeRequest("POST", baseURL+path, params, c.cfg.AccessToken, c.cfg.AccessTokenSecret)
}
