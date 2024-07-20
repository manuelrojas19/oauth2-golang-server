package authmethodtype

type TokenEndpointAuthMethod string

const (
	ClientSecretBasic TokenEndpointAuthMethod = "client_secret_basic"
	ClientSecretPost  TokenEndpointAuthMethod = "client_secret_post"
	None              TokenEndpointAuthMethod = "none"
)
