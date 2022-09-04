package hcl_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	// Modules
	httpserver "github.com/mutablelogic/terraform-provider-nginx/pkg/httpserver"
	cty "github.com/zclconf/go-cty/cty"
	slices "golang.org/x/exp/slices"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/hcl"
)

func Test_Cty_001(t *testing.T) {
	var tests = []struct {
		In  cty.Value
		Out any
	}{
		{cty.StringVal("hello"), "hello"},
		{cty.NilVal, ""},
		{cty.NilVal, int(0)},
		{cty.NumberIntVal(0), int(0)},
		{cty.NumberIntVal(1), int(1)},
		{cty.NumberIntVal(math.MinInt), int(math.MinInt)},
		{cty.NumberIntVal(math.MaxInt), int(math.MaxInt)},
		{cty.NilVal, int8(0)},
		{cty.NumberIntVal(0), int8(0)},
		{cty.NumberIntVal(1), int8(1)},
		{cty.NumberIntVal(math.MinInt8), int8(math.MinInt8)},
		{cty.NumberIntVal(math.MaxInt8), int8(math.MaxInt8)},
		{cty.NilVal, int16(0)},
		{cty.NumberIntVal(0), int16(0)},
		{cty.NumberIntVal(1), int16(1)},
		{cty.NumberIntVal(math.MinInt16), int16(math.MinInt16)},
		{cty.NumberIntVal(math.MaxInt16), int16(math.MaxInt16)},
		{cty.NilVal, int32(0)},
		{cty.NumberIntVal(0), int32(0)},
		{cty.NumberIntVal(1), int32(1)},
		{cty.NumberIntVal(math.MinInt32), int32(math.MinInt32)},
		{cty.NumberIntVal(math.MaxInt32), int32(math.MaxInt32)},
		{cty.NilVal, int64(0)},
		{cty.NumberIntVal(0), int64(0)},
		{cty.NumberIntVal(1), int64(1)},
		{cty.NumberIntVal(math.MinInt64), int64(math.MinInt64)},
		{cty.NumberIntVal(math.MaxInt64), int64(math.MaxInt64)},
		{cty.NilVal, uint(0)},
		{cty.NumberUIntVal(0), uint(0)},
		{cty.NumberUIntVal(1), uint(1)},
		{cty.NumberUIntVal(math.MaxUint), uint(math.MaxUint)},
		{cty.NilVal, uint8(0)},
		{cty.NumberUIntVal(0), uint8(0)},
		{cty.NumberUIntVal(1), uint8(1)},
		{cty.NumberUIntVal(math.MaxUint8), uint8(math.MaxUint8)},
		{cty.NilVal, uint16(0)},
		{cty.NumberUIntVal(0), uint16(0)},
		{cty.NumberUIntVal(1), uint16(1)},
		{cty.NumberUIntVal(math.MaxUint16), uint16(math.MaxUint16)},
		{cty.NilVal, uint32(0)},
		{cty.NumberUIntVal(0), uint32(0)},
		{cty.NumberUIntVal(1), uint32(1)},
		{cty.NumberUIntVal(math.MaxUint32), uint32(math.MaxUint32)},
		{cty.NilVal, uint64(0)},
		{cty.NumberUIntVal(0), uint64(0)},
		{cty.NumberUIntVal(1), uint64(1)},
		{cty.NumberUIntVal(math.MaxUint64), uint64(math.MaxUint64)},
		{cty.NilVal, float32(0)},
		{cty.NilVal, float64(0)},
		{cty.NumberFloatVal(math.MaxFloat32), float32(math.MaxFloat32)},
		{cty.NumberFloatVal(math.MaxFloat64), float64(math.MaxFloat64)},
		{cty.NumberFloatVal(-math.MaxFloat32), float32(-math.MaxFloat32)},
		{cty.NumberFloatVal(-math.MaxFloat64), float64(-math.MaxFloat64)},
		{cty.NumberFloatVal(math.Inf(-1)), math.Inf(-1)},
		{cty.NumberFloatVal(math.Inf(+1)), math.Inf(+1)},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			switch val := test.Out.(type) {
			case string:
				var out string
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %q", out)
				}
			case int:
				var out int
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case int8:
				var out int8
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case int16:
				var out int16
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case int32:
				var out int32
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case int64:
				var out int64
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case uint:
				var out uint
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case uint8:
				var out uint8
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case uint16:
				var out uint16
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case uint32:
				var out uint32
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case uint64:
				var out uint64
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case float32:
				var out float32
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			case float64:
				var out float64
				if err := FromCtyValue(test.In, &out); err != nil {
					t.Error(err)
				} else if out != val {
					t.Errorf("unexpected value returned: %v (expected %v)", out, val)
				}
			default:
				t.Fatal("Unexpected type for test", i, "=>", reflect.TypeOf(test.Out))
			}
		})
	}
}

func Test_Cty_002(t *testing.T) {
	var tests = []struct {
		In  cty.Value
		Out any
	}{
		{cty.NilVal, nil},
		{cty.ListValEmpty(cty.String), []string{}},
		{cty.ListVal([]cty.Value{cty.StringVal("a")}), []string{"a"}},
		{cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}), []string{"a", "b"}},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out []string
			if err := FromCtyValue(test.In, &out); err != nil {
				t.Error(err)
			} else if test.Out == nil {
				if out != nil {
					t.Errorf("unexpected value returned: %v (expected %v)", out, test.Out)
				}
			} else if slices.Compare(out, test.Out.([]string)) != 0 {
				t.Errorf("unexpected value returned: %q (expected %q)", out, test.Out)
			}
		})
	}
}

func Test_Cty_003(t *testing.T) {
	var tests = []struct {
		In  cty.Value
		Out any
	}{
		{cty.NilVal, nil},
		{cty.ListVal([]cty.Value{cty.StringVal("a")}), [1]string{"a"}},
		{cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}), [2]string{"a", "b"}},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out []string
			if err := FromCtyValue(test.In, &out); err != nil {
				t.Error(err)
			} else if test.Out == nil {
				if out != nil {
					t.Errorf("unexpected value returned: %v (expected %v)", out, test.Out)
				}
			}
		})
	}
}

func Test_Cty_004(t *testing.T) {
	var tests = []struct {
		In  cty.Value
		Out any
	}{
		{cty.NilVal, nil},
		{cty.MapValEmpty(cty.String), map[string]string{}},
		{cty.MapVal(map[string]cty.Value{"a": cty.StringVal("a")}), map[string]string{"a": "a"}},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out map[string]string
			if err := FromCtyValue(test.In, &out); err != nil {
				t.Error(err)
			} else if test.Out == nil {
				if out != nil {
					t.Errorf("unexpected value returned: %v (expected %v)", out, test.Out)
				}
			} else if reflect.DeepEqual(out, test.Out.(map[string]string)) != true {
				t.Errorf("unexpected value returned: %v (expected %v)", out, test.Out)
			}
		})
	}
}

func Test_Cty_005(t *testing.T) {
	type testStruct struct {
		A string
	}

	var tests = []struct {
		In  cty.Value
		Out testStruct
	}{
		{cty.NilVal, testStruct{}},
		{cty.ObjectVal(map[string]cty.Value{}), testStruct{}},
		{cty.ObjectVal(map[string]cty.Value{"A": cty.StringVal("A")}), testStruct{A: "A"}},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out testStruct
			out.A = "default"
			if err := FromCtyValue(test.In, &out); err != nil {
				t.Error(err)
			} else if reflect.DeepEqual(out, test.Out) != true {
				t.Errorf("unexpected value returned: %v (expected %v)", out, test.Out)
			}
		})
	}
}

func Test_Cty_006(t *testing.T) {
	var tests = []struct {
		In  cty.Value
		Out httpserver.ServerConfig
	}{
		{cty.NilVal, httpserver.ServerConfig{}},
		{cty.ObjectVal(map[string]cty.Value{"Label": cty.StringVal("test")}), httpserver.ServerConfig{Label: "test"}},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out httpserver.ServerConfig
			if err := FromCtyValue(test.In, &out); err != nil {
				t.Error(err)
			} else if reflect.DeepEqual(out, test.Out) != true {
				t.Errorf("unexpected value returned: %v (expected %v)", out, test.Out)
			}
		})
	}
}

func Test_Cty_007(t *testing.T) {
	type varblock struct {
		Label       string
		Type        string
		Default     any
		Description string
	}

	var tests = []struct {
		In  cty.Value
		Out varblock
	}{
		{cty.NilVal, varblock{}},
		{cty.ObjectVal(map[string]cty.Value{"Label": cty.StringVal("test")}), varblock{Label: "test"}},
		{cty.ObjectVal(map[string]cty.Value{"Default": cty.StringVal("test")}), varblock{Default: "test"}},
		{cty.ObjectVal(map[string]cty.Value{"Default": cty.NumberFloatVal(math.MaxFloat64)}), varblock{Default: math.MaxFloat64}},
		{cty.ObjectVal(map[string]cty.Value{"Default": cty.BoolVal(true)}), varblock{Default: true}},
		{cty.ObjectVal(map[string]cty.Value{"Default": cty.BoolVal(false)}), varblock{Default: false}},
		{cty.ObjectVal(map[string]cty.Value{"Default": cty.NumberUIntVal(math.MaxUint64)}), varblock{Default: uint64(math.MaxUint64)}},
		{cty.ObjectVal(map[string]cty.Value{"Default": cty.NumberIntVal(math.MaxInt64)}), varblock{Default: int64(math.MaxInt64)}},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out varblock
			if err := FromCtyValue(test.In, &out); err != nil {
				t.Error(err)
			} else if reflect.DeepEqual(out, test.Out) != true {
				t.Errorf("unexpected value returned: %v (expected %v)", out, test.Out)
			}
		})
	}
}

func Test_Cty_008(t *testing.T) {
	type varblock struct {
		Label   string
		Type    string
		Default any
	}

	var tests = []struct {
		In  cty.Value
		Out varblock
	}{
		{cty.NilVal, varblock{}},
		{cty.TupleVal([]cty.Value{cty.StringVal("label"), cty.StringVal("type"), cty.StringVal("default")}), varblock{Label: "label", Type: "type", Default: "default"}},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out varblock
			if err := FromCtyValue(test.In, &out); err != nil {
				t.Error(err)
			} else if reflect.DeepEqual(out, test.Out) != true {
				t.Errorf("unexpected value returned: %v (expected %v)", out, test.Out)
			}
		})
	}
}
