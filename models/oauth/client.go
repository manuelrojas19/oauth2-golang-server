package oauth

type Client struct {
	ClientId                string
	ClientSecret            string
	ClientName              string
	GrantTypes              []GrantType
	ResponseTypes           []ResponseType
	TokenEndpointAuthMethod TokenEndpointAuthMethod
	RedirectUris            []string
}
