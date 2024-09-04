package patch

import (
	"encoding/json"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
)

func ConvertFromBetaToV3(beta []sailpointBeta.JsonPatchOperation) ([]sailpointV3.JsonPatchOperation, error) {
	v3 := make([]sailpointV3.JsonPatchOperation, len(beta))
	for i, b := range beta {
		value, err := convertValue(b.Value)
		if err != nil {
			return nil, err
		}
		v3[i] = sailpointV3.JsonPatchOperation{
			Op:    b.Op,
			Path:  b.Path,
			Value: value,
		}
	}
	return v3, nil
}

func convertValue(value *sailpointBeta.JsonPatchOperationValue) (*sailpointV3.JsonPatchOperationValue, error) {
	if value == nil {
		return nil, nil
	}
	var arrayInnerPtr *[]sailpointV3.ArrayInner
	if value.ArrayOfArrayInner != nil {
		jsonBytes, err := json.Marshal(value.ArrayOfArrayInner)
		if err != nil {
			return nil, err
		}
		var arrayInner []sailpointV3.ArrayInner
		err = json.Unmarshal(jsonBytes, &arrayInner)
		if err != nil {
			return nil, err
		}
		arrayInnerPtr = &arrayInner
	}
	return &sailpointV3.JsonPatchOperationValue{
		ArrayOfArrayInner:       arrayInnerPtr,
		Int32:                   value.Int32,
		MapmapOfStringinterface: value.MapmapOfStringinterface,
		String:                  value.String,
		Bool:                    value.Bool,
	}, nil
}
