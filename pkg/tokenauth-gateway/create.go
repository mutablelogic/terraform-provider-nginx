package tokenauth_gateway

import (
	"net/http"

	// Modules
	context "github.com/mutablelogic/terraform-provider-nginx/pkg/context"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"
)

func (plugin *gateway) CreateHandler(w http.ResponseWriter, r *http.Request) {
	params := context.ReqParams(r)
	if len(params) != 1 {
		util.ServeError(w, http.StatusBadRequest)
		return
	}

	name := params[0]
	if plugin.Exists(name) {
		util.ServeError(w, http.StatusBadRequest)
	} else if value, err := plugin.Create(name); err != nil {
		util.ServeError(w, http.StatusInternalServerError, err.Error())
	} else {
		util.ServeJSON(w, value, http.StatusCreated, 2)
	}
}
