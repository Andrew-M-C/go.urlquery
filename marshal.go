package urlquery

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
)

func marshalToValues(in interface{}) (kv url.Values, err error) {
	v, err := validateMarshalParam(in)
	if err != nil {
		return nil, err
	}

	kv = url.Values{}
	_ = v

	// TODO:
	return kv, nil
}

func validateMarshalParam(in interface{}) (v reflect.Value, err error) {
	if in == nil {
		err = errors.New("no data provided")
		return
	}

	v = reflect.ValueOf(in)
	t := v.Type()

	if k := t.Kind(); k == reflect.Struct {
		// OK
		return v, nil

	} else if k == reflect.Ptr {
		if v.IsNil() {
			err = errors.New("nil pointer of a struct is not supported")
			return
		}
		t = t.Elem()
		if t.Kind() != reflect.Struct {
			err = fmt.Errorf("invalid type of input: %v", t)
			return
		}

		return v.Elem(), nil
	}

	err = fmt.Errorf("invalid type of input: %v", t)
	return
}
