package tokenauth_gateway

import (
	"net/http"

	// Modules
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"
)

func (plugin *gateway) CreateHandler(w http.ResponseWriter, r *http.Request) {
	// Enumerate tokens
	tokens := plugin.Enumerate()
	if tokens == nil {
		util.ServeError(w, http.StatusInternalServerError)
		return
	}

	// Create response
	result := make([]Token, 0, len(tokens))
	for name, atime := range tokens {
		result = append(result, Token{Name: name, AccessTime: atime})
	}

	// Serve response
	util.ServeJSON(w, result, http.StatusOK, 2)
}