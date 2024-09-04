package patch

import (
	sailpointBeta "github.com/sailpoint-oss/golang-sdk/v2/api_beta"
)

var _ patchBuilder = &WorkflowPatchBuilder{}

type WorkflowPatchBuilder struct {
	abstractPatchBuilder
	modified, current *sailpointBeta.CreateWorkflowRequest
}

func NewWorkflowPatchBuilder(modified, current *sailpointBeta.CreateWorkflowRequest) *WorkflowPatchBuilder {
	v := &WorkflowPatchBuilder{
		modified: modified,
		current:  current,
	}
	v.abstractPatchBuilder = abstractPatchBuilder{}
	v.abstractPatchBuilder.defineValuesToCompare = v.defineValuesToCompare
	return v
}

func (pb *WorkflowPatchBuilder) defineValuesToCompare() {
	pb.valuesToCompare = []comparableValues{
		{modifiedVal: pb.modified.Name, currentVal: pb.current.Name, path: "/name"},
		{modifiedVal: pb.modified.Description, currentVal: pb.current.Description, path: "/description"},
		{modifiedVal: pb.modified.Enabled, currentVal: pb.current.Enabled, path: "/enabled"},
		{modifiedVal: pb.modified.Definition, currentVal: pb.current.Definition, path: "/definition"},
		{modifiedVal: pb.modified.Trigger, currentVal: pb.current.Trigger, path: "/trigger"},
	}
}
