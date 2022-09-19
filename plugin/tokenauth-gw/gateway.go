package main

import (
	"fmt"

	// Module imports
	plugin "github.com/mutablelogic/terraform-provider-nginx/plugin"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type gateway struct {
	plugin.TokenAuth
	label string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(c Config, t plugin.TokenAuth) (Task, error) {
	this := new(gateway)
	this.label = c.Label
	this.TokenAuth = t

	// Return success
	return this, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (plugin *gateway) String() string {
	str := "<tokenauth-gateway"
	str += fmt.Sprintf(" label=%q", plugin.label)
	return str + ">"
}
