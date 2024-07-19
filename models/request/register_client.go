package request

type RegisterClientRequest struct {
	RedirectUris string `json:"redirect_uris"`
}

func NewRegisterClientRequest() *RegisterClientRequest {
	return &RegisterClientRequest{}
}
