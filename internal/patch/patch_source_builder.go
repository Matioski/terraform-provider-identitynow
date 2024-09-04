package patch

import (
	sailpointV3 "github.com/sailpoint-oss/golang-sdk/v2/api_v3"
)

var _ patchBuilder = &SourcePatchBuilder{}

type SourcePatchBuilder struct {
	abstractPatchBuilder
	modified, current *sailpointV3.Source
}

func NewSourcePatchBuilder(modified, current *sailpointV3.Source) *SourcePatchBuilder {
	v := &SourcePatchBuilder{
		modified: modified,
		current:  current,
	}
	v.abstractPatchBuilder = abstractPatchBuilder{}
	v.abstractPatchBuilder.defineValuesToCompare = v.defineValuesToCompare
	return v
}

func (pb *SourcePatchBuilder) defineValuesToCompare() {
	pb.valuesToCompare = []comparableValues{
		{modifiedVal: pb.modified.Name, currentVal: pb.current.Name, path: "/name"},
		{modifiedVal: pb.modified.Description, currentVal: pb.current.Description, path: "/description"},
		{modifiedVal: pb.modified.Features, currentVal: pb.current.Features, path: "/features"},
		{modifiedVal: pb.modified.ConnectorAttributes, currentVal: pb.current.ConnectorAttributes, path: "/connectorAttributes"},
		{modifiedVal: pb.modified.DeleteThreshold, currentVal: pb.current.DeleteThreshold, path: "/deleteThreshold"},
		{modifiedVal: pb.modified.ManagerCorrelationMapping, currentVal: pb.current.ManagerCorrelationMapping, path: "/managerCorrelationMapping"},
	}

	pb.referencesToCompare = []comparableValues{
		{modifiedVal: &pb.modified.Owner, currentVal: &pb.current.Owner, path: "/owner"},
		{modifiedVal: pb.modified.Cluster, currentVal: pb.current.Cluster, path: "/cluster"},
		{modifiedVal: pb.modified.AccountCorrelationConfig, currentVal: pb.current.AccountCorrelationConfig, path: "/accountCorrelationConfig"},
		{modifiedVal: pb.modified.AccountCorrelationRule, currentVal: pb.current.AccountCorrelationRule, path: "/accountCorrelationRule"},
		{modifiedVal: pb.modified.ManagerCorrelationRule, currentVal: pb.current.ManagerCorrelationRule, path: "/managerCorrelationRule"},
		{modifiedVal: pb.modified.BeforeProvisioningRule, currentVal: pb.current.BeforeProvisioningRule, path: "/beforeProvisioningRule"},
		{modifiedVal: pb.modified.ManagementWorkgroup, currentVal: pb.current.ManagementWorkgroup, path: "/managementWorkgroup"},
	}
}
