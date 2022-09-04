package context

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type contextType uint

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	contextNone contextType = iota
	contextName
	contextLabel
	contextPrefix
	contextParams
)

///////////////////////////////////////////////////////////////////////////////
// CREATE CONTEXT

func WithPrefixParams(ctx context.Context, prefix string, params []string) context.Context {
	return context.WithValue(context.WithValue(ctx, contextParams, params), contextPrefix, prefix)
}

func WithName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, contextName, name)
}

func WithNameLabel(ctx context.Context, name, label string) context.Context {
	return context.WithValue(context.WithValue(ctx, contextName, name), contextLabel, label)
}

///////////////////////////////////////////////////////////////////////////////
// RETURN VALUES FROM CONTEXT

func Name(ctx context.Context) string {
	return contextString(ctx, contextName)
}

func Label(ctx context.Context) string {
	return contextString(ctx, contextLabel)
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

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func DumpContext(ctx context.Context, w io.Writer) {
	fmt.Fprintf(w, "<context")
	if value, ok := ctx.Value(contextName).(string); ok {
		fmt.Fprintf(w, " name=%q", value)
	}
	if value, ok := ctx.Value(contextLabel).(string); ok {
		fmt.Fprintf(w, " label=%q", value)
	}
	if value, ok := ctx.Value(contextPrefix).(string); ok {
		fmt.Fprintf(w, " prefix=%q", value)
	}
	if value, ok := ctx.Value(contextParams).([]string); ok {
		fmt.Fprintf(w, " params=%q", value)
	}
	fmt.Fprintf(w, ">")
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func contextString(ctx context.Context, key contextType) string {
	if value, ok := ctx.Value(key).(string); ok {
		return value
	} else {
		return ""
	}
}
