package tokenauth_gateway

import (
	"fmt"
	"net/http"
)

func (plugin *gateway) AuthenticateHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("TODO: AuthenticateHandler")
		fn(w, r)
	}
}

func (plugin *gateway) AuthenticateAdminHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("TODO: AuthenticateAdminHandler")
		fn(w, r)
	}
}
