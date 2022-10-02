/*
Package config parses configuration files

Currently it parses a set of JSON files. The format is:

	  {
		"resource": "<name of plugin>",
		"label": "<label for resource>",
		....
	  }

The fields names should be strings, and the field values are
interpreted as follows:

	true, false: boolean
	123: integer
	123.456: float
	"10s": duration
	"hello": string
	["a", "b", "c"]: list
	{"a": 1, "b": 2}: map
*/
package config
