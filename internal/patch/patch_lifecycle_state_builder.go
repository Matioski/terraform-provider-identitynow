package patch

import (
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
)

var _ patchBuilder = &LifecycleStatePatchBuilder{}

type LifecycleStatePatchBuilder struct {
	abstractPatchBuilder
	modified, current *sailpointV3.LifecycleState
}

func NewLifecycleStatePatchBuilder(modified, current *sailpointV3.LifecycleState) *LifecycleStatePatchBuilder {
	v := &LifecycleStatePatchBuilder{
		modified: modified,
		current:  current,
	}
	v.abstractPatchBuilder = abstractPatchBuilder{}
	v.abstractPatchBuilder.defineValuesToCompare = v.defineValuesToCompare
	return v
}

func (pb *LifecycleStatePatchBuilder) defineValuesToCompare() {
	pb.valuesToCompare = []comparableValues{
		{modifiedVal: pb.modified.Name, currentVal: pb.current.Name, path: "/name"},
		{modifiedVal: pb.modified.Enabled, currentVal: pb.current.Enabled, path: "/enabled"},
		{modifiedVal: pb.modified.TechnicalName, currentVal: pb.current.TechnicalName, path: "/technicalName"},
		{modifiedVal: pb.modified.Description, currentVal: pb.current.Description, path: "/description"},
		{modifiedVal: pb.modified.EmailNotificationOption, currentVal: pb.current.EmailNotificationOption, path: "/emailNotificationOption"},
		{modifiedVal: pb.modified.AccountActions, currentVal: pb.current.AccountActions, path: "/accountActions"},
		{modifiedVal: pb.modified.AccessProfileIds, currentVal: pb.current.AccessProfileIds, path: "/accessProfileIds"},
		{modifiedVal: pb.modified.IdentityState, currentVal: pb.current.IdentityState, path: "/identityState"},
	}
}
