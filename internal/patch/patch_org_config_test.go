//go:build !integration

package patch

import (
	"testing"

	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
)

func OrgConfig_ExpectedResult() []sailpointBeta.JsonPatchOperation {
	return []sailpointBeta.JsonPatchOperation{
		{
			Op:   "replace",
			Path: "/timeZone",
			Value: &sailpointBeta.UpdateMultiHostSourcesRequestInnerValue{
				String: sailpointBeta.PtrString("UTC"),
			},
		},
	}
}
func Test_OrgConfig(t *testing.T) {
	mod := sailpointBeta.OrgConfig{
		TimeZone: sailpointBeta.PtrString("UTC"),
	}
	cur := sailpointBeta.OrgConfig{
		TimeZone: sailpointBeta.PtrString("Europe/Zurich"),
	}

	patch, err := NewOrgConfigPatchBuilder(&mod, &cur).GenerateJsonPatch()
	expectedResults := OrgConfig_ExpectedResult()

	assertResults(t, err, patch, expectedResults)
}
