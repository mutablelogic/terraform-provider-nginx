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
