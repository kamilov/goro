package reader

import (
	"encoding"
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

type FormReader struct{}

const (
	FormMaxMemory = 32 << 20
	tagForm       = "form"
	skipTagValue  = "-"
)

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func (r FormReader) Read(request *http.Request, data any) error {
	_ = request.ParseMultipartForm(FormMaxMemory)
	return ReadFormData(request.Form, data)
}

func ReadFormData(form map[string][]string, data any) error {
	value := reflect.ValueOf(data)

	if value.Kind() != reflect.Ptr || value.IsNil() {
		return errors.New("data must be a pointer")
	}

	value = indirect(value)

	if value.Kind() != reflect.Struct {
		return errors.New("data must be a pointer to a struct")
	}

	return readForm(form, "", value)
}

func readForm(form map[string][]string, prefix string, value reflect.Value) error {
	value = indirect(value)
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		tag := field.Tag.Get(tagForm)

		if !field.Anonymous && field.PkgPath == "" || tag == skipTagValue {
			continue
		}

		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		name := tag

		if name == "" && !field.Anonymous {
			name = field.Name
		}

		if name != "" && prefix != "" {
			name = prefix + "." + name
		}

		if ok, err := readFormFieldKnownType(form, name, value.Field(i)); err != nil {
			return err
		} else if ok {
			continue
		}

		if fieldType.Kind() != reflect.Struct {
			if err := readFormField(form, name, value.Field(i)); err != nil {
				return err
			}
			continue
		}

		if name == "" {
			name = prefix
		}

		if err := readForm(form, name, value.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func readFormField(form map[string][]string, name string, field reflect.Value) error {
	value, ok := form[name]

	if !ok {
		return nil
	}

	field = indirect(field)

	if field.Kind() != reflect.Slice {
		return setFormFieldValue(field, value[0])
	}

	n := len(value)
	slice := reflect.MakeSlice(field.Type(), n, n)

	for i := 0; i < n; i++ {
		if err := setFormFieldValue(slice.Index(i), value[i]); err != nil {
			return err
		}
	}

	field.Set(slice)

	return nil
}

func setFormFieldValue(index reflect.Value, value string) error {
	switch index.Kind() {
	case reflect.Bool:
		if value == "" {
			value = "false"
		}

		v, err := strconv.ParseBool(value)

		if err != nil {
			return err
		}

		index.SetBool(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			value = "0"
		}

		v, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return err
		}

		index.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			value = "0"
		}

		v, err := strconv.ParseUint(value, 10, 64)

		if err != nil {
			return err
		}

		index.SetUint(v)

	case reflect.Float32, reflect.Float64:
		if value == "" {
			value = "0"
		}

		v, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return err
		}

		index.SetFloat(v)

	case reflect.String:
		index.SetString(value)
		return nil
	}

	return errors.New("Unknown type: " + index.Kind().String())
}

func readFormFieldKnownType(form map[string][]string, name string, field reflect.Value) (bool, error) {
	value, ok := form[name]

	if !ok {
		return false, nil
	}

	field = indirect(field)
	fieldType := field.Type()

	if fieldType.Implements(textUnmarshalerType) {
		return true, field.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value[0]))
	} else if reflect.PtrTo(fieldType).Implements(textUnmarshalerType) {
		return true, field.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value[0]))
	}

	return false, nil
}

func indirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		value = value.Elem()
	}
	return value
}
