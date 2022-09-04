package hcl

import (
	"context"
	"fmt"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type varblock struct {
	Label       string `hcl:"label,label"`
	Type        string `hcl:"type"`
	Default     any    `hcl:"default,optional"`
	Description string `hcl:"description,optional"`
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (varblock) Name() string {
	return "var"
}

func (v varblock) New(context.Context, Provider) (Task, error) {
	return nil, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v *varblock) String() string {
	str := "<var"
	if v.Label != "" {
		str += fmt.Sprintf(" label=%q", v.Label)
	}
	if v.Type != "" {
		str += fmt.Sprintf(" type=%q", v.Type)
	}
	switch v.Default.(type) {
	case string:
		str += fmt.Sprintf(" default=%q", v.Default)
	default:
		str += fmt.Sprintf(" default=%v", v.Default)
	}
	if v.Description != "" {
		str += fmt.Sprintf(" description=%q", v.Description)
	}
	return str + ">"
}
