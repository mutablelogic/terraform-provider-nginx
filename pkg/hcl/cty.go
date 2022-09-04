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

	return fromCtyValue(val, v, make(cty.Path, 0))
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func fromCtyValue(val cty.Value, target reflect.Value, path cty.Path) error {
	t := val.Type()

	// Set target to element
	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}
	if !target.CanSet() {
		return ErrBadParameter.With("cannot set target of type", target.Type())
	}

	// Where value is nil, return zero-value for type
	if val.IsNull() {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	switch t {
	case cty.Bool:
		return fromCtyBool(val, target, path)
	case cty.Number:
		return fromCtyNumber(val, target, path)
	case cty.String:
		return fromCtyString(val, target, path)
	}
	switch {
	case t.IsListType() || t.IsSetType():
		return fromCtyList(val, target, path)
	case t.IsMapType():
		return fromCtyMap(val, target, path)
	case t.IsObjectType():
		return fromCtyObject(val, target, path)
		/*
			case t.IsTupleType():
				return fromCtyTuple(val, target, path)
			case t.IsCapsuleType():
				return fromCtyCapsule(val, target, path)
		*/
	}

	return ErrBadParameter.With("unsupported source type: ", t.GoString())
}

func fromCtyBool(val cty.Value, target reflect.Value, path cty.Path) error {
	switch {
	case target.Kind() == reflect.Bool:
		target.SetBool(val.True())
	case target.Type() == typeAny:
		target.Set(reflect.ValueOf(val.True()))
	default:
		return path.NewErrorf("bool type is required")
	}
	return nil
}

func fromCtyNumber(val cty.Value, target reflect.Value, path cty.Path) error {
	switch target.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fromCtyNumberInt(val, target, path)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fromCtyNumberUInt(val, target, path)
	case reflect.Float32, reflect.Float64:
		return fromCtyNumberFloat(val, target, path)
	default:
		return path.NewErrorf("number type is required, not %v", target.Kind())
	}
}

func fromCtyNumberInt(val cty.Value, target reflect.Value, path cty.Path) error {
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
		return path.NewErrorf("integer type is required, not %v", target.Kind())
	}

	target.SetInt(iv)
	return nil
}

func fromCtyNumberUInt(val cty.Value, target reflect.Value, path cty.Path) error {
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
		return path.NewErrorf("unsigned integer type is required, not %v", target.Kind())
	}

	target.SetUint(iv)
	return nil
}

func fromCtyNumberFloat(val cty.Value, target reflect.Value, path cty.Path) error {
	switch {
	case target.Kind() == reflect.Float32 || target.Kind() == reflect.Float64:
		fv, _ := val.AsBigFloat().Float64()
		target.SetFloat(fv)
	case target.Type() == typeAny:
		fv, _ := val.AsBigFloat().Float64()
		target.Set(reflect.ValueOf(fv))
	default:
		return path.NewErrorf("float type is required, not %v", target.Kind())
	}
	return nil
}

func fromCtyString(val cty.Value, target reflect.Value, path cty.Path) error {
	switch {
	case target.Kind() == reflect.String:
		target.SetString(val.AsString())
	case target.Type() == typeAny:
		target.Set(reflect.ValueOf(val.AsString()))
	default:
		return path.NewErrorf("string type is required, not %q", target.Kind())
	}

	// Return success
	return nil
}

func fromCtyList(val cty.Value, target reflect.Value, path cty.Path) error {
	var result error
	path = append(path, nil)
	i := 0

	switch target.Kind() {
	case reflect.Slice:
		length := val.LengthInt()
		tv := reflect.MakeSlice(target.Type(), length, length)
		if !val.ForEachElement(func(key cty.Value, val cty.Value) bool {
			path[len(path)-1] = cty.IndexStep{Key: key}
			if err := fromCtyValue(val, tv.Index(i), path); err != nil {
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
			result = multierror.Append(result, path.NewErrorf("must be a list of length %d", target.Len()))
		} else {
			val.ForEachElement(func(key cty.Value, val cty.Value) bool {
				path[len(path)-1] = cty.IndexStep{Key: key}
				if err := fromCtyValue(val, target.Index(i), path); err != nil {
					result = multierror.Append(result, err)
					return true
				}
				i++
				return false
			})
		}
	default:
		result = multierror.Append(result, path.NewErrorf("list type is required, not %v", target.Kind()))
	}

	path = path[:len(path)-1]
	return result
}

func fromCtyMap(val cty.Value, target reflect.Value, path cty.Path) error {
	var result error
	path = append(path, nil)

	if target.Kind() != reflect.Map {
		return path.NewErrorf("map type is required, not %v", target.Kind())
	}

	tv := reflect.MakeMap(target.Type())
	et := target.Type().Elem()
	path = append(path, nil)
	if !val.ForEachElement(func(key, val cty.Value) bool {
		path[len(path)-1] = cty.IndexStep{Key: key}

		ks := key.AsString()
		targetElem := reflect.New(et)
		if err := fromCtyValue(val, targetElem, path); err != nil {
			result = multierror.Append(result, err)
			return true
		} else {
			tv.SetMapIndex(reflect.ValueOf(ks), targetElem.Elem())
			return false
		}
	}) {
		target.Set(tv)
	}

	path = path[:len(path)-1]
	return result
}

func fromCtyObject(val cty.Value, target reflect.Value, path cty.Path) error {
	var result error
	path = append(path, nil)

	if target.Kind() != reflect.Struct {
		return path.NewErrorf("struct type is required, not %v", target.Kind())
	}

	// Make zero-version of the target struct
	target.Set(reflect.New(target.Type()).Elem())

	// Set attributes
	for key := range val.Type().AttributeTypes() {
		path[len(path)-1] = cty.GetAttrStep{Name: key}
		field := target.FieldByName(key)
		if !field.IsValid() {
			result = multierror.Append(result, path.NewErrorf("unsupported attribute %q", key))
		} else if err := fromCtyValue(val.GetAttr(key), field, path); err != nil {
			result = multierror.Append(result, err)
		}
	}

	path = path[:len(path)-1]
	return result
}
