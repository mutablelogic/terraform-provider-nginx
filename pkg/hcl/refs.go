package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type refs struct {
	references map[string][]string
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewRefs() *refs {
	r := new(refs)
	r.references = make(map[string][]string)
	return r
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (r *refs) String() string {
	str := "<refs"
	for k, v := range r.references {
		str += fmt.Sprintf(" %v=%q", k, v)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (r *refs) CreateReferences(body hcl.Body, spec hcldec.Spec) error {
	varName := varblock{}.Name() // "var"
	traversal := hcldec.Variables(body, spec)
	for _, t := range traversal {
		ts := t.SimpleSplit()
		refName := ts.RootName()
		if len(ts.Rel) != 1 {
			return ErrBadParameter.Withf("invalid reference: %q", t.SourceRange())
		}
		refIdentifier := ts.Rel[0].(hcl.TraverseAttr).Name

		if refName == varName {
			// 'var' references are stored at the top level
			r.references[refIdentifier] = nil
		} else {
			r.references[refName] = append(r.references[refName], refIdentifier)
		}
	}
	return nil
}

func (r *refs) Context() *hcl.EvalContext {
	ctx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
	}
	for k := range r.references {
		ctx.Variables[k] = r.objectValForName(k)
	}
	return ctx
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (r *refs) objectValForName(name string) cty.Value {
	values := r.references[name]
	if values == nil {
		return cty.DynamicVal.Mark(name)
	}
	result := make(map[string]cty.Value)
	for _, label := range values {
		result[label] = cty.DynamicVal.Mark(label)
	}
	return cty.ObjectVal(result)
}
