package authgw

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	// Modules
	httpserver "github.com/mutablelogic/terraform-provider-nginx/pkg/httpserver"
	tokenauth "github.com/mutablelogic/terraform-provider-nginx/pkg/tokenauth"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type TokenAuth interface {
	Exists(name string) bool
	Create(name string) (string, error)
	Revoke(name string) error
	Enumerate() map[string]time.Time
	Matches(value string) string
}

type authgw struct {
	TokenAuth
	prefix string
}

type token struct {
	Token      string    `json:"token"`
	Admin      bool      `json:"admin,omitempty"`
	AccessTime time.Time `json:"access_time,omitempty"`
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	reEnumerateTokens   = regexp.MustCompile(`^/$`)
	reCreateRevokeToken = regexp.MustCompile(`^/([a-zA-Z][a-zA-Z0-9_\-]+)$`)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func New(auth TokenAuth, prefix string) *authgw {
	return &authgw{auth, prefix}
}

func (a *authgw) Run(ctx context.Context, kernel Kernel) error {
	// Add middleware
	if err := kernel.AddMiddleware("token-auth", a.AuthenticateRequest); err != nil {
		return err
	}
	if err := kernel.AddMiddleware("token-admin-auth", a.AuthenticateAdminRequest); err != nil {
		return err
	}
	// Add routes
	if err := kernel.AddHandler(a.prefix, reEnumerateTokens, a.EnumerateTokens); err != nil {
		return err
	}
	if err := kernel.AddHandler(a.prefix, reCreateRevokeToken, a.CreateRevokeToken, http.MethodDelete, http.MethodPost); err != nil {
		return err
	}

	// Wait until cancel
	<-ctx.Done()
	return ctx.Err()
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (*authgw) C() <-chan Event {
	return nil
}

/////////////////////////////////////////////////////////////////////
// HTTP HANDLERS

// Handler to enumerate existing tokens
func (a *authgw) EnumerateTokens(w http.ResponseWriter, req *http.Request) {
	var result []token
	for name, accessTime := range a.Enumerate() {
		result = append(result, token{name, name == tokenauth.AdminToken, accessTime})
	}
	httpserver.ServeJSON(w, result, http.StatusOK, 2)
}

// Handler to create or revoke a token
func (a *authgw) CreateRevokeToken(w http.ResponseWriter, req *http.Request) {
	// Get the token to create or revoke
	token := httpserver.ReqParams(req)
	if len(token) != 1 {
		httpserver.ServeError(w, http.StatusBadRequest)
		return
	}
	switch req.Method {
	case http.MethodDelete:
		if err := a.Revoke(token[0]); err != nil {
			httpserver.ServeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		if value, err := a.Create(token[0]); err != nil {
			httpserver.ServeError(w, http.StatusInternalServerError, err.Error())
			return
		} else {
			httpserver.ServeText(w, value, http.StatusOK)
		}
	default:
		httpserver.ServeError(w, http.StatusMethodNotAllowed)
	}
}

/////////////////////////////////////////////////////////////////////
// HTTP MIDDLEWARE

// Middleware to authenticate a request
func (a *authgw) AuthenticateRequest(child http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Token ") {
			httpserver.ServeError(w, http.StatusUnauthorized)
			return
		}
		if token := a.Matches(strings.TrimPrefix(auth, "Token ")); token == "" {
			httpserver.ServeError(w, http.StatusUnauthorized)
			return
		} else {
			req := req.WithContext(ctxWithToken(req.Context(), token))
			child(w, req)
		}
	}
}

// Middleware to authenticate an admin request
func (a *authgw) AuthenticateAdminRequest(child http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if token := ReqToken(req); token != tokenauth.AdminToken {
			httpserver.ServeError(w, http.StatusUnauthorized)
			return
		} else {
			child(w, req)
		}
	}
}
