package env

import (
	"os"
	"reflect"
	"strconv"
)

const (
	EnvTag = "env"
)

func ReadEnv[Type any]() Type {
	var container Type
	_value := reflect.ValueOf(&container).Elem()
	updateEnvRecoursive(_value)
	return container
}

func updateEnvRecoursive(_value reflect.Value) {
	for n := range _value.NumField() {

		field := _value.Field(n)

		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
		}

		if field.Kind() == reflect.Interface {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			updateEnvRecoursive(field)
			continue
		}

		env, ok := _value.Type().Field(n).Tag.Lookup(EnvTag)
		if !ok {
			continue
		}

		val, ok := os.LookupEnv(env)
		if !ok {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(val)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				continue
			}
			field.SetInt(int64(v))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				continue
			}
			field.SetUint(v)

		case reflect.Bool:
			v, err := strconv.ParseBool(val)
			if err != nil {
				continue
			}
			field.SetBool(v)
		}

	}
}
