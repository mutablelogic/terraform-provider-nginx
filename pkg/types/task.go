package types

import (
	// Modules

	"strconv"

	iface "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Task struct {
	iface.Task
	Ref string // Reference to resolve into a task
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (t *Task) UnmarshalJSON(data []byte) error {
	if v, err := strconv.Unquote(string(data)); err != nil {
		return err
	} else {
		t.Ref = v
	}
	return nil
}
