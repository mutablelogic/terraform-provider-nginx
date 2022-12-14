package context

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
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
	contextAdmin
	contextAddress
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func ContextForSignal(signals ...os.Signal) context.Context {
	if len(signals) == 0 {
		return nil
	}

	ch := make(chan os.Signal)
	ctx, cancel := context.WithCancel(context.Background())

	// Send message on channel when signal received
	signal.Notify(ch, signals...)

	// When any signal received, call cancel
	go func() {
		<-ch
		cancel()
	}()

	// Return success
	return ctx
}

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

func WithAdmin(ctx context.Context, admin bool) context.Context {
	return context.WithValue(ctx, contextAdmin, admin)
}

func WithAddress(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, contextAddress, addr)
}

///////////////////////////////////////////////////////////////////////////////
// RETURN VALUES FROM CONTEXT

func Name(ctx context.Context) string {
	return contextString(ctx, contextName)
}

func Label(ctx context.Context) string {
	return contextString(ctx, contextLabel)
}

func Admin(ctx context.Context) bool {
	return contextBool(ctx, contextAdmin)
}

func Address(ctx context.Context) string {
	return contextString(ctx, contextAddress)
}

func ReqParams(req *http.Request) []string {
	if value, ok := req.Context().Value(contextParams).([]string); ok {
		return value
	} else {
		return nil
	}
}

func ReqPrefix(req *http.Request) string {
	return contextString(req.Context(), contextPrefix)
}

func ReqName(req *http.Request) string {
	return contextString(req.Context(), contextName)
}

func ReqAdmin(req *http.Request) bool {
	return contextBool(req.Context(), contextAdmin)
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
	if value, ok := ctx.Value(contextBool).(bool); ok {
		fmt.Fprintf(w, " admin=%v", value)
	}
	if value, ok := ctx.Value(contextAddress).(string); ok {
		fmt.Fprintf(w, " address=%q", value)
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

func contextBool(ctx context.Context, key contextType) bool {
	if value, ok := ctx.Value(key).(bool); ok {
		return value
	} else {
		return false
	}
}
