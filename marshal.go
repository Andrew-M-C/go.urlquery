package urlquery

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func marshalToValues(in interface{}) (kv url.Values, err error) {
	v, err := validateMarshalParam(in)
	if err != nil {
		return nil, err
	}

	t := v.Type()
	numField := t.NumField() // 结构体下所有字段的数量

	kv = url.Values{}

	// 迭代每一个字段
	for i := 0; i < numField; i++ {
		fv := v.Field(i) // field value
		ft := t.Field(i) // field type

		if ft.Anonymous {
			// TODO: 后文再处理
			continue
		}
		if !fv.CanInterface() {
			continue
		}

		tg := readTag(&ft, "url")
		if tg.Name() == "-" {
			continue
		}

		str, ok := readFieldVal(&fv, tg)
		if !ok {
			continue
		}
		if str == "" && tg.Has("omitempty") {
			continue
		}

		// 写 KV 值
		kv.Set(tg.Name(), str)
	}

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

func readFieldVal(v *reflect.Value, tag tags) (s string, ok bool) {
	switch v.Type().Kind() {
	default:
		return "", false
	case reflect.String:
		return v.String(), true
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return strconv.FormatInt(v.Int(), 10), true
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return strconv.FormatUint(v.Uint(), 10), true
	case reflect.Bool:
		return fmt.Sprintf("%v", v.Bool()), true
	case reflect.Float64, reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), true
	}
}

type tags []string

func readTag(ft *reflect.StructField, tag string) tags {
	tg := ft.Tag.Get(tag)

	// 如果 tag 配置非空，则返回
	if tg != "" {
		res := strings.Split(tg, ",")
		if res[0] != "" {
			return res
		}
		return append(tags{ft.Name}, res[1:]...)
	}
	// 如果 tag 配置为空，则返回字段名
	return tags{ft.Name}
}

// Name 表示当前 tag 所定义的第一个字段，这个字段必须是名称
func (tg tags) Name() string {
	return tg[0]
}

// Has 判断当前 tag 是否配置了某些额外参数值，比如 omitempty
func (tg tags) Has(opt string) bool {
	for i := 1; i < len(tg); i++ {
		t := tg[i]
		if t == opt {
			return true
		}
	}
	return false
}
