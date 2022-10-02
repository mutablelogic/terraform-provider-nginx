package provider_test

import (
	"context"
	"testing"

	iface "github.com/mutablelogic/terraform-provider-nginx"
	"github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
)

type config1 struct {
	Label string
	t1    iface.Task
	t2    iface.Task
}

func (c config1) Name() string {
	return "config1"
}

func (c config1) New(ctx context.Context, p iface.Provider) (iface.Task, error) {
	return provider.Config{Label: c.Label}.New(ctx, p)
}

func Test_Deps_001(t *testing.T) {
	t1, err := provider.Config{Label: "t1"}.New(context.Background(), nil)
	if err != nil {
		t.Fatal("t1:", err)
	}
	t2, err := provider.Config{Label: "t2"}.New(context.Background(), nil)
	if err != nil {
		t.Fatal("t2:", err)
	}
	config := config1{
		Label: t.Name(),
		t1:    t1,
		t2:    t2,
	}
	deps, err := provider.ReorderPlugins(config)
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) == 0 {
		t.Fatal("Expected dependencies")
	}
}
