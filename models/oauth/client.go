package oauth

type Client struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
}

func NewClient(clientId string, clientSecret string, redirectURI string) *Client {
	return &Client{ClientId: clientId, ClientSecret: clientSecret, RedirectURI: redirectURI}
}
