package authgw

import (
	"context"
	"net/http"
)

type contextType uint

const (
	contextNone contextType = iota
	contextToken
)

func ctxWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, contextToken, token)
}

func ReqToken(req *http.Request) string {
	if value, ok := req.Context().Value(contextToken).(string); ok {
		return value
	} else {
		return ""
	}
}
