package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// StructToMap 结构体转map
func StructToMap(obj interface{}) map[string]interface{} {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
		objValue = objValue.Elem()
	}
	
	data := make(map[string]interface{})
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		value := objValue.Field(i)
		
		// 跳过私有字段
		if !value.CanInterface() {
			continue
		}
		
		// 获取json标签
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		if jsonTag == "" {
			jsonTag = field.Name
		}
		
		data[jsonTag] = value.Interface()
	}
	return data
}

// MapToStruct map转结构体
func MapToStruct(m map[string]interface{}, obj interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

// AnyToStruct 任意类型转结构体
func AnyToStruct(data interface{}, obj interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, obj)
}

// InterfaceToString 接口转字符串
func InterfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case int:
		return IntToString(v)
	case int64:
		return Int64ToString(v)
	case uint:
		return UintToString(v)
	case float64:
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		bytes, _ := json.Marshal(v)
		return string(bytes)
	}
}
