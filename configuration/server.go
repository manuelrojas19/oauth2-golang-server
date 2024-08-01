package configuration

import "net/http"

func NewServeMux() *http.ServeMux {
	return http.NewServeMux()
}
