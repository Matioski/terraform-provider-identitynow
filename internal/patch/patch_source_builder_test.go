//go:build !integration

package patch

import (
	"encoding/json"
	"strconv"
	"testing"

	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
	"github.com/stretchr/testify/assert"
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
	var newArray []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"n\"]"), &newArray)
	var innerArray []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"C\",\"D\"]"), &innerArray)
	var arrayChanged []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"a\",\"c\",\"c\",\"g\"]"), &arrayChanged)
	inner3Added := map[string]interface{}{
		"anotherParam":  "asd",
		"anotherParam2": true,
	}
	inner3AddedValue := sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&inner3Added)
	return []sailpointBeta.JsonPatchOperation{
		{
			Op:   "replace",
			Path: "/connectorAttributes/inner1/inner2/asd",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("test-update"),
			},
		},
		{
			Op:   "add",
			Path: "/connectorAttributes/inner1/inner2/innerArray",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
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
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				Int32: sailpointBeta.PtrInt32(123),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/boolValue",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
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
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				ArrayOfArrayInner: &arrayChanged,
			},
		},
		{
			Op:    "add",
			Path:  "/connectorAttributes/inner1/inner3Added",
			Value: &inner3AddedValue,
		},
		{
			Op:   "remove",
			Path: "/connectorAttributes/inner1/inner4Removed",
		},
		{
			Op:   "add",
			Path: "/connectorAttributes/newArray",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				ArrayOfArrayInner: &newArray,
			},
		},
		{
			Op:   "add",
			Path: "/connectorAttributes/newValue",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("new"),
			},
		},
		{
			Op:   "remove",
			Path: "/connectorAttributes/removeArray",
		},
		{
			Op:   "remove",
			Path: "/connectorAttributes/removeValue",
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
		"removeArray":  nil,
		"removeValue":  nil,
		"newArray":     []string{"n"},
		"newValue":     "new",
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
		"removeArray":    []string{"o"},
		"removeValue":    "val",
		"newArray":       nil,
		"newValue":       nil,
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

	cluster := map[string]interface{}{
		"type": "CLUSTER",
		"id":   "newCluster",
		"name": "newClusterName",
	}
	clusterValue := sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&cluster)

	accountCorrConfig := map[string]interface{}{
		"id": "newAccountCorrelationConfig",
	}
	accountCorrConfigValue := sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&accountCorrConfig)

	accountCorrRule := map[string]interface{}{
		"id": "newAccountCorrelationRule",
	}
	accountCorrelationRuleValue := sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&accountCorrRule)

	managerCorrMapping := map[string]interface{}{
		"accountAttributeName":  "newAccountAttribute",
		"identityAttributeName": "newIdentityAttribute",
	}
	managerCorrMappingValue := sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&managerCorrMapping)

	managerCorrRule := map[string]interface{}{
		"id": "newManagerCorrelationRule",
	}
	managerCorrRuleValue := sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&managerCorrRule)

	anagementWorkgroup := map[string]interface{}{
		"id": "newManagementWorkgroup",
	}
	managementWorkgroupValue := sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&anagementWorkgroup)

	beforeProvisioningRule := map[string]interface{}{
		"id": "newBeforeProvisioningRule",
	}
	beforeProvisioningRulevalue := sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&beforeProvisioningRule)

	return []sailpointBeta.JsonPatchOperation{
		{
			Op:   "add",
			Path: "/name",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newName"),
			},
		},
		{
			Op:   "add",
			Path: "/description",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newDescription"),
			},
		},
		{
			Op:   "replace",
			Path: "/owner/id",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newOwner"),
			},
		},
		{
			Op:    "add",
			Path:  "/cluster",
			Value: &clusterValue,
		},
		{
			Op:    "add",
			Path:  "/accountCorrelationConfig",
			Value: &accountCorrConfigValue,
		},
		{
			Op:    "add",
			Path:  "/accountCorrelationRule",
			Value: &accountCorrelationRuleValue,
		},
		{
			Op:    "add",
			Path:  "/managerCorrelationMapping",
			Value: &managerCorrMappingValue,
		},
		{
			Op:    "add",
			Path:  "/managerCorrelationRule",
			Value: &managerCorrRuleValue,
		},
		{
			Op:   "replace",
			Path: "/features",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				ArrayOfArrayInner: &features,
			},
		},
		{
			Op:   "add",
			Path: "/deleteThreshold",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				Int32: sailpointBeta.PtrInt32(99),
			},
		},
		{
			Op:    "add",
			Path:  "/managementWorkgroup",
			Value: &managementWorkgroupValue,
		},
		{
			Op:    "add",
			Path:  "/beforeProvisioningRule",
			Value: &beforeProvisioningRulevalue,
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
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newName"),
			},
		},
		{
			Op:   "replace",
			Path: "/description",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newDescription"),
			},
		},
		{
			Op:   "replace",
			Path: "/owner/id",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newOwner"),
			},
		},
		{
			Op:   "replace",
			Path: "/cluster/id",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newCluster"),
			},
		},
		{
			Op:   "replace",
			Path: "/accountCorrelationConfig/id",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newAccountCorrelationConfig"),
			},
		},
		{
			Op:   "replace",
			Path: "/accountCorrelationRule/id",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newAccountCorrelationRule"),
			},
		},
		{
			Op:   "replace",
			Path: "/managerCorrelationMapping/accountAttributeName",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newAccountAttribute"),
			},
		},
		{
			Op:   "replace",
			Path: "/managerCorrelationMapping/identityAttributeName",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newIdentityAttribute"),
			},
		},
		{
			Op:   "replace",
			Path: "/managerCorrelationRule/id",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newManagerCorrelationRule"),
			},
		},

		{
			Op:   "replace",
			Path: "/features",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				ArrayOfArrayInner: &features,
			},
		},
		{
			Op:   "replace",
			Path: "/deleteThreshold",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				Int32: sailpointBeta.PtrInt32(99),
			},
		},
		{
			Op:   "replace",
			Path: "/managementWorkgroup/id",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newManagementWorkgroup"),
			},
		},
		{
			Op:   "replace",
			Path: "/beforeProvisioningRule/id",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newBeforeProvisioningRule"),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/a",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("aa"),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/b",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				Int32: sailpointBeta.PtrInt32(12),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/c",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("12.2"),
			},
		},
		{
			Op:   "replace",
			Path: "/connectorAttributes/d",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
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
