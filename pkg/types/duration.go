package types

import (
	"strconv"
	"time"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Duration time.Duration

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (d *Duration) UnmarshalJSON(data []byte) error {
	if v, err := strconv.Unquote(string(data)); err != nil {
		return err
	} else if v_, err := time.ParseDuration(v); err != nil {
		return err
	} else {
		*d = Duration(v_)
		return nil
	}
}
