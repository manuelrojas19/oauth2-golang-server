package commands

type Authorization struct {
	ClientId     string
	Scope        string
	RedirectUri  string
	ResponseType string
	SessionId    string
}
