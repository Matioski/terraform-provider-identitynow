package patch

import (
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
)

var _ patchBuilder = &OrgConfigPatchBuilder{}

type OrgConfigPatchBuilder struct {
	abstractPatchBuilder
	modified, current *sailpointBeta.OrgConfig
}

func NewOrgConfigPatchBuilder(modified, current *sailpointBeta.OrgConfig) *OrgConfigPatchBuilder {
	v := &OrgConfigPatchBuilder{
		modified: modified,
		current:  current,
	}
	v.abstractPatchBuilder = abstractPatchBuilder{}
	v.abstractPatchBuilder.defineValuesToCompare = v.defineValuesToCompare
	return v
}

func (pb *OrgConfigPatchBuilder) defineValuesToCompare() {
	pb.valuesToCompare = []comparableValues{
		{modifiedVal: pb.modified.TimeZone, currentVal: pb.current.TimeZone, path: "/timeZone"},
	}
}
