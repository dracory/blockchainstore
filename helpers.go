package blockchainstore

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// anyToMapStringString converts any interface to map[string]string
func anyToMapStringString(data any) map[string]string {
	if data == nil {
		return nil
	}

	t := reflect.TypeOf(data)

	// Is it a map?
	kind := t.Kind()
	if kind != reflect.Map {
		return nil
	}

	elemKeyKind := t.Key().Kind()
	elemKValueKind := t.Elem().Kind()

	// Is it a match?
	if elemKeyKind == reflect.String && elemKValueKind == reflect.String {
		return data.(map[string]string)
	}

	// Is it map[string]any? Convert to map[string]string
	if elemKeyKind == reflect.String && elemKValueKind == reflect.Interface {
		mapStringInterface := data.(map[string]interface{})
		mapStringString := map[string]string{}
		for key, value := range mapStringInterface {
			mapStringString[key] = toString(value)
		}
		return mapStringString
	}

	// Last attempt to convert, may panic and need another condition
	return data.(map[string]string)
}

// toString converts an interface to string
func toString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v

	case nil:
		return ""

	case []byte:
		return btos(v)

	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)

	case float64:
		return strconv.FormatFloat(v, 'f', 4, 64)

	default:
		return fmt.Sprint(v)
	}
}

// btos converts bytes to string
func btos(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// fromJSON unmarshals jsonString to interface
func fromJSON(jsonString string, valueDefault interface{}) (interface{}, error) {
	var e interface{}

	jsonError := json.Unmarshal([]byte(jsonString), &e)

	if jsonError != nil {
		return valueDefault, jsonError
	}

	return e, nil
}
