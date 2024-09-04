package patch

import (
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
)

type IdentityProfilePatchBuilder struct {
	abstractPatchBuilder
	modified, current *sailpointBeta.IdentityProfile
}

func NewIdentityProfilePatchBuilder(modified, current *sailpointBeta.IdentityProfile) *IdentityProfilePatchBuilder {
	v := &IdentityProfilePatchBuilder{
		modified: modified,
		current:  current,
	}
	v.abstractPatchBuilder = abstractPatchBuilder{}
	v.abstractPatchBuilder.defineValuesToCompare = v.defineValuesToCompare
	return v
}

func (pb *IdentityProfilePatchBuilder) defineValuesToCompare() {
	pb.valuesToCompare = []comparableValues{
		{modifiedVal: pb.modified.Name, currentVal: pb.current.Name, path: "/name"},
		{modifiedVal: pb.modified.Description.Get(), currentVal: pb.current.Description.Get(), path: "/description"},
		{modifiedVal: pb.modified.Priority, currentVal: pb.current.Priority, path: "/priority"},
		{modifiedVal: pb.modified.IdentityRefreshRequired, currentVal: pb.current.IdentityRefreshRequired, path: "/identityRefreshRequired"},
		{modifiedVal: pb.modified.HasTimeBasedAttr, currentVal: pb.current.HasTimeBasedAttr, path: "/hasTimeBasedAttr"},
		{modifiedVal: pb.modified.IdentityAttributeConfig, currentVal: pb.current.IdentityAttributeConfig, path: "/identityAttributeConfig"},
	}

	pb.referencesToCompare = []comparableValues{
		{modifiedVal: pb.modified.Owner, currentVal: pb.current.Owner, path: "/owner"},
		{modifiedVal: &pb.modified.AuthoritativeSource, currentVal: &pb.current.AuthoritativeSource, path: "/authoritativeSource"},
	}
}
