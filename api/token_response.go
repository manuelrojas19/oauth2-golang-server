package api

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func NewTokenResponse(accessToken string, tokenType string, refreshToken string) *TokenResponse {
	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		RefreshToken: refreshToken,
	}
}
