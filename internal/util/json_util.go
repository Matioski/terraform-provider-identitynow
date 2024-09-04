package util

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"reflect"
)

func MarshalToJsonType(v any, diagnostics *diag.Diagnostics) jsontypes.Exact {
	if v == nil {
		return jsontypes.NewExactNull()
	}
	jsonBytes, err := json.Marshal(&v)
	if err != nil {
		diagnostics.AddError(
			"Error Processing JSON value",
			fmt.Sprintf("Cannot marshal Transform Attributes to JSON '%v': %s", v, err.Error()),
		)
		return jsontypes.NewExactValue("")
	}
	return jsontypes.NewExactValue(string(jsonBytes))
}

func MarshalToJsonTypeNormalized(v any, diagnostics *diag.Diagnostics) jsontypes.Normalized {
	if v == nil {
		return jsontypes.NewNormalizedNull()
	}
	jsonBytes, err := json.Marshal(&v)
	if err != nil {
		diagnostics.AddError(
			"Error Processing JSON value",
			fmt.Sprintf("Cannot marshal Transform Attributes to JSON '%v': %s", v, err.Error()),
		)
		return jsontypes.NewNormalizedValue("")
	}
	return jsontypes.NewNormalizedValue(string(jsonBytes))
}

func MarshalToJsonTypeWithDefinedSchema(mapToConvert map[string]interface{}, originalSchemaJson jsontypes.Exact, diagnostics *diag.Diagnostics) jsontypes.Exact {
	if mapToConvert == nil {
		return jsontypes.NewExactNull()
	}
	var referenceSchema map[string]interface{}
	if !originalSchemaJson.IsNull() && !originalSchemaJson.IsUnknown() {
		diagnostic := originalSchemaJson.Unmarshal(&referenceSchema)
		diagnostics.Append(diagnostic...)
		if diagnostic.HasError() {
			return jsontypes.NewExactValue("")
		}
	}
	reducedCopy := make(map[string]interface{})
	copyMapValues(mapToConvert, reducedCopy, referenceSchema)
	return MarshalToJsonType(reducedCopy, diagnostics)
}

func MarshalToJsonTypeWithDefinedSchemaNormalized(mapToConvert map[string]interface{}, originalSchemaJson jsontypes.Normalized, diagnostics *diag.Diagnostics) jsontypes.Normalized {
	if mapToConvert == nil {
		return jsontypes.NewNormalizedNull()
	}
	var referenceSchema map[string]interface{}
	if !originalSchemaJson.IsNull() && !originalSchemaJson.IsUnknown() {
		diagnostic := originalSchemaJson.Unmarshal(&referenceSchema)
		diagnostics.Append(diagnostic...)
		if diagnostic.HasError() {
			return jsontypes.NewNormalizedValue("")
		}
	}
	reducedCopy := make(map[string]interface{})
	copyMapValues(mapToConvert, reducedCopy, referenceSchema)
	return MarshalToJsonTypeNormalized(reducedCopy, diagnostics)
}

func copyMapValues(currentMap, result, referenceSchema map[string]interface{}) map[string]interface{} {
	for key, value := range currentMap {
		if _, ok := referenceSchema[key]; ok {
			if value == nil {
				result[key] = nil
			} else if reflect.TypeOf(value).Kind() == reflect.Map {
				result[key] = copyMapValues(value.(map[string]interface{}), make(map[string]interface{}), referenceSchema[key].(map[string]interface{}))
			} else if reflect.TypeOf(value).Kind() == reflect.Slice {
				result[key] = copySliceValues(value.([]interface{}), make([]interface{}, 0), referenceSchema[key].([]interface{}))
			} else {
				result[key] = value
			}
		}
	}
	return result
}

func copySliceValues(currentSlice, result, referenceSchema []interface{}) []interface{} {
	for i, value := range currentSlice {
		if i < len(referenceSchema) {
			if reflect.TypeOf(value).Kind() == reflect.Map {
				result = append(result, copyMapValues(value.(map[string]interface{}), make(map[string]interface{}), referenceSchema[i].(map[string]interface{})))
			} else if reflect.TypeOf(value).Kind() == reflect.Slice {
				result = append(result, copySliceValues(value.([]interface{}), make([]interface{}, 0), referenceSchema[i].([]interface{})))
			} else {
				result = append(result, value)
			}
		}
	}
	return result
}

func UnmarshalJsonType(jsonObj jsontypes.Exact, diagnostics *diag.Diagnostics) map[string]interface{} {
	if jsonObj.IsNull() || jsonObj.IsUnknown() {
		return nil
	}
	var obj map[string]interface{}
	diagnostic := jsonObj.Unmarshal(&obj)
	diagnostics.Append(diagnostic...)
	if diagnostic.HasError() {
		return make(map[string]interface{})
	}
	return obj
}

func UnmarshalJsonTypeNormalized(jsonObj jsontypes.Normalized, diagnostics *diag.Diagnostics) map[string]interface{} {
	if jsonObj.IsNull() || jsonObj.IsUnknown() {
		return nil
	}
	var obj map[string]interface{}
	diagnostic := jsonObj.Unmarshal(&obj)
	diagnostics.Append(diagnostic...)
	if diagnostic.HasError() {
		return make(map[string]interface{})
	}
	return obj
}
