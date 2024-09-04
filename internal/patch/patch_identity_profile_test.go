//go:build !integration

package patch

import (
	"encoding/json"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	"testing"
)

func IdentityProfile_MultipleAttributeAdd_ExpectedResult() []sailpointBeta.JsonPatchOperation {
	test := map[string]interface{}{
		"identityAttributeName": "newIdentityAttributeName",
		"transformDefinition": map[string]interface{}{
			"type": "newType",
			"attributes": map[string]interface{}{
				"test": "test",
			},
		},
	}
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
			Op:   "add",
			Path: "/owner",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"id": "newOwner",
				},
			},
		},
		{
			Op:   "replace",
			Path: "/authoritativeSource/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newAuthoritativeSource"),
			},
		},
		{
			Op:   "add",
			Path: "/priority",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Int32: sailpointBeta.PtrInt32(99),
			},
		},
		{
			Op:   "add",
			Path: "/identityRefreshRequired",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Bool: sailpointBeta.PtrBool(true),
			},
		},
		{
			Op:   "add",
			Path: "/hasTimeBasedAttr",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Bool: sailpointBeta.PtrBool(true),
			},
		},
		{
			Op:   "add",
			Path: "/identityAttributeConfig",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"enabled":             true,
					"attributeTransforms": []interface{}{test},
				},
			},
		},
	}
}
func Test_IdentityProfile_MultipleAttributeAdd(t *testing.T) {
	mod := sailpointBeta.IdentityProfile{
		Name:        "newName",
		Description: *sailpointBeta.NewNullableString(sailpointBeta.PtrString("newDescription")),
		Owner: *sailpointBeta.NewNullableIdentityProfileAllOfOwner(&sailpointBeta.IdentityProfileAllOfOwner{
			Id: sailpointBeta.PtrString("newOwner"),
		}),
		Priority:                sailpointBeta.PtrInt64(99),
		IdentityRefreshRequired: sailpointBeta.PtrBool(true),
		HasTimeBasedAttr:        sailpointBeta.PtrBool(true),
		IdentityAttributeConfig: &sailpointBeta.IdentityAttributeConfig{
			Enabled: sailpointBeta.PtrBool(true),
			AttributeTransforms: []sailpointBeta.IdentityAttributeTransform{{
				IdentityAttributeName: sailpointBeta.PtrString("newIdentityAttributeName"),
				TransformDefinition: &sailpointBeta.TransformDefinition{
					Type: sailpointBeta.PtrString("newType"),
					Attributes: map[string]interface{}{
						"test": "test",
					},
				},
			},
			},
		},
		AuthoritativeSource: sailpointBeta.IdentityProfileAllOfAuthoritativeSource{Id: sailpointBeta.PtrString("newAuthoritativeSource")},
	}
	cur := sailpointBeta.IdentityProfile{}

	patch, err := NewIdentityProfilePatchBuilder(&mod, &cur).GenerateJsonPatch()
	expectedResults := IdentityProfile_MultipleAttributeAdd_ExpectedResult()

	assertResults(t, err, patch, expectedResults)
}

func IdentityProfile_MultipleAttributeReplace_ExpectedResult() []sailpointBeta.JsonPatchOperation {
	var attributeTransforms []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[{\"identityAttributeName\":\"newIdentityAttributeName\",\"transformDefinition\":{\"type\":\"newType\",\"attributes\":{\"test\":\"test\"}}}]"), &attributeTransforms)
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
			Path: "/authoritativeSource/id",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newAuthoritativeSource"),
			},
		},
		{
			Op:   "replace",
			Path: "/priority",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Int32: sailpointBeta.PtrInt32(99),
			},
		},
		{
			Op:   "replace",
			Path: "/identityRefreshRequired",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Bool: sailpointBeta.PtrBool(true),
			},
		},
		{
			Op:   "replace",
			Path: "/hasTimeBasedAttr",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Bool: sailpointBeta.PtrBool(true),
			},
		},
		{
			Op:   "replace",
			Path: "/identityAttributeConfig/enabled",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Bool: sailpointBeta.PtrBool(true),
			},
		},
		{
			Op:   "replace",
			Path: "/identityAttributeConfig/attributeTransforms",
			Value: &sailpointBeta.JsonPatchOperationValue{
				ArrayOfArrayInner: &attributeTransforms,
			},
		},
	}
}

func Test_IdentityProfile_MultipleAttributeReplace(t *testing.T) {
	mod := sailpointBeta.IdentityProfile{
		Name:                    "newName",
		Description:             *sailpointBeta.NewNullableString(sailpointBeta.PtrString("newDescription")),
		Owner:                   *sailpointBeta.NewNullableIdentityProfileAllOfOwner(&sailpointBeta.IdentityProfileAllOfOwner{Id: sailpointBeta.PtrString("newOwner")}),
		Priority:                sailpointBeta.PtrInt64(99),
		IdentityRefreshRequired: sailpointBeta.PtrBool(true),
		HasTimeBasedAttr:        sailpointBeta.PtrBool(true),
		IdentityAttributeConfig: &sailpointBeta.IdentityAttributeConfig{
			Enabled: sailpointBeta.PtrBool(true),
			AttributeTransforms: []sailpointBeta.IdentityAttributeTransform{{
				IdentityAttributeName: sailpointBeta.PtrString("newIdentityAttributeName"),
				TransformDefinition: &sailpointBeta.TransformDefinition{
					Type: sailpointBeta.PtrString("newType"),
					Attributes: map[string]interface{}{
						"test": "test",
					},
				},
			},
			},
		},
		AuthoritativeSource: sailpointBeta.IdentityProfileAllOfAuthoritativeSource{Id: sailpointBeta.PtrString("newAuthoritativeSource")},
	}
	cur := sailpointBeta.IdentityProfile{
		Name:                    "Name",
		Description:             *sailpointBeta.NewNullableString(sailpointBeta.PtrString("Description")),
		Owner:                   *sailpointBeta.NewNullableIdentityProfileAllOfOwner(&sailpointBeta.IdentityProfileAllOfOwner{Id: sailpointBeta.PtrString("Owner")}),
		Priority:                sailpointBeta.PtrInt64(1),
		IdentityRefreshRequired: sailpointBeta.PtrBool(false),
		HasTimeBasedAttr:        sailpointBeta.PtrBool(false),
		IdentityAttributeConfig: &sailpointBeta.IdentityAttributeConfig{
			Enabled: sailpointBeta.PtrBool(false),
			AttributeTransforms: []sailpointBeta.IdentityAttributeTransform{{
				IdentityAttributeName: sailpointBeta.PtrString("dentityAttributeName"),
				TransformDefinition: &sailpointBeta.TransformDefinition{
					Type: sailpointBeta.PtrString("Type"),
					Attributes: map[string]interface{}{
						"test": "notTest",
					},
				},
			},
			},
		},
		AuthoritativeSource: sailpointBeta.IdentityProfileAllOfAuthoritativeSource{Id: sailpointBeta.PtrString("AuthoritativeSource")},
	}

	patch, err := NewIdentityProfilePatchBuilder(&mod, &cur).GenerateJsonPatch()
	expectedResults := IdentityProfile_MultipleAttributeReplace_ExpectedResult()

	assertResults(t, err, patch, expectedResults)
}

func IdentityProfile_MultipleAttributeRemove_ExpectedResult() []sailpointBeta.JsonPatchOperation {
	var attributeTransforms []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[]"), &attributeTransforms)
	return []sailpointBeta.JsonPatchOperation{
		{
			Op:   "remove",
			Path: "/name",
		},
		{
			Op:   "remove",
			Path: "/description",
		},
		{
			Op:   "remove",
			Path: "/owner",
		},
		{
			Op:   "remove",
			Path: "/priority",
		},
		{
			Op:   "remove",
			Path: "/identityRefreshRequired",
		},
		{
			Op:   "remove",
			Path: "/hasTimeBasedAttr",
		},
		{
			Op:   "replace",
			Path: "/identityAttributeConfig/attributeTransforms",
			Value: &sailpointBeta.JsonPatchOperationValue{
				ArrayOfArrayInner: &attributeTransforms,
			},
		},
	}
}

func Test_IdentityProfile_MultipleAttributeRemove(t *testing.T) {
	mod := sailpointBeta.IdentityProfile{
		IdentityAttributeConfig: &sailpointBeta.IdentityAttributeConfig{
			Enabled:             sailpointBeta.PtrBool(false),
			AttributeTransforms: make([]sailpointBeta.IdentityAttributeTransform, 0),
		},
	}
	cur := sailpointBeta.IdentityProfile{
		Name:                    "Name",
		Description:             *sailpointBeta.NewNullableString(sailpointBeta.PtrString("Description")),
		Owner:                   *sailpointBeta.NewNullableIdentityProfileAllOfOwner(&sailpointBeta.IdentityProfileAllOfOwner{Id: sailpointBeta.PtrString("Owner")}),
		Priority:                sailpointBeta.PtrInt64(1),
		IdentityRefreshRequired: sailpointBeta.PtrBool(false),
		HasTimeBasedAttr:        sailpointBeta.PtrBool(false),
		IdentityAttributeConfig: &sailpointBeta.IdentityAttributeConfig{
			Enabled: sailpointBeta.PtrBool(false),
			AttributeTransforms: []sailpointBeta.IdentityAttributeTransform{{
				IdentityAttributeName: sailpointBeta.PtrString("dentityAttributeName"),
				TransformDefinition: &sailpointBeta.TransformDefinition{
					Type: sailpointBeta.PtrString("Type"),
					Attributes: map[string]interface{}{
						"test": "notTest",
					},
				},
			},
			},
		},
	}

	patch, err := NewIdentityProfilePatchBuilder(&mod, &cur).GenerateJsonPatch()
	expectedResults := IdentityProfile_MultipleAttributeRemove_ExpectedResult()

	assertResults(t, err, patch, expectedResults)
}
