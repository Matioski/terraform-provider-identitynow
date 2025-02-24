//go:build !integration

package patch

import (
	"encoding/json"
	"testing"

	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
)

func LifecycleState_ExpectedResult() []sailpointBeta.JsonPatchOperation {
	var accountActions []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[{\"action\":\"ENABLE\",\"sourceIds\": [\"source1\", \"source2\"]}]"), &accountActions)
	var accessProfileIds []sailpointBeta.ArrayInner
	json.Unmarshal([]byte("[\"access1\", \"access2\"]"), &accessProfileIds)
	emailNotifOpts := map[string]interface{}{
		"emailAddressList":    []interface{}{"email1", "email2"},
		"notifyAllAdmins":     true,
		"notifyManagers":      true,
		"notifySpecificUsers": true,
	}
	emailNotifOptsValue := sailpointBeta.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&emailNotifOpts)

	return []sailpointBeta.JsonPatchOperation{
		{
			Op:   "add",
			Path: "/description",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("newDescription"),
			},
		},
		{
			Op:   "replace",
			Path: "/enabled",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				Bool: sailpointBeta.PtrBool(true),
			},
		},
		{
			Op:   "replace",
			Path: "/name",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("nameUpd"),
			},
		},
		{
			Op:   "replace",
			Path: "/technicalName",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("techNameUpd"),
			},
		},
		{
			Op:    "add",
			Path:  "/emailNotificationOption",
			Value: &emailNotifOptsValue,
		},
		{
			Op:   "replace",
			Path: "/accountActions",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				ArrayOfArrayInner: &accountActions,
			},
		},
		{
			Op:   "replace",
			Path: "/accessProfileIds",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				ArrayOfArrayInner: &accessProfileIds,
			},
		},
		{
			Op:   "replace",
			Path: "/identityState",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("INACTIVE_SHORT_TERM"),
			},
		},
	}
}
func Test_LifecycleState(t *testing.T) {
	activeString := "ACTIVE"
	inactiveString := "INACTIVE_SHORT_TERM"

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
		IdentityState:    *sailpointV3.NewNullableString(&inactiveString),
	}
	cur := sailpointV3.LifecycleState{
		Name:          "name",
		TechnicalName: "techName",
		Enabled:       sailpointV3.PtrBool(false),
		IdentityState: *sailpointV3.NewNullableString(&activeString),
	}

	patch, err := NewLifecycleStatePatchBuilder(&mod, &cur).GenerateJsonPatch()
	expectedResults := LifecycleState_ExpectedResult()

	assertResults(t, err, patch, expectedResults)
}
