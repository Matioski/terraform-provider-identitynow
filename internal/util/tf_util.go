package util

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"strings"
)

func ConvertTFModelToMap(obj interface{}) map[string]interface{} {
	attributes := make(map[string]interface{})
	objRef := reflect.ValueOf(obj)
	for j := 0; j < objRef.NumField(); j++ {
		field := objRef.Type().Field(j)
		tag := strings.Split(field.Tag.Get("json"), ",")[0]
		switch field.Type {
		case reflect.TypeOf(types.String{}):
			tfField := objRef.Field(j).Interface().(types.String)
			if !tfField.IsNull() && !tfField.IsUnknown() {
				attributes[tag] = tfField.ValueString()
			}
		case reflect.TypeOf(types.Bool{}):
			tfField := objRef.Field(j).Interface().(types.Bool)
			if !tfField.IsNull() && !tfField.IsUnknown() {
				attributes[tag] = tfField.ValueBool()
			}
		}
	}
	return attributes
}

func GetTFStringPointer(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	return value.ValueStringPointer()
}
