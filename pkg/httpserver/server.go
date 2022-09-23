package httpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	// Modules
	multierror "github.com/hashicorp/go-multierror"
	fcgi "github.com/mutablelogic/terraform-provider-nginx/pkg/fcgi"
	"github.com/mutablelogic/terraform-provider-nginx/pkg/provider"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

type httpserver struct {
	Router
	label string
	srv   *http.Server
	fcgi  *fcgi.Server
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create the server
func NewWithConfig(c Config) (*httpserver, error) {
	this := new(httpserver)

	// Set label
	if cfg.Label == "" {
		this.label = cfg.Name()
	} else {
		this.label = cfg.Label
	}

	// Obtain router
	if cfg.Router == nil {
		fmt.Println("Creating a new router since doesn't exist", cfg)
		if router, err := provider.New(ctx, RouterConfig{
			Label: this.label + "-router",
		}); err != nil {
			return nil, err
		} else {
			cfg.Router = router
		}
	}

	// Check that router is a handler and a router
	if _, ok := cfg.Router.(http.Handler); !ok {
		return nil, ErrInternalAppError.With("invalid router")
	} else if _, ok := cfg.Router.(Router); !ok {
		return nil, ErrInternalAppError.With("invalid router")
	} else {
		this.Router = cfg.Router.(Router)
	}

	// Check addr for being (host, port). If not, then run as FCGI server
	if _, _, err := net.SplitHostPort(cfg.Addr); cfg.Addr != "" && err != nil {
		if err := this.fcgiserver(cfg.Addr, cfg.Router.(http.Handler)); err != nil {
			return nil, err
		} else {
			return this, nil
		}
	}

	// If either key or cert is non-nil then create a TLSConfig
	var tlsconfig *tls.Config
	if cfg.TLS != nil {
		if cert, err := tls.LoadX509KeyPair(cfg.TLS.Cert, cfg.TLS.Key); err != nil {
			return nil, err
		} else {
			tlsconfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
		}
	}

	// If addr is empty, then set depending on whether it's SSL or not
	if cfg.Addr == "" {
		if tlsconfig == nil {
			cfg.Addr = ":http"
		} else {
			cfg.Addr = ":https"
		}
	}

	// Create net server
	if err := this.netserver(cfg.Addr, tlsconfig, cfg.Timeout, cfg.Router.(http.Handler)); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *server) Run(ctx context.Context) error {
	var result error
	go func() {
		<-ctx.Done()
		if err := this.stop(); err != nil {
			result = multierror.Append(result, err)
		}
	}()
	if err := this.runInForeground(); err != nil && errors.Is(err, http.ErrServerClosed) == false {
		result = multierror.Append(result, err)
	}
	return result
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *server) String() string {
	str := "<httpserver"
	if this.label != "" {
		str += fmt.Sprintf(" label=%q", this.label)
	}
	if this.fcgi != nil {
		str += fmt.Sprintf(" fcgi=%q", this.fcgi.Addr)
	} else {
		str += fmt.Sprintf(" addr=%q", this.srv.Addr)
		if this.srv.TLSConfig != nil {
			str += " tls=true"
		}
		if this.srv.ReadHeaderTimeout != 0 {
			str += fmt.Sprintf(" read_timeout=%v", this.srv.ReadHeaderTimeout)
		}
	}
	if this.Router != nil {
		str += fmt.Sprintf(" router=%v", this.Router)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *httpserver) fcgiserver(path string, handler http.Handler) error {
	// Create server
	this.fcgi = &fcgi.Server{}
	this.fcgi.Network = "unix"
	this.fcgi.Addr = path
	this.fcgi.Handler = handler

	// Return success
	return nil
}

func (this *httpserver) netserver(addr string, config *tls.Config, timeout time.Duration, handler http.Handler) error {
	// Set up server
	this.srv = &http.Server{}
	if config != nil {
		this.srv.TLSConfig = config
	}

	// Set default timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	// Set server parameters
	this.srv.Addr = addr
	this.srv.Handler = handler
	this.srv.ReadHeaderTimeout = timeout
	this.srv.IdleTimeout = timeout

	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// START AND STOP

func (this *httpserver) runInForeground() error {
	if this.fcgi != nil {
		return this.fcgi.ListenAndServe()
	} else if this.srv.TLSConfig != nil {
		return this.srv.ListenAndServeTLS("", "")
	} else {
		return this.srv.ListenAndServe()
	}
}

func (this *httpserver) stop() error {
	if this.fcgi != nil {
		return this.fcgi.Close()
	} else {
		return this.srv.Close()
	}
}
