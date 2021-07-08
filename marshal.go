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

		readFieldToKV(&fv, &ft, kv, "")
	}

	return kv, nil
}

func readFieldToKV(fv *reflect.Value, ft *reflect.StructField, kv url.Values, keyPrefix string) {
	if ft.Anonymous {
		numField := fv.NumField()
		for i := 0; i < numField; i++ {
			ffv := fv.Field(i)
			fft := ft.Type.Field(i)

			readFieldToKV(&ffv, &fft, kv, keyPrefix)
		}
		return
	}
	if !fv.CanInterface() {
		return
	}

	tg := readTag(ft, "url")
	if tg.Name() == "-" {
		return
	}

	// 写 KV 值
	readFieldValToKV(fv, tg, kv, keyPrefix)
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

func readFieldValToKV(v *reflect.Value, tg tags, kv url.Values, keyPrefix string) {
	key := tg.Name()
	if keyPrefix != "" {
		key = keyPrefix + "." + key
	}

	val := ""
	var vals []string
	omitempty := tg.Has("omitempty")
	isSliceOrArray := false

	switch v.Type().Kind() {
	default:
		omitempty = true
	case reflect.String:
		val = v.String()
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		val = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		val = strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		val = fmt.Sprintf("%v", v.Bool())
	case reflect.Float64, reflect.Float32:
		val = strconv.FormatFloat(v.Float(), 'f', -1, 64)

	case reflect.Slice, reflect.Array:
		isSliceOrArray = true
		elemTy := v.Type().Elem()
		switch elemTy.Kind() {
		default:
			// 什么也不做，omitempty 对数组而言没有意义
		case reflect.String:
			vals = readStringArray(v)
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			vals = readIntArray(v)
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			vals = readUintArray(v)
		case reflect.Bool:
			vals = readBoolArray(v)
		case reflect.Float64, reflect.Float32:
			vals = readFloatArray(v)
		}

	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		elem := v.Elem()
		v = &elem
		fallthrough

	case reflect.Struct:
		t := v.Type()
		numField := t.NumField()

		for i := 0; i < numField; i++ {
			fv := v.Field(i) // field value
			ft := t.Field(i) // field type

			readFieldToKV(&fv, &ft, kv, key)
		}
		return // 不再往下走，而是由被递归的函数来完成 kv.Set
	}

	// 数组使用 Add 函数
	if isSliceOrArray {
		for _, v := range vals {
			kv.Add(key, v)
		}
		return
	}

	if val == "" && omitempty {
		return
	}
	kv.Set(key, val)
}

func readStringArray(v *reflect.Value) (vals []string) {
	count := v.Len()

	for i := 0; i < count; i++ {
		child := v.Index(i)
		s := child.String()
		vals = append(vals, s)
	}

	return
}

func readIntArray(v *reflect.Value) (vals []string) {
	count := v.Len()

	for i := 0; i < count; i++ {
		child := v.Index(i)
		v := child.Int()
		vals = append(vals, strconv.FormatInt(v, 10))
	}

	return
}

func readUintArray(v *reflect.Value) (vals []string) {
	count := v.Len()

	for i := 0; i < count; i++ {
		child := v.Index(i)
		v := child.Uint()
		vals = append(vals, strconv.FormatUint(v, 10))
	}

	return
}

func readBoolArray(v *reflect.Value) (vals []string) {
	count := v.Len()

	for i := 0; i < count; i++ {
		child := v.Index(i)
		if child.Bool() {
			vals = append(vals, "true")
		} else {
			vals = append(vals, "false")
		}
	}

	return
}

func readFloatArray(v *reflect.Value) (vals []string) {
	count := v.Len()

	for i := 0; i < count; i++ {
		child := v.Index(i)
		v := child.Float()
		vals = append(vals, strconv.FormatFloat(v, 'f', -1, 64))
	}

	return
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
