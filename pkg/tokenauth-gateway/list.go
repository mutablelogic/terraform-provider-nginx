package tokenauth_gateway

import (
	"net/http"

	// Modules
	tokenauth "github.com/mutablelogic/terraform-provider-nginx/pkg/tokenauth"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"
)

type Token struct {
	tokenauth.Token
	Name string `json:"name"`
}

func (plugin *gateway) ListHandler(w http.ResponseWriter, r *http.Request) {
	// Enumerate tokens
	tokens := plugin.Enumerate()
	if tokens == nil {
		util.ServeError(w, http.StatusInternalServerError)
		return
	}

	// Create response
	result := make([]Token, 0, len(tokens))
	for name, time := range tokens {
		result = append(result, Token{Name: name, Token: tokenauth.Token{Time: time}})
	}

	// Serve response
	util.ServeJSON(w, result, http.StatusOK, 2)
}
