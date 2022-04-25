package reflection

import "reflect"

func StructType(v interface{}) reflect.Type {
	t := TypeOf(v)
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}

	return t
}

func StructValue(v interface{}) reflect.Value {
	va := ValueOf(v)
	if va.Kind() == reflect.Ptr && !va.IsNil() {
		return va.Elem()
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

func TypeOf(v interface{}) reflect.Type {
	if va, ok := v.(reflect.Type); ok {
		return va
	}

	return reflect.TypeOf(v)
}

func ValueOf(v interface{}) reflect.Value {
	if va, ok := v.(reflect.Value); ok {
		return va
	}

	return reflect.ValueOf(v)
}

func SubStructOf(child, parent interface{}) int {
	ct := StructType(child)
	pt := StructType(parent)

	if ct.Kind() == reflect.Struct {
		for i := 0; i < ct.NumField(); i++ {
			field := ct.Field(i)
			if field.Anonymous && StructType(field.Type) == pt {
				return i
			}
		}
	}

	return -1
}

func Implements(abs interface{}, of reflect.Type) bool {
	t := TypeOf(abs)
	return t.Implements(of)
}

func InterfaceType[T interface{}]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
