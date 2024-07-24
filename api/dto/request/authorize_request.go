package request

// AuthorizeRequest represents the request to authorize a client
type AuthorizeRequest struct {
	ResponseType string
	ClientId     string
	RedirectUri  string
	Scope        string
	State        string
}
