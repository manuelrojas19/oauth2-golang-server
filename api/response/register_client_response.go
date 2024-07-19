package response

type RegisterClientResponse struct {
	ClientId                string `json:"client_id"`
	ClientName              string `json:"client_name"`
	GrantTypes              string `json:"grant_types"`
	ResponseTypes           string `json:"response_types"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method"`
	RedirectUris            string `json:"redirect_uris"`
}
