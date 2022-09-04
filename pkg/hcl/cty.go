package hcl

import (
	"math"
	"math/big"
	"reflect"

	// Modules
	"github.com/hashicorp/go-multierror"
	"github.com/zclconf/go-cty/cty"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// FromCtyValue decodes values from a cty.Value into golang native value
func FromCtyValue(val cty.Value, target interface{}) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return ErrBadParameter.With("target must be a pointer")
	} else if v.IsNil() {
		return ErrBadParameter.With("target is nil")
	}

	return fromCtyValue(val, v)
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func fromCtyValue(val cty.Value, target reflect.Value) error {
	t := val.Type()

	// Set target to element
	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}
	if !target.CanSet() {
		if target.IsValid() {
			return ErrBadParameter.With("cannot set target of type", target.Type())
		} else {
			return ErrBadParameter.With("cannot set target")
		}
	}

	// Where value is nil, return zero-value for type
	if val.IsNull() {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	switch t {
	case cty.Bool:
		return fromCtyBool(val, target)
	case cty.Number:
		return fromCtyNumber(val, target)
	case cty.String:
		return fromCtyString(val, target)
	}
	switch {
	case t.IsListType() || t.IsSetType():
		return fromCtyList(val, target)
	case t.IsMapType():
		return fromCtyMap(val, target)
	case t.IsObjectType():
		return fromCtyObject(val, target)
	case t.IsTupleType():
		return fromCtyTuple(val, target)
	default:
		return ErrBadParameter.With("unsupported source type: ", t.GoString())
	}
}

func fromCtyBool(val cty.Value, target reflect.Value) error {
	switch {
	case target.Kind() == reflect.Bool:
		target.SetBool(val.True())
	case target.Type() == typeAny:
		target.Set(reflect.ValueOf(val.True()))
	default:
		return ErrBadParameter.With("bool type is required")
	}

	// Return success
	return nil
}

func fromCtyNumber(val cty.Value, target reflect.Value) error {
	switch target.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fromCtyNumberInt(val, target)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fromCtyNumberUInt(val, target)
	case reflect.Float32, reflect.Float64:
		return fromCtyNumberFloat(val, target)
	default:
		return ErrBadParameter.Withf("number type is required, not %v", target.Kind())
	}
}

func fromCtyNumberInt(val cty.Value, target reflect.Value) error {
	var min, max int64
	switch target.Type().Bits() {
	case 8:
		min = math.MinInt8
		max = math.MaxInt8
	case 16:
		min = math.MinInt16
		max = math.MaxInt16
	case 32:
		min = math.MinInt32
		max = math.MaxInt32
	case 64:
		min = math.MinInt64
		max = math.MaxInt64
	}

	iv, accuracy := val.AsBigFloat().Int64()
	if accuracy != big.Exact || iv < min || iv > max {
		return ErrBadParameter.Withf("integer type is required, not %v", target.Kind())
	} else {
		target.SetInt(iv)
	}

	// Return success
	return nil
}

func fromCtyNumberUInt(val cty.Value, target reflect.Value) error {
	var max uint64
	switch target.Type().Bits() {
	case 8:
		max = math.MaxUint8
	case 16:
		max = math.MaxUint16
	case 32:
		max = math.MaxUint32
	case 64:
		max = math.MaxUint64
	}

	iv, accuracy := val.AsBigFloat().Uint64()
	if accuracy != big.Exact || iv > max {
		return ErrBadParameter.Withf("unsigned integer type is required, not %v", target.Kind())
	} else {
		target.SetUint(iv)
	}

	// Return success
	return nil
}

func fromCtyNumberFloat(val cty.Value, target reflect.Value) error {
	switch {
	case target.Kind() == reflect.Float32 || target.Kind() == reflect.Float64:
		fv, _ := val.AsBigFloat().Float64()
		target.SetFloat(fv)
	case target.Type() == typeAny:
		fv, _ := val.AsBigFloat().Float64()
		target.Set(reflect.ValueOf(fv))
	default:
		return ErrBadParameter.Withf("float type is required, not %v", target.Kind())
	}

	// Return success
	return nil
}

func fromCtyString(val cty.Value, target reflect.Value) error {
	switch {
	case target.Kind() == reflect.String:
		target.SetString(val.AsString())
	case target.Type() == typeAny:
		target.Set(reflect.ValueOf(val.AsString()))
	default:
		return ErrBadParameter.Withf("string type is required, not %q", target.Kind())
	}

	// Return success
	return nil
}

func fromCtyList(val cty.Value, target reflect.Value) error {
	var result error
	var i int

	switch target.Kind() {
	case reflect.Slice:
		length := val.LengthInt()
		tv := reflect.MakeSlice(target.Type(), length, length)
		if !val.ForEachElement(func(key cty.Value, val cty.Value) bool {
			if err := fromCtyValue(val, tv.Index(i)); err != nil {
				result = multierror.Append(result, err)
				return true
			}
			i++
			return false
		}) {
			target.Set(tv)
		}
	case reflect.Array:
		length := val.LengthInt()
		if length != target.Len() {
			result = multierror.Append(result, ErrBadParameter.Withf("must be a list of length %d", target.Len()))
		} else {
			val.ForEachElement(func(key cty.Value, val cty.Value) bool {
				if err := fromCtyValue(val, target.Index(i)); err != nil {
					result = multierror.Append(result, err)
					return true
				}
				i++
				return false
			})
		}
	default:
		result = multierror.Append(result, ErrBadParameter.Withf("list type is required, not %v", target.Kind()))
	}

	// Return success
	return result
}

func fromCtyMap(val cty.Value, target reflect.Value) error {
	var result error

	if target.Kind() != reflect.Map {
		return ErrBadParameter.Withf("map type is required, not %v", target.Kind())
	}

	tv := reflect.MakeMap(target.Type())
	et := target.Type().Elem()
	if !val.ForEachElement(func(key, val cty.Value) bool {

		ks := key.AsString()
		targetElem := reflect.New(et)
		if err := fromCtyValue(val, targetElem); err != nil {
			result = multierror.Append(result, err)
			return true
		} else {
			tv.SetMapIndex(reflect.ValueOf(ks), targetElem.Elem())
			return false
		}
	}) {
		target.Set(tv)
	}

	// Return success
	return result
}

func fromCtyObject(val cty.Value, target reflect.Value) error {
	var result error

	if target.Kind() != reflect.Struct {
		return ErrBadParameter.Withf("struct type is required, not %v", target.Kind())
	}

	// Make zero-version of the target struct
	target.Set(reflect.New(target.Type()).Elem())

	// Set attributes
	for key := range val.Type().AttributeTypes() {
		field := target.FieldByName(key)
		if !field.IsValid() {
			result = multierror.Append(result, ErrBadParameter.Withf("unsupported attribute %q", key))
		} else if err := fromCtyValue(val.GetAttr(key), field); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return success
	return result
}

func fromCtyTuple(val cty.Value, target reflect.Value) error {
	var result error

	if target.Kind() != reflect.Struct {
		return ErrBadParameter.Withf("struct type is required, not %v", target.Kind())
	}

	elemTypes := val.Type().TupleElementTypes()
	fieldCount := target.Type().NumField()
	if fieldCount != len(elemTypes) {
		return ErrBadParameter.Withf("a tuple of %d elements is required", fieldCount)
	}

	for i := range elemTypes {
		ev := val.Index(cty.NumberIntVal(int64(i)))
		targetField := target.Field(i)
		if err := fromCtyValue(ev, targetField); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return success
	return result
}
