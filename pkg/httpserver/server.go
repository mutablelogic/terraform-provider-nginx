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

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Addr    string        // Address or path for binding HTTP server
	TLS     *TLS          // TLS parameters
	Timeout time.Duration // Read timeout on HTTP requests
}

type TLS struct {
	Key  string // Path to TLS Private Key
	Cert string // Path to TLS Certificate
}

type server struct {
	srv  *http.Server
	fcgi *fcgi.Server
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultTimeout = 10 * time.Second
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create the module
func (cfg Config) New(handler http.Handler) (*server, error) {
	this := new(server)

	if handler == nil {
		handler = http.DefaultServeMux
	}

	// Check addr for being (host, port). If not, then run as FCGI server
	if _, _, err := net.SplitHostPort(cfg.Addr); cfg.Addr != "" && err != nil {
		if err := this.fcgiserver(cfg.Addr, handler); err != nil {
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
	if err := this.netserver(cfg.Addr, tlsconfig, cfg.Timeout, handler); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *server) Run(ctx context.Context, _ Kernel) error {
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
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (*server) C() <-chan Event {
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *server) fcgiserver(path string, handler http.Handler) error {
	// Create server
	this.fcgi = &fcgi.Server{}
	this.fcgi.Network = "unix"
	this.fcgi.Addr = path
	this.fcgi.Handler = handler

	// Return success
	return nil
}

func (this *server) netserver(addr string, config *tls.Config, timeout time.Duration, handler http.Handler) error {
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

func (this *server) runInForeground() error {
	if this.fcgi != nil {
		return this.fcgi.ListenAndServe()
	} else if this.srv.TLSConfig != nil {
		return this.srv.ListenAndServeTLS("", "")
	} else {
		return this.srv.ListenAndServe()
	}
}

func (this *server) stop() error {
	if this.fcgi != nil {
		return this.fcgi.Close()
	} else {
		return this.srv.Close()
	}
}
