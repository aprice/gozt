package gozt

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var errCannotSubReference = "can only sub-reference maps, structs, and arrays"
var errNoSuchIndex = "index %q could not be found"

type scope struct {
	parent *scope
	root   interface{}
	path   string
}

func newScope(parent *scope, root interface{}, path string) *scope {
	return &scope{
		parent,
		root,
		path,
	}
}

func newRootScope(data interface{}) *scope {
	return newScope(nil, data, "")
}

func newChildScope(parent *scope, data interface{}, k string) *scope {
	return &scope{
		parent,
		data,
		strings.Join([]string{parent.path, k}, "."),
	}
}

func (s *scope) substitute(in string) (string, error) {
	escaping := false
	inVar := false
	out := bytes.NewBuffer(make([]byte, 0, len(in)))
	ref := bytes.NewBuffer(make([]byte, 0, 32))
	var last rune
	for _, c := range in {
		var to *bytes.Buffer
		write := true
		if inVar {
			to = ref
		} else {
			to = out
		}
		if escaping {
			to.WriteRune(c)
			continue
		}

		if c == '\\' {
			escaping = true
			write = false
		}
		if !inVar && c == '$' {
			write = false
		} else if last == '$' && c == '{' {
			inVar = true
			write = false
		} else if last == '$' {
			to.WriteRune('$')
		} else if inVar && c == '}' {
			val, err := s.resolveReference(ref.String())
			if err != nil {
				return "", err
			}
			out.WriteString(fmt.Sprint(val))
			ref.Reset()
			inVar = false
			write = false
		}

		if write {
			to.WriteRune(c)
		}
		last = c
	}
	return out.String(), nil
}

func (s *scope) resolveReference(ref string) (interface{}, error) {
	if ref == "" {
		return s.root, nil
	}

	scope := s
	for scope != nil {
		v, err := getProp(ref, scope.root)
		if err == nil {
			return v, nil
		}
		scope = scope.parent
	}

	return nil, fmt.Errorf(errNoSuchIndex, ref)
}

func (s *scope) resolveBoolean(ref string) (bool, error) {
	v, err := s.resolveReference(ref)
	if err != nil {
		return false, err
	}
	return isTruthy(v), nil
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	rv := reflect.ValueOf(v)
	k := rv.Kind()
	if k == reflect.Array || k == reflect.Slice || k == reflect.Map {
		return rv.Len() > 0
	}
	str := fmt.Sprint(v)
	return str != "" && str != "0" && !strings.EqualFold(str, "false")
}

func getProp(ref string, root interface{}) (interface{}, error) {
	if ref == "" {
		return root, nil
	}

	parts := strings.Split(ref, ".")
	parent := reflect.ValueOf(root)
	for _, k := range parts {
		if parent.Kind() == reflect.Interface {
			parent = parent.Elem()
		}
		switch parent.Kind() {
		case reflect.Map:
			v := parent.MapIndex(reflect.ValueOf(k))
			if !v.IsValid() {
				return reflect.ValueOf(nil), fmt.Errorf(errNoSuchIndex, k)
			}
			parent = v
		case reflect.Struct:
			v := parent.FieldByName(k)
			if !v.IsValid() {
				return reflect.ValueOf(nil), fmt.Errorf(errNoSuchIndex, k)
			}
			parent = v
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			i, err := strconv.Atoi(k)
			if err != nil {
				return reflect.ValueOf(nil), err
			}
			if i > parent.Len() {
				return reflect.ValueOf(nil), fmt.Errorf(errNoSuchIndex, k)
			}
			parent = parent.Index(i)
		default:
			return reflect.ValueOf(nil), fmt.Errorf(errCannotSubReference)
		}
	}

	return parent.Interface(), nil
}
