package httpserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type contextType uint

const (
	contextNone contextType = iota
	contextPrefix
	contextParams
)

func ctxWithPrefixParams(ctx context.Context, prefix string, params []string) context.Context {
	return context.WithValue(context.WithValue(ctx, contextParams, params), contextPrefix, prefix)
}

func ReqParams(req *http.Request) []string {
	if value, ok := req.Context().Value(contextParams).([]string); ok {
		return value
	} else {
		return nil
	}
}

func ReqPrefix(req *http.Request) string {
	if value, ok := req.Context().Value(contextPrefix).(string); ok {
		return value
	} else {
		return ""
	}
}

func DumpContext(ctx context.Context, w io.Writer) {
	fmt.Fprintf(w, "<context")
	if value, ok := ctx.Value(contextPrefix).(string); ok {
		fmt.Fprintf(w, " prefix=%q", value)
	}
	if value, ok := ctx.Value(contextParams).([]string); ok {
		fmt.Fprintf(w, " params=%q", value)
	}
	fmt.Fprintf(w, ">")
}
