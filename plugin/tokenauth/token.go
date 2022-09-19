package main

import (
	"fmt"
	"time"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type token struct {
	Token string    `json:"token"`
	Time  time.Time `json:"atime"`
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t *token) String() string {
	str := "<tokenauth-token"
	str += fmt.Sprintf(" token=%q", t.Token)
	str += fmt.Sprintf(" last_accessed=%q", t.Time.Format(time.RFC3339))
	return str + ">"
}
