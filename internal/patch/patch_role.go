package patch

import (
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
)

type RolePatchBuilder struct {
	abstractPatchBuilder
	modified, current *sailpointV3.Role
}

func NewRolePatchBuilder(modified, current *sailpointV3.Role) *RolePatchBuilder {
	v := &RolePatchBuilder{
		modified: modified,
		current:  current,
	}
	v.abstractPatchBuilder = abstractPatchBuilder{}
	v.abstractPatchBuilder.defineValuesToCompare = v.defineValuesToCompare
	return v
}

func (pb *RolePatchBuilder) defineValuesToCompare() {
	pb.valuesToCompare = []comparableValues{
		{modifiedVal: pb.modified.Name, currentVal: pb.current.Name, path: "/name"},
		{modifiedVal: pb.modified.Description.Get(), currentVal: pb.current.Description.Get(), path: "/description"},
		{modifiedVal: pb.modified.AccessProfiles, currentVal: pb.current.AccessProfiles, path: "/accessProfiles"},
		{modifiedVal: pb.modified.Entitlements, currentVal: pb.current.Entitlements, path: "/entitlements"},
		{modifiedVal: pb.modified.Membership.Get(), currentVal: pb.current.Membership.Get(), path: "/membership"},
		{modifiedVal: pb.modified.Enabled, currentVal: pb.current.Enabled, path: "/enabled"},
		{modifiedVal: pb.modified.Requestable, currentVal: pb.current.Requestable, path: "/requestable"},
		{modifiedVal: pb.modified.AccessRequestConfig, currentVal: pb.current.AccessRequestConfig, path: "/accessRequestConfig"},
		{modifiedVal: pb.modified.RevocationRequestConfig, currentVal: pb.current.RevocationRequestConfig, path: "/revocationRequestConfig"},
		{modifiedVal: pb.modified.Segments, currentVal: pb.current.Segments, path: "/segments"},
	}

	pb.referencesToCompare = []comparableValues{
		{modifiedVal: pb.modified.Owner, currentVal: pb.current.Owner, path: "/owner"},
	}
}
