package dialect

import (
	"fmt"
	"reflect"
)

//ClickhouseType clickhouse
type ClickhouseType struct{}

//DataTypeOf clickhouse 类型转换
func (ck *ClickhouseType) DataTypeOf(typ reflect.Type) string {
	switch typ.Kind() {
	case reflect.Int8:
		return "Int8"
	case reflect.Int16:
		return "Int16"
	case reflect.Int32, reflect.Int:
		return "Int32"
	case reflect.Int64:
		return "Int64"
	case reflect.Uint8:
		return "UInt8"
	case reflect.Uint16:
		return "UInt16"
	case reflect.Uint32:
		return "UInt32"
	case reflect.Uint64:
		return "UInt64"
	case reflect.Float32:
		return "Float32"
	case reflect.Float64:
		return "Float64"
	case reflect.String:
		return "String"
	default:
		if typ.Name() == "Time" {
			return "DateTime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Name(), typ.Kind()))
}
