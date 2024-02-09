package polygo

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
)

type Decoder[T any] struct {
	Type      T
	FieldName string
	TypeMap   map[string]reflect.Type
}

func NewDecoder[T any](field string) *Decoder[T] {
	return &Decoder[T]{
		FieldName: field,
		TypeMap:   map[string]reflect.Type{},
	}
}

func (d *Decoder[T]) Register(value string, v any) *Decoder[T] {
	d.TypeMap[value] = reflect.TypeOf(v)
	return d
}

func (d *Decoder[T]) UnmarshalArray(b []byte) ([]T, error) {
	res := gjson.ParseBytes(b)
	return d.unmarshalArray("", res)
}

func (d *Decoder[T]) UnmarshalInnerArray(path string, b []byte) ([]T, error) {
	res := gjson.ParseBytes(b)
	return d.unmarshalArray(path, res)
}

func (d *Decoder[T]) unmarshalArray(path string, res gjson.Result) ([]T, error) {
	if path != "" {
		res = res.Get(path)
	}

	if !res.IsArray() {
		return nil, errors.New("object is not an array")
	}

	arr := []T{}
	var err error

	res.ForEach(func(key, value gjson.Result) bool {
		a, errParse := d.unmarshalObject("", value)
		if errParse != nil {
			err = errParse
			return false
		}

		arr = append(arr, a)
		return true
	})

	return arr, err
}

func (d *Decoder[T]) UnmarshalObject(b []byte) (T, error) {
	res := gjson.ParseBytes(b)
	return d.unmarshalObject("", res)
}

func (d *Decoder[T]) UnmarshalInnerObject(path string, b []byte) (T, error) {
	res := gjson.ParseBytes(b)
	return d.unmarshalObject(path, res)
}

func (d *Decoder[T]) unmarshalObject(path string, res gjson.Result) (T, error) {
	var zero T

	if path != "" {
		res = res.Get(path)
	}

	if res.IsArray() {
		return zero, errors.New("cannot unmarshal object: JSON is array")
	}

	fieldValue := res.Get(d.FieldName).String()
	// TODO check if string
	if fieldValue == "" {
		return zero, fmt.Errorf("field '%s' not found", d.FieldName)
	}

	matchedType, found := d.TypeMap[fieldValue]
	if !found {
		return zero, fmt.Errorf("type '%s' not registered", fieldValue)
	}

	var v reflect.Value
	if matchedType.Kind() == reflect.Ptr {
		v = reflect.New(matchedType.Elem())
	} else {
		v = reflect.New(matchedType)
	}

	err := json.Unmarshal([]byte(res.Raw), v.Interface())
	if err != nil {
		return zero, err
	}

	a, ok := v.Interface().(T)
	if !ok {
		return zero, errors.New("error casting")
	}

	return a, nil
}
