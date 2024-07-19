package handlers

import "net/http"

type RegisterHandler interface {
	Register(http.ResponseWriter, *http.Request)
}
