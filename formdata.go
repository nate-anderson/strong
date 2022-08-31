package strong

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
)

func unmarshalFormIntoStruct(form url.Values, out interface{}, outType reflect.Type) error {
	outVal := reflect.ValueOf(out)
	if outType.Kind() != reflect.Pointer {
		return errors.New("destination must be a pointer to a struct if not a map")
	}
	elem := outVal.Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := outType.Elem().Field(i)
		formTag := field.Tag.Get("form")
		if formTag == "" || !field.IsExported() {
			continue
		}

		fieldVal := reflect.ValueOf(form.Get(formTag))
		elem.FieldByName(field.Name).Set(fieldVal)
	}
	return nil
}

func unmarshalFormIntoMap(form url.Values, out interface{}, outType reflect.Type) error {
	if m, ok := out.(map[string]string); !ok {
		return fmt.Errorf("destination must be map[string]string, got %T", out)
	} else {
		for key := range form {
			m[key] = form.Get(key)
		}
	}

	return nil
}

// unmarshalForm unmarshals a url.Values into a struct pointer or a map[string]any
func unmarshalForm(form url.Values, out interface{}) error {
	outType := reflect.TypeOf(out)
	if outType.Kind() == reflect.Pointer {
		if outType.Elem().Kind() == reflect.Map {
			return errors.New("cannot unmarshal into pointer to map: pass map instead")
		}
		return unmarshalFormIntoStruct(form, out, outType)
	}
	if outType.Kind() == reflect.Map {
		return unmarshalFormIntoMap(form, out, outType)
	}
	return fmt.Errorf("form data can only be unmarshaled into struct or map, got %s", outType.Kind().String())
}
