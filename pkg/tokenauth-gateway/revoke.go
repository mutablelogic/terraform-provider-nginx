package tokenauth_gateway

import (
	"net/http"

	// Modules
	context "github.com/mutablelogic/terraform-provider-nginx/pkg/context"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"
)

func (plugin *gateway) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	admin := context.ReqAdmin(r)
	params := context.ReqParams(r)
	if len(params) != 1 {
		util.ServeError(w, http.StatusBadRequest)
		return
	}
	if !admin {
		util.ServeError(w, http.StatusUnauthorized)
		return
	}

	name := params[0]
	if !plugin.Exists(name) {
		util.ServeError(w, http.StatusNotFound)
	} else if err := plugin.Revoke(name); err != nil {
		util.ServeError(w, http.StatusInternalServerError, err.Error())
	} else {
		// Serve emoty page
		util.ServeEmpty(w, http.StatusOK)
	}
}
