package fsauth

// BaseAuth holds the OAuth 1.0 consumer credentials shared across client types.
type BaseAuth struct {
	consumerKey    string
	consumerSecret string
}
