package handlers

import "net/http"

//AuthenticationIface abstracts the functionlities on authentication handler
type AuthenticationIface interface {
	GenerateToken(w http.ResponseWriter, r *http.Request)
	IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler
}
