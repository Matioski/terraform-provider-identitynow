//go:build !integration

package patch

import (
	"encoding/json"
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
	"testing"
)

func LifecycleState_ExpectedResult() []sailpointBeta.JsonPatchOperation {
	var accountActions []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[{\"action\":\"ENABLE\",\"sourceIds\": [\"source1\", \"source2\"]}]"), &accountActions)
	var accessProfileIds []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"access1\", \"access2\"]"), &accessProfileIds)
	return []sailpointBeta.JsonPatchOperation{
		{
			Op:   "add",
			Path: "/description",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("newDescription"),
			},
		},
		{
			Op:   "replace",
			Path: "/enabled",
			Value: &sailpointBeta.JsonPatchOperationValue{
				Bool: sailpointBeta.PtrBool(true),
			},
		},
		{
			Op:   "replace",
			Path: "/name",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("nameUpd"),
			},
		},
		{
			Op:   "replace",
			Path: "/technicalName",
			Value: &sailpointBeta.JsonPatchOperationValue{
				String: sailpointBeta.PtrString("techNameUpd"),
			},
		},
		{
			Op:   "add",
			Path: "/emailNotificationOption",
			Value: &sailpointBeta.JsonPatchOperationValue{
				MapmapOfStringinterface: &map[string]interface{}{
					"emailAddressList":    []interface{}{"email1", "email2"},
					"notifyAllAdmins":     true,
					"notifyManagers":      true,
					"notifySpecificUsers": true,
				},
			},
		},
		{
			Op:   "replace",
			Path: "/accountActions",
			Value: &sailpointBeta.JsonPatchOperationValue{
				ArrayOfArrayInner: &accountActions,
			},
		},
		{
			Op:   "replace",
			Path: "/accessProfileIds",
			Value: &sailpointBeta.JsonPatchOperationValue{
				ArrayOfArrayInner: &accessProfileIds,
			},
		},
	}
}
func Test_LifecycleState(t *testing.T) {
	mod := sailpointV3.LifecycleState{
		Name:          "nameUpd",
		Enabled:       sailpointV3.PtrBool(true),
		TechnicalName: "techNameUpd",
		Description:   sailpointV3.PtrString("newDescription"),

		EmailNotificationOption: &sailpointV3.EmailNotificationOption{
			NotifyManagers:      sailpointV3.PtrBool(true),
			NotifyAllAdmins:     sailpointV3.PtrBool(true),
			NotifySpecificUsers: sailpointV3.PtrBool(true),
			EmailAddressList:    []string{"email1", "email2"},
		},
		AccountActions: []sailpointV3.AccountAction{
			{
				Action:    sailpointV3.PtrString("ENABLE"),
				SourceIds: []string{"source1", "source2"},
			},
		},
		AccessProfileIds: []string{"access1", "access2"},
	}
	cur := sailpointV3.LifecycleState{
		Name:          "name",
		TechnicalName: "techName",
		Enabled:       sailpointV3.PtrBool(false),
	}

	patch, err := NewLifecycleStatePatchBuilder(&mod, &cur).GenerateJsonPatch()
	expectedResults := LifecycleState_ExpectedResult()

	assertResults(t, err, patch, expectedResults)
}
