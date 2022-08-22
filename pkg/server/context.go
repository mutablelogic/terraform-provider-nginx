package server

import (
	"context"
	"net/http"
)

type contextType uint

const (
	contextNone contextType = iota
	contextParams
)

func ctxWithParams(ctx context.Context, params []string) context.Context {
	return context.WithValue(ctx, contextParams, params)
}

func ReqParams(req *http.Request) []string {
	if value, ok := req.Context().Value(contextParams).([]string); ok {
		return value
	} else {
		return nil
	}
}
