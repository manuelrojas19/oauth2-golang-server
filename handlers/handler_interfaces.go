package handlers

import "net/http"

type Handler interface {
	Handler(http.ResponseWriter, *http.Request)
}
