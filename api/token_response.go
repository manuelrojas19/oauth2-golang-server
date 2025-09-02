package api

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func NewTokenResponse(accessToken string, tokenType string, accessTokenExpiresIn int, refreshToken string, scope string) *TokenResponse {
	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		ExpiresIn:    accessTokenExpiresIn,
		RefreshToken: refreshToken,
		Scope:        scope,
	}
}
