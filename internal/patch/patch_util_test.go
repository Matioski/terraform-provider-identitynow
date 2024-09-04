package patch

import (
	"encoding/json"
	"fmt"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertFromBetaToV3(t *testing.T) {
	var arrayBeta []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"TEST\"]"), &arrayBeta)
	var arrayV3 []sailpointV3.ArrayInner
	json.Unmarshal([]byte("[\"TEST\"]"), &arrayV3)

	type args struct {
		beta []sailpointBeta.JsonPatchOperation
	}
	tests := []struct {
		name    string
		args    args
		want    []sailpointV3.JsonPatchOperation
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "TestConvertFromBetaToV3",
			args: args{beta: []sailpointBeta.JsonPatchOperation{
				{
					Op:   "add",
					Path: "/attrString",
					Value: &sailpointBeta.JsonPatchOperationValue{
						String: sailpointBeta.PtrString("newValue"),
					},
				},
				{
					Op:   "replace",
					Path: "/attrBool",
					Value: &sailpointBeta.JsonPatchOperationValue{
						Bool: sailpointBeta.PtrBool(true),
					},
				},
				{
					Op:   "replace",
					Path: "/attrInt",
					Value: &sailpointBeta.JsonPatchOperationValue{
						Int32: sailpointBeta.PtrInt32(222),
					},
				},
				{
					Op:   "replace",
					Path: "/attrArray",
					Value: &sailpointBeta.JsonPatchOperationValue{
						ArrayOfArrayInner: &arrayBeta,
					},
				},
				{
					Op:   "add",
					Path: "/attrMap",
					Value: &sailpointBeta.JsonPatchOperationValue{
						MapmapOfStringinterface: &map[string]interface{}{
							"id": "newId",
						},
					},
				},
			}},
			want: []sailpointV3.JsonPatchOperation{
				{
					Op:   "add",
					Path: "/attrString",
					Value: &sailpointV3.JsonPatchOperationValue{
						String: sailpointV3.PtrString("newValue"),
					},
				},
				{
					Op:   "replace",
					Path: "/attrBool",
					Value: &sailpointV3.JsonPatchOperationValue{
						Bool: sailpointV3.PtrBool(true),
					},
				},
				{
					Op:   "replace",
					Path: "/attrInt",
					Value: &sailpointV3.JsonPatchOperationValue{
						Int32: sailpointV3.PtrInt32(222),
					},
				},
				{
					Op:   "replace",
					Path: "/attrArray",
					Value: &sailpointV3.JsonPatchOperationValue{
						ArrayOfArrayInner: &arrayV3,
					},
				},
				{
					Op:   "add",
					Path: "/attrMap",
					Value: &sailpointV3.JsonPatchOperationValue{
						MapmapOfStringinterface: &map[string]interface{}{
							"id": "newId",
						},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertFromBetaToV3(tt.args.beta)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertFromBetaToV3(%v)", tt.args.beta)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvertFromBetaToV3(%v)", tt.args.beta)
		})
	}
}
