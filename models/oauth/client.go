package oauth

type Client struct {
	ClientId                string
	ClientSecret            string
	ClientName              string
	GrantTypes              string
	ResponseTypes           string
	TokenEndpointAuthMethod string
	RedirectUris            string
}
