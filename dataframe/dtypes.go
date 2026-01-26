// Package dataframe provides DataFrame and Series data structures
// similar to Python's pandas library for data analysis in Go.
package dataframe

import (
	"fmt"
	"reflect"
	"time"
)

// DType represents the data type of a Series
type DType int

const (
	// DTypeUnknown represents unknown data type
	DTypeUnknown DType = iota
	// DTypeInt64 represents 64-bit integer
	DTypeInt64
	// DTypeFloat64 represents 64-bit floating point
	DTypeFloat64
	// DTypeString represents string type
	DTypeString
	// DTypeBool represents boolean type
	DTypeBool
	// DTypeDateTime represents date/time type
	DTypeDateTime
	// DTypeObject represents any type (interface{})
	DTypeObject
)

// String returns the string representation of DType
func (d DType) String() string {
	switch d {
	case DTypeInt64:
		return "int64"
	case DTypeFloat64:
		return "float64"
	case DTypeString:
		return "string"
	case DTypeBool:
		return "bool"
	case DTypeDateTime:
		return "datetime"
	case DTypeObject:
		return "object"
	default:
		return "unknown"
	}
}

// InferDType infers the DType from a Go value
func InferDType(v interface{}) DType {
	if v == nil {
		return DTypeObject
	}

	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return DTypeInt64
	case float32, float64:
		return DTypeFloat64
	case string:
		return DTypeString
	case bool:
		return DTypeBool
	case time.Time:
		return DTypeDateTime
	default:
		return DTypeObject
	}
}

// InferDTypeFromSlice infers the DType from a slice of values
func InferDTypeFromSlice(values []interface{}) DType {
	if len(values) == 0 {
		return DTypeObject
	}

	// Find first non-nil value
	var firstType DType = DTypeObject
	for _, v := range values {
		if v != nil {
			firstType = InferDType(v)
			break
		}
	}

	return firstType
}

// ConvertToType converts a value to the specified DType
func ConvertToType(v interface{}, dtype DType) (interface{}, error) {
	if v == nil {
		return nil, nil
	}

	switch dtype {
	case DTypeInt64:
		return toInt64(v)
	case DTypeFloat64:
		return toFloat64(v)
	case DTypeString:
		return toString(v)
	case DTypeBool:
		return toBool(v)
	case DTypeDateTime:
		return toDateTime(v)
	default:
		return v, nil
	}
}

func toInt64(v interface{}) (int64, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(rv.Float()), nil
	case reflect.String:
		var result int64
		_, err := fmt.Sscanf(rv.String(), "%d", &result)
		return result, err
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

func toFloat64(v interface{}) (float64, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return rv.Float(), nil
	case reflect.String:
		var result float64
		_, err := fmt.Sscanf(rv.String(), "%f", &result)
		return result, err
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

func toString(v interface{}) (string, error) {
	return fmt.Sprintf("%v", v), nil
}

func toBool(v interface{}) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() != 0, nil
	case float32, float64:
		return reflect.ValueOf(v).Float() != 0, nil
	case string:
		return val != "" && val != "0" && val != "false", nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

func toDateTime(v interface{}) (time.Time, error) {
	switch val := v.(type) {
	case time.Time:
		return val, nil
	case string:
		// Try common date formats
		formats := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02",
			"2006/01/02",
			"01/02/2006",
			"02-01-2006",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, val); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot parse '%s' as datetime", val)
	default:
		return time.Time{}, fmt.Errorf("cannot convert %T to datetime", v)
	}
}

// IsNA checks if a value is considered as NA (Not Available)
func IsNA(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == "" || val == "NA" || val == "NaN" || val == "null"
	case float64:
		return val != val // NaN check
	case float32:
		return val != val // NaN check
	}
	return false
}
