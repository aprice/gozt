package gozt

import (
	"reflect"
	"testing"
)

func TestSubstitute(t *testing.T) {
	tests := map[string]struct {
		corpus   string
		scope    *scope
		expected string
	}{
		"simple": {
			"Let's ${foo} the ${baz}",
			newRootScope(map[string]string{"foo": "bar", "baz": "qux"}),
			"Let's bar the qux",
		},
	}
	for k, v := range tests {
		t.Run(k, func(tt *testing.T) {
			actual, err := v.scope.substitute(v.corpus)
			if err != nil {
				tt.Error(err)
			} else if v.expected != actual {
				tt.Errorf("Expected: %v, actual: %v", v.expected, actual)
			}
		})
	}
}

func TestResolveReference(t *testing.T) {
	tests := map[string]struct {
		ref      string
		model    interface{}
		expected interface{}
	}{
		"map":    {"test", map[string]string{"test": "data"}, "data"},
		"subMap": {"test.sub", map[string]interface{}{"test": map[string]string{"sub": "foo"}}, "foo"},
		"array":  {"test.1", map[string]interface{}{"test": [...]string{"foo", "bar"}}, "bar"},
		"slice":  {"test.1", map[string]interface{}{"test": []string{"foo", "baz"}}, "baz"},
		"struct": {"test.Field", map[string]interface{}{"test": struct{ Field string }{"qux"}}, "qux"},
		"int":    {"test", map[string]interface{}{"test": 42}, 42},
		"float":  {"test", map[string]interface{}{"test": 4.2}, 4.2},
	}
	for k, v := range tests {
		t.Run(k, func(tt *testing.T) {
			s := newRootScope(v.model)
			actual, err := s.resolveReference(v.ref)
			if err != nil {
				tt.Error(err)
			} else if v.expected != actual {
				tt.Errorf("Expected: %v, actual: %v", v.expected, actual)
			}
		})
	}
}

func TestGetProp(t *testing.T) {
	tests := map[string]struct {
		ref      string
		model    interface{}
		expected interface{}
	}{
		"map":    {"test", map[string]string{"test": "data"}, "data"},
		"subMap": {"test.sub", map[string]interface{}{"test": map[string]string{"sub": "foo"}}, "foo"},
		"array":  {"test.1", map[string]interface{}{"test": [...]string{"foo", "bar"}}, "bar"},
		"slice":  {"test.1", map[string]interface{}{"test": []string{"foo", "baz"}}, "baz"},
		"struct": {"test.Field", map[string]interface{}{"test": struct{ Field string }{"qux"}}, "qux"},
		"int":    {"test", map[string]interface{}{"test": 42}, 42},
		"float":  {"test", map[string]interface{}{"test": 4.2}, 4.2},
	}
	for k, v := range tests {
		t.Run(k, func(tt *testing.T) {
			actual, err := getProp(v.ref, v.model)
			if err != nil {
				tt.Error(err)
			} else if !reflect.DeepEqual(v.expected, actual) {
				tt.Errorf("Expected: %v, actual: %v", v.expected, actual)
			}
		})
	}
}

func TestIsTruthy(t *testing.T) {
	tests := map[string]struct {
		in  interface{}
		out bool
	}{
		"true":          {true, true},
		"false":         {false, false},
		"nil":           {nil, false},
		"struct":        {struct{}{}, true},
		"pointer":       {new(struct{}), true},
		"int0":          {0, false},
		"int1":          {1, true},
		"float0":        {0.0, false},
		"float1":        {1.0, true},
		"stringempty":   {"", false},
		"stringfalse":   {"false", false},
		"stringtrue":    {"true", true},
		"emptyslice":    {[]int{}, false},
		"nonemptyslice": {[]int{1}, true},
	}
	for k, v := range tests {
		t.Run(k, func(tt *testing.T) {
			actual := isTruthy(v.in)
			if v.out != actual {
				tt.Errorf("isTruthy(%v) expected %v, actual %v", v.in, v.out, actual)
			}
		})
	}
}
