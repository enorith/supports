package reflection

import "reflect"

func StructType(v interface{}) reflect.Type {
	if t, ok := v.(reflect.Type); ok {
		return t
	}

	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

func StructValue(v interface{}) reflect.Value {
	if va, ok := v.(reflect.Value); ok {
		return va
	}

	va := reflect.ValueOf(v)

	if va.Kind() == reflect.Ptr {
		va = va.Elem()
	}

	return va
}

func TypeString(abstract interface{}) (typeString string) {

	if t, ok := abstract.(string); ok {
		typeString = t
	} else if t, ok := abstract.(reflect.Type); ok {
		typeString = t.String()
	} else {
		typeString = reflect.TypeOf(abstract).String()
	}

	return
}