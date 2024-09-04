package source

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	"reflect"
	"testing"
)

func Test_sourceResource_getConnectorAttributes(t *testing.T) {
	type fields struct {
		apiClient *sailpoint.APIClient
	}
	type args struct {
		model       *sourceModel
		diagnostics *diag.Diagnostics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			name: "Test Merge Of Connector Attributes including array merge", fields: fields{apiClient: nil},
			args: struct {
				model       *sourceModel
				diagnostics *diag.Diagnostics
			}{
				model: &sourceModel{
					ConnectorAttributes:            jsontypes.NewNormalizedValue("{\"key\":\"value\", \"arrayObject\":[{\"firstElement\":\"firstValue\"},{\"secondElement\":\"secondValue\"}], \"innerObject\":{\"innerKey\":\"innerValue\"}, \"array\":[\"a\",\"b\",\"c\"]}"),
					ConnectorAttributesCredentials: jsontypes.NewExactValue("{\"password\":\"pass\", \"arrayObject\":[{\"innerPass1\":\"pass\"},{\"innerPass2\":\"pass\"}],\"innerObject\":{\"innerPass\":\"pass\"}}"),
				},
				diagnostics: &diag.Diagnostics{}},
			want: map[string]interface{}{
				"key": "value",
				"arrayObject": []interface{}{
					map[string]interface{}{
						"firstElement": "firstValue",
						"innerPass1":   "pass",
					},
					map[string]interface{}{
						"secondElement": "secondValue",
						"innerPass2":    "pass",
					},
				},
				"innerObject": map[string]interface{}{
					"innerKey":  "innerValue",
					"innerPass": "pass",
				},
				"array":    []interface{}{"a", "b", "c"},
				"password": "pass",
			},
		},
		{
			name: "Test Complex Array merges", fields: fields{apiClient: nil},
			args: struct {
				model       *sourceModel
				diagnostics *diag.Diagnostics
			}{
				model: &sourceModel{
					ConnectorAttributes:            jsontypes.NewNormalizedValue("{\"first\": [\"a\", \"b\"], \"second\": [\"c\"], \"innerArrayFirst\": [{\"a\": \"a\", \"b\": \"b\"}, {\"c\": \"c\"}]}"),
					ConnectorAttributesCredentials: jsontypes.NewExactValue("{\"first\": [\"c\"], \"second\": [\"d\", \"e\"], \"innerArrayFirst\": [{\"f\": \"f\"}, {\"g\": \"g\", \"d\": \"d\"}, {\"s\": \"s\"}]}"),
				},
				diagnostics: &diag.Diagnostics{}},
			want: map[string]interface{}{
				"innerArrayFirst": []interface{}{
					map[string]interface{}{
						"a": "a",
						"b": "b",
						"f": "f",
					},
					map[string]interface{}{
						"c": "c",
						"g": "g",
						"d": "d",
					},
					map[string]interface{}{
						"s": "s",
					},
				},
				"first":  []interface{}{"a", "b", "c"},
				"second": []interface{}{"c", "d", "e"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &sourceResource{
				apiClient: tt.fields.apiClient,
			}
			if got := r.getConnectorAttributes(tt.args.model, tt.args.diagnostics); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConnectorAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}
