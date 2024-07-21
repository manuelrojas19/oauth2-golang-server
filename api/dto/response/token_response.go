package response

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewTokenResponse(accessToken string, tokenType string, expiresIn int) *TokenResponse {
	return &TokenResponse{AccessToken: accessToken, TokenType: tokenType, ExpiresIn: expiresIn}
}
