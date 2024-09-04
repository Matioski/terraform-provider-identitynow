//go:build !integration

package patch

import (
	"encoding/json"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_Source_OwnerPatch_Match_NoPatch(t *testing.T) {
	modOwn := sailpointV3.SourceOwner{
		Id:   sailpointV3.PtrString("123"),
		Type: sailpointV3.PtrString("compareInnerObject"),
	}
	curOwn := sailpointV3.SourceOwner{
		Id:   sailpointV3.PtrString("123"),
		Type: sailpointV3.PtrString("compareInnerObject"),
	}

	mod := sailpointV3.Source{Owner: modOwn}
	cur := sailpointV3.Source{Owner: curOwn}
	patcOp, err := NewSourcePatchBuilder(&mod, &cur).GenerateJsonPatch()
	if err != nil {
		t.Error(err)
	}
	if len(patcOp) != 0 {
		t.Error("Expected no patch operations")
	}
}

func Test_Source_OwnerPatch_ReplacePatch(t *testing.T) {
	modOwn := sailpointV3.SourceOwner{
		Id:   sailpointV3.PtrString("222"),
		Type: sailpointV3.PtrString("compareInnerObject"),
	}
	curOwn := sailpointV3.SourceOwner{
		Id:   sailpointV3.PtrString("123"),
		Type: sailpointV3.PtrString("compareInnerObject"),
	}

	mod := sailpointV3.Source{Owner: modOwn}
	cur := sailpointV3.Source{Owner: curOwn}
	patcOp, err := NewSourcePatchBuilder(&mod, &cur).GenerateJsonPatch()
	if err != nil {
		t.Error(err)
	}

	if len(patcOp) != 1 {
		t.Error("Expected 1 patch operations but got ", strconv.Itoa(len(patcOp)))
	}
	assert := assert.New(t)
	assert.Equal("/owner/id", patcOp[0].Path)
	assert.Equal("replace", patcOp[0].Op)
	assert.Equal("222", *patcOp[0].Value.String)
}

func connectorAttributesPatch_expectedResult() []sailpointBeta.JsonPatchOperation {
	var innerArray []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"C\",\"D\"]"), &innerArray)
	var arrayChanged []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"a\",\"c\",\"c\",\"g\"]"), &arrayChanged)
	inner3Added := map[string]interface{}{
		"anotherParam":  "asd",
		"anotherParam2": true,
	}
	return []sailpointBeta.JsonPatchOperation{
		{
			Op:   "replace",
			Path: "/connectorAttributes/inner1/inner2/asd",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("test-update"),
			},
		},
		{
			Op:   "add",
			Path: "/connectorAttributes/inner1/inner2/innerArray",
			Value: &sailpointBeta.JsonPatchOperationValue{
				ArrayOfArrayInner: &innerArray,
			},
		},
		{
			Op:   "remove",
			Path: "/connectorAttributes/inner1/innerString",
		},
		{
			Op:   "remove",
			Path: "/connectorAttributes/inner1/innerString",
		},
		{
			Op:   "add",
			Path: "/connectorAttributes/intValue",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Int32: sailpointBeta.PtrInt32(123),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/boolValue",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Bool: sailpointBeta.PtrBool(true),
			},
		},
		{
			Op:   "remove",
			Path: "/connectorAttributes/oldStringValue",
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/arrayChanged",
			Value: &sailpointBeta.JsonPatchOperationValue{
				ArrayOfArrayInner: &arrayChanged,
			},
		},
		{
			Op:   "add",
			Path: "/connectorAttributes/inner1/inner3Added",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &inner3Added,
			},
		},
		{
			Op:   "remove",
			Path: "/connectorAttributes/inner1/inner4Removed",
		},
	}
}

func Test_Source_ConnectorAttributesPatch_ReplaceMultiple(t *testing.T) {
	modAttr := map[string]interface{}{
		"inner1": map[string]interface{}{
			"inner2": map[string]interface{}{
				"asd":        "test-update",
				"innerArray": []string{"C", "D"},
			},
			"inner3Added": map[string]interface{}{
				"anotherParam":  "asd",
				"anotherParam2": true,
			},
		},
		"intValue":     123,
		"boolValue":    true,
		"newStringVal": "newString",
		"array":        []string{"a", "b", "c"},
		"arrayChanged": []string{"a", "c", "c", "g"},
	}

	curAttr := map[string]interface{}{
		"inner1": map[string]interface{}{
			"inner2": map[string]interface{}{
				"asd": "test",
			},
			"innerString": "innerStringRemoved",
			"inner4Removed": map[string]interface{}{
				"removed": "removed",
			},
		},
		"boolValue":      false,
		"oldStringValue": "oldValue",
		"array":          []string{"a", "b", "c"},
		"arrayChanged":   []string{"a", "b", "c"},
	}

	mod := sailpointV3.Source{ConnectorAttributes: modAttr}
	cur := sailpointV3.Source{ConnectorAttributes: curAttr}
	patcOp, err := NewSourcePatchBuilder(&mod, &cur).GenerateJsonPatch()
	expectedResults := connectorAttributesPatch_expectedResult()
	assertResults(t, err, patcOp, expectedResults)
}

func MultipleAttributeAdd_expectedResult() []sailpointBeta.JsonPatchOperation {

	var features []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"AUTHENTICATE\"]"), &features)
	return []sailpointBeta.JsonPatchOperation{
		{
			Op:   "add",
			Path: "/name",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newName"),
			},
		},
		{
			Op:   "add",
			Path: "/description",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newDescription"),
			},
		},
		{
			Op:   "replace",
			Path: "/owner/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newOwner"),
			},
		},
		{
			Op:   "add",
			Path: "/cluster",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"type": "CLUSTER",
					"id":   "newCluster",
					"name": "newClusterName",
				},
			},
		},
		{
			Op:   "add",
			Path: "/accountCorrelationConfig",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"id": "newAccountCorrelationConfig",
				},
			},
		},
		{
			Op:   "add",
			Path: "/accountCorrelationRule",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"id": "newAccountCorrelationRule",
				},
			},
		},
		{
			Op:   "add",
			Path: "/managerCorrelationMapping",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"accountAttributeName":  "newAccountAttribute",
					"identityAttributeName": "newIdentityAttribute",
				},
			},
		},
		{
			Op:   "add",
			Path: "/managerCorrelationRule",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"id": "newManagerCorrelationRule",
				},
			},
		},
		{
			Op:   "replace",
			Path: "/features",
			Value: &sailpointBeta.JsonPatchOperationValue{
				ArrayOfArrayInner: &features,
			},
		},
		{
			Op:   "add",
			Path: "/deleteThreshold",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Int32: sailpointBeta.PtrInt32(99),
			},
		},
		{
			Op:   "add",
			Path: "/managementWorkgroup",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"id": "newManagementWorkgroup",
				},
			},
		},
		{
			Op:   "add",
			Path: "/beforeProvisioningRule",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"id": "newBeforeProvisioningRule",
				},
			},
		},
	}
}

func Test_Source_MultipleAttributeAdd(t *testing.T) {
	mod := sailpointV3.Source{
		Name:        "newName",
		Description: sailpointV3.PtrString("newDescription"),
		Owner: sailpointV3.SourceOwner{
			Id: sailpointV3.PtrString("newOwner"),
		},
		Cluster: *sailpointV3.NewNullableSourceCluster(&sailpointV3.SourceCluster{
			Type: "CLUSTER",
			Id:   "newCluster",
			Name: "newClusterName",
		}),
		AccountCorrelationConfig: *sailpointV3.NewNullableSourceAccountCorrelationConfig(&sailpointV3.SourceAccountCorrelationConfig{
			Id: sailpointV3.PtrString("newAccountCorrelationConfig"),
		}),
		AccountCorrelationRule: *sailpointV3.NewNullableSourceAccountCorrelationRule(&sailpointV3.SourceAccountCorrelationRule{
			Id: sailpointV3.PtrString("newAccountCorrelationRule"),
		}),
		ManagerCorrelationMapping: &sailpointV3.SourceManagerCorrelationMapping{
			AccountAttributeName:  sailpointV3.PtrString("newAccountAttribute"),
			IdentityAttributeName: sailpointV3.PtrString("newIdentityAttribute"),
		},
		ManagerCorrelationRule: *sailpointV3.NewNullableSourceManagerCorrelationRule(&sailpointV3.SourceManagerCorrelationRule{
			Id: sailpointV3.PtrString("newManagerCorrelationRule"),
		}),
		Features: []string{
			"AUTHENTICATE",
		},
		DeleteThreshold: sailpointV3.PtrInt32(99),
		ManagementWorkgroup: *sailpointV3.NewNullableSourceManagementWorkgroup(&sailpointV3.SourceManagementWorkgroup{
			Id: sailpointV3.PtrString("newManagementWorkgroup"),
		}),
		BeforeProvisioningRule: *sailpointV3.NewNullableSourceBeforeProvisioningRule(&sailpointV3.SourceBeforeProvisioningRule{
			Id: sailpointV3.PtrString("newBeforeProvisioningRule"),
		}),
	}
	cur := sailpointV3.Source{}
	patch, err := NewSourcePatchBuilder(&mod, &cur).GenerateJsonPatch()
	expectedResults := MultipleAttributeAdd_expectedResult()
	assertResults(t, err, patch, expectedResults)
}

func MultipleAttributeReplace_expectedResult() []sailpointBeta.JsonPatchOperation {
	var features []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"AUTHENTICATE\"]"), &features)
	return []sailpointBeta.JsonPatchOperation{
		{
			Op:   "replace",
			Path: "/name",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newName"),
			},
		},
		{
			Op:   "replace",
			Path: "/description",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newDescription"),
			},
		},
		{
			Op:   "replace",
			Path: "/owner/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newOwner"),
			},
		},
		{
			Op:   "replace",
			Path: "/cluster/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newCluster"),
			},
		},
		{
			Op:   "replace",
			Path: "/accountCorrelationConfig/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newAccountCorrelationConfig"),
			},
		},
		{
			Op:   "replace",
			Path: "/accountCorrelationRule/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newAccountCorrelationRule"),
			},
		},
		{
			Op:   "replace",
			Path: "/managerCorrelationMapping/accountAttributeName",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newAccountAttribute"),
			},
		},
		{
			Op:   "replace",
			Path: "/managerCorrelationMapping/identityAttributeName",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newIdentityAttribute"),
			},
		},
		{
			Op:   "replace",
			Path: "/managerCorrelationRule/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newManagerCorrelationRule"),
			},
		},

		{
			Op:   "replace",
			Path: "/features",
			Value: &sailpointBeta.JsonPatchOperationValue{
				ArrayOfArrayInner: &features,
			},
		},
		{
			Op:   "replace",
			Path: "/deleteThreshold",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Int32: sailpointBeta.PtrInt32(99),
			},
		},
		{
			Op:   "replace",
			Path: "/managementWorkgroup/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newManagementWorkgroup"),
			},
		},
		{
			Op:   "replace",
			Path: "/beforeProvisioningRule/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newBeforeProvisioningRule"),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/a",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("aa"),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/b",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Int32: sailpointBeta.PtrInt32(12),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/c",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("12.2"),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/d",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Bool: sailpointBeta.PtrBool(true),
			},
		},
	}
}

func Test_Source_MultipleAttributeReplace(t *testing.T) {
	mod := sailpointV3.Source{
		Name:        "newName",
		Description: sailpointV3.PtrString("newDescription"),
		Owner: sailpointV3.SourceOwner{
			Id: sailpointV3.PtrString("newOwner"),
		},
		Cluster: *sailpointV3.NewNullableSourceCluster(&sailpointV3.SourceCluster{
			Id: "newCluster",
		}),
		AccountCorrelationConfig: *sailpointV3.NewNullableSourceAccountCorrelationConfig(&sailpointV3.SourceAccountCorrelationConfig{
			Id: sailpointV3.PtrString("newAccountCorrelationConfig"),
		}),
		AccountCorrelationRule: *sailpointV3.NewNullableSourceAccountCorrelationRule(&sailpointV3.SourceAccountCorrelationRule{
			Id: sailpointV3.PtrString("newAccountCorrelationRule"),
		}),
		ManagerCorrelationMapping: &sailpointV3.SourceManagerCorrelationMapping{
			AccountAttributeName:  sailpointV3.PtrString("newAccountAttribute"),
			IdentityAttributeName: sailpointV3.PtrString("newIdentityAttribute"),
		},
		ManagerCorrelationRule: *sailpointV3.NewNullableSourceManagerCorrelationRule(&sailpointV3.SourceManagerCorrelationRule{
			Id: sailpointV3.PtrString("newManagerCorrelationRule"),
		}),
		Features: []string{
			"AUTHENTICATE",
		},
		DeleteThreshold: sailpointV3.PtrInt32(99),
		ManagementWorkgroup: *sailpointV3.NewNullableSourceManagementWorkgroup(&sailpointV3.SourceManagementWorkgroup{
			Id: sailpointV3.PtrString("newManagementWorkgroup"),
		}),
		BeforeProvisioningRule: *sailpointV3.NewNullableSourceBeforeProvisioningRule(&sailpointV3.SourceBeforeProvisioningRule{
			Id: sailpointV3.PtrString("newBeforeProvisioningRule"),
		}),
		ConnectorAttributes: map[string]interface{}{
			"a": "aa",
			"b": 12,
			"c": 12.2,
			"d": true,
		},
	}
	cur := sailpointV3.Source{
		Name:        "name",
		Description: sailpointV3.PtrString("description"),
		Owner: sailpointV3.SourceOwner{
			Id: sailpointV3.PtrString("owner"),
		},
		Cluster: *sailpointV3.NewNullableSourceCluster(&sailpointV3.SourceCluster{
			Id: "cluster",
		}),
		AccountCorrelationConfig: *sailpointV3.NewNullableSourceAccountCorrelationConfig(&sailpointV3.SourceAccountCorrelationConfig{
			Id: sailpointV3.PtrString("accountCorrelationConfig"),
		}),
		AccountCorrelationRule: *sailpointV3.NewNullableSourceAccountCorrelationRule(&sailpointV3.SourceAccountCorrelationRule{
			Id: sailpointV3.PtrString("accountCorrelationRule"),
		}),
		ManagerCorrelationMapping: &sailpointV3.SourceManagerCorrelationMapping{
			AccountAttributeName:  sailpointV3.PtrString("accountAttribute"),
			IdentityAttributeName: sailpointV3.PtrString("identityAttribute"),
		},
		ManagerCorrelationRule: *sailpointV3.NewNullableSourceManagerCorrelationRule(&sailpointV3.SourceManagerCorrelationRule{
			Id: sailpointV3.PtrString("managerCorrelationRule"),
		}),
		Features: []string{
			"DISCOVER_SCHEMA",
		},
		DeleteThreshold: sailpointV3.PtrInt32(1),
		ManagementWorkgroup: *sailpointV3.NewNullableSourceManagementWorkgroup(&sailpointV3.SourceManagementWorkgroup{
			Id: sailpointV3.PtrString("managementWorkgroup"),
		}),
		BeforeProvisioningRule: *sailpointV3.NewNullableSourceBeforeProvisioningRule(&sailpointV3.SourceBeforeProvisioningRule{
			Id: sailpointV3.PtrString("beforeProvisioningRule"),
		}),
		ConnectorAttributes: map[string]interface{}{
			"a": "a",
			"b": 11,
			"c": 11.1,
			"d": false,
		},
	}
	patch, err := NewSourcePatchBuilder(&mod, &cur).GenerateJsonPatch()
	expectedResults := MultipleAttributeReplace_expectedResult()
	assertResults(t, err, patch, expectedResults)
}

func assertResults(t *testing.T, err error, patcOp []sailpointBeta.JsonPatchOperation, expectedResults []sailpointBeta.JsonPatchOperation) {
	if err != nil {
		t.Error(err)
	}

	if len(patcOp) != len(expectedResults) {
		t.Error("Expected "+strconv.Itoa(len(expectedResults))+" patch operations but got ", strconv.Itoa(len(patcOp)))
	}
	assert := assert.New(t)
	for _, expected := range expectedResults {
		found := false
		for _, actual := range patcOp {
			if actual.Path == expected.Path {
				found = true
				assert.Equal(expected, actual)
			}
		}
		if !found {
			t.Errorf("Did not find expected patch operation %+v", expected)
		}
	}
}
