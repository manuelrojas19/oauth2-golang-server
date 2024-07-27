package commands

type Authorize struct {
	ClientId     string
	Scope        string
	RedirectUri  string
	ResponseType string
	SessionId    string
}
