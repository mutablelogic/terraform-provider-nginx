package hcl

import (
	"reflect"
	"strings"
	"time"

	// Module imports
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	TagName = "hcl"
)

const (
	hclTagAttr     = "attr"
	hclTagBlock    = "block"
	hclTagLabel    = "label"
	hclTagOptional = "optional"
)

var (
	typeString     = reflect.TypeOf("")
	typeListString = reflect.TypeOf([]string{})
	typeDuration   = reflect.TypeOf(time.Second)
	typeTask       = reflect.TypeOf((*Task)(nil)).Elem()
	typeAny        = reflect.TypeOf((*interface{})(nil)).Elem()
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func SpecForType(t reflect.Type, tag, name string) (hcldec.Spec, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		return specForBlock(t, tag, name, 0, 0)
	default:
		return nil, ErrBadParameter.With("Invalid type:", t)
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func specForBlock(t reflect.Type, tag, name string, min, max int) (hcldec.Spec, error) {
	n := t.NumField()
	result := make([]hcldec.Spec, 0, n)
	label := 0
	for i := 0; i < n; i++ {
		field := t.Field(i)

		// Obtain name and kind of tag
		name, kind := tagNameKind(field.Tag, tag)
		if name == "" {
			continue
		}

		// Append to spec
		switch kind {
		case hclTagAttr:
			if spec, err := specForAttr(field.Type, name, true); err != nil {
				return nil, err
			} else {
				result = append(result, spec)
			}
		case hclTagOptional:
			if spec, err := specForAttr(field.Type, name, false); err != nil {
				return nil, err
			} else {
				result = append(result, spec)
			}
		case hclTagBlock:
			if field.Type.Kind() == reflect.Ptr {
				// Zero or one block
				if spec, err := specForBlock(field.Type.Elem(), tag, name, 0, 1); err != nil {
					return nil, err
				} else {
					result = append(result, spec)
				}
			} else {
				// Exactly one block
				if spec, err := specForBlock(field.Type.Elem(), tag, name, 1, 1); err != nil {
					return nil, err
				} else {
					result = append(result, spec)
				}
			}
		case hclTagLabel:
			spec, err := specForLabel(field.Type, name, label)
			if err != nil {
				return nil, err
			}
			result = append(result, spec)
			label++
		default:
			return nil, ErrBadParameter.Withf("invalid hcl field tag %q on %q", kind, field.Name)
		}
	}

	// If exactly one block, return a BlockSpec
	return &hcldec.BlockListSpec{
		TypeName: name,
		Nested:   hcldec.TupleSpec(result),
		MinItems: min,
		MaxItems: max,
	}, nil
}

func specForAttr(t reflect.Type, name string, required bool) (*hcldec.AttrSpec, error) {
	ty := typeForAttr(t)
	if ty == cty.NilType {
		return nil, ErrBadParameter.Withf("unsupported type %q for attribute %q", t, name)
	}

	// Return success
	return &hcldec.AttrSpec{
		Name:     name,
		Type:     ty,
		Required: required,
	}, nil
}

func specForLabel(t reflect.Type, name string, index int) (*hcldec.BlockLabelSpec, error) {
	ty := typeForAttr(t)
	if ty != cty.String {
		return nil, ErrBadParameter.Withf("unsupported type %q for label %q", t, name)
	}
	return &hcldec.BlockLabelSpec{
		Index: index,
		Name:  name,
	}, nil
}

func tagNameKind(field reflect.StructTag, name string) (string, string) {
	tag := field.Get(name)
	if tag == "" {
		return "", ""
	}
	comma := strings.Index(tag, ",")
	if comma != -1 {
		return tag[:comma], tag[comma+1:]
	} else {
		return tag, hclTagAttr
	}
}

func typeForAttr(t reflect.Type) cty.Type {
	switch t {
	case typeString:
		return cty.String
	case typeListString:
		return cty.List(cty.String)
	case typeDuration:
		return cty.DynamicPseudoType
	case typeTask:
		return cty.DynamicPseudoType
	case typeAny:
		return cty.DynamicPseudoType
	}
	// By default, return NilType for unsupported types
	return cty.NilType
}
