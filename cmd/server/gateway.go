package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sync"

	// Moudule imports
	config "github.com/mutablelogic/terraform-provider-nginx/pkg/config"
	server "github.com/mutablelogic/terraform-provider-nginx/pkg/server"
)

type Runner interface {
	Run(context.Context) error
	List() []config.Object
}

type Gateway struct {
	sync.Mutex
	Runner
	http.Handler
}

var (
	reList   = regexp.MustCompile(`^/$`)
	reConfig = regexp.MustCompile(`^/([a-zA-Z0-9_\-]+)$`)
)

func NewGateway(available, enabled string) (*Gateway, error) {
	this := new(Gateway)

	// Create a runner
	runner, err := config.Config{AvailablePath: available, EnabledPath: enabled}.NewRunner()
	if err != nil {
		return nil, err
	} else {
		this.Runner = runner
	}

	// Create a router
	router := server.NewRouter()
	router.AddRoute(reList, this.APIList)
	router.AddRoute(reConfig, this.APIConfig, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch)
	this.Handler = router

	// Report runner events
	go func() {
		for evt := range runner.C() {
			fmt.Println("evt=", evt)
		}
	}()

	return this, nil
}

func (g *Gateway) APIList(w http.ResponseWriter, req *http.Request) {
	g.Lock()
	defer g.Unlock()
	server.ServeJSON(w, g.Runner.List(), http.StatusOK, 2)
}

func (g *Gateway) APIConfig(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Operate", req.Method, server.ReqParams(req))
}
