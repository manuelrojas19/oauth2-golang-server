package errors

const (
	ErrUserNotAuthenticated    = "user not authenticated"
	ErrConsentRequired         = "user consent required"
	ErrUnsupportedResponseType = "the authorization server does not support obtaining an authorization code using this method"
	ErrInvalidRedirectUri      = "redirect URI is not registered for client"
)
