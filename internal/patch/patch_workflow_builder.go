package patch

import (
	"fmt"

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
	// Base fields
	pb.valuesToCompare = []comparableValues{
		{modifiedVal: pb.modified.Name, currentVal: pb.current.Name, path: "/name"},
		{modifiedVal: pb.modified.Description, currentVal: pb.current.Description, path: "/description"},
		{modifiedVal: pb.modified.Enabled, currentVal: pb.current.Enabled, path: "/enabled"},
		{modifiedVal: pb.modified.Definition, currentVal: pb.current.Definition, path: "/definition"},
	}
	pb.referencesToCompare = []comparableValues{}

	// Trigger comparison
	newTriggerName := pb.getModifiedWfTriggerName()
	currentTriggerName := pb.getCurrentWfTriggerName()
	newTriggerDescription := pb.getModifiedWfTriggerDescription()
	currentTriggerDescription := pb.getCurrentWfTriggerDescription()
	newCronString := pb.getModifiedtWfTriggerCronString()
	currentCronString := pb.getCurrenttWfTriggerCronString()
	newTriggerId := pb.getModifiedWfTriggerId()
	currentTriggerId := pb.getCurrentWfTriggerId()
	newTriggerFilter := pb.getModifiedWfTriggerFilter()
	currentTriggerFilter := pb.getCurrentWfTriggerFilter()

	pb.valuesToCompare = append(pb.valuesToCompare, comparableValues{modifiedVal: pb.modified.Trigger, currentVal: pb.current.Trigger, path: "/trigger"})
	if pb.modified.Trigger != nil {
		pb.valuesToCompare = append(pb.valuesToCompare, comparableValues{modifiedVal: newTriggerId, currentVal: currentTriggerId, path: "/trigger/attributes/id"})
		pb.valuesToCompare = append(pb.valuesToCompare, comparableValues{modifiedVal: newTriggerFilter, currentVal: currentTriggerFilter, path: "/trigger/attributes/filter.$"})
		pb.valuesToCompare = append(pb.valuesToCompare, comparableValues{modifiedVal: newTriggerDescription, currentVal: currentTriggerDescription, path: "/trigger/attributes/description"})
		pb.valuesToCompare = append(pb.valuesToCompare, comparableValues{modifiedVal: newTriggerName, currentVal: currentTriggerName, path: "/trigger/attributes/name"})
		pb.valuesToCompare = append(pb.valuesToCompare, comparableValues{modifiedVal: newCronString, currentVal: currentCronString, path: "/trigger/attributes/cronString"})
	}
}

func (pb *WorkflowPatchBuilder) getModifiedWfTriggerName() string {
	result := ""
	if pb.modified.Trigger != nil && pb.modified.Trigger.Type == "EXTERNAL" && pb.modified.Trigger.GetAttributes().ExternalAttributes != nil && pb.modified.Trigger.GetAttributes().ExternalAttributes.Name != nil {
		result = *pb.modified.Trigger.GetAttributes().ExternalAttributes.Name
	}
	return result
}

func (pb *WorkflowPatchBuilder) getCurrentWfTriggerName() string {
	result := ""
	if pb.current.Trigger != nil && pb.current.Trigger.Type == "EXTERNAL" && pb.current.Trigger.GetAttributes().ExternalAttributes != nil && pb.current.Trigger.GetAttributes().ExternalAttributes.Name != nil {
		result = *pb.current.Trigger.GetAttributes().ExternalAttributes.Name
	}
	return result
}

func (pb *WorkflowPatchBuilder) getModifiedWfTriggerDescription() string {
	result := ""
	if pb.modified.Trigger != nil {
		if pb.modified.Trigger.Type == "EVENT" && pb.modified.Trigger.GetAttributes().EventAttributes != nil && pb.modified.Trigger.GetAttributes().EventAttributes.Description != nil {
			result = *pb.modified.Trigger.GetAttributes().EventAttributes.Description
		} else if pb.modified.Trigger.Type == "EXTERNAL" && pb.modified.Trigger.GetAttributes().ExternalAttributes != nil && pb.modified.Trigger.GetAttributes().ExternalAttributes.Description != nil {
			result = *pb.modified.Trigger.GetAttributes().ExternalAttributes.Description
		}
	}
	return result
}

func (pb *WorkflowPatchBuilder) getCurrentWfTriggerDescription() string {
	result := ""
	if pb.current.Trigger != nil {
		if pb.current.Trigger.Type == "EVENT" && pb.current.Trigger.GetAttributes().EventAttributes != nil && pb.current.Trigger.GetAttributes().EventAttributes.Description != nil {
			result = *pb.current.Trigger.GetAttributes().EventAttributes.Description
		} else if pb.current.Trigger.Type == "EXTERNAL" && pb.current.Trigger.GetAttributes().ExternalAttributes != nil && pb.current.Trigger.GetAttributes().ExternalAttributes.Description != nil {
			result = *pb.current.Trigger.GetAttributes().ExternalAttributes.Description
		}
	}
	return result
}

func (pb *WorkflowPatchBuilder) getModifiedtWfTriggerCronString() string {
	result := ""
	if pb.modified.Trigger != nil && pb.modified.Trigger.GetAttributes().ScheduledAttributes != nil && pb.modified.Trigger.GetAttributes().ScheduledAttributes.CronString != nil {
		fmt.Printf("*pb.modified.Trigger.GetAttributes().ScheduledAttributes.CronString = %s\n", *pb.modified.Trigger.GetAttributes().ScheduledAttributes.CronString)
		result = *pb.modified.Trigger.GetAttributes().ScheduledAttributes.CronString
	}
	return result
}

func (pb *WorkflowPatchBuilder) getCurrenttWfTriggerCronString() string {
	result := ""
	if pb.current.Trigger != nil && pb.current.Trigger.GetAttributes().ScheduledAttributes != nil && pb.current.Trigger.GetAttributes().ScheduledAttributes.CronString != nil {
		fmt.Printf("*pb.current.Trigger.GetAttributes().ScheduledAttributes.CronString = %s\n", *pb.current.Trigger.GetAttributes().ScheduledAttributes.CronString)
		result = *pb.current.Trigger.GetAttributes().ScheduledAttributes.CronString
	}
	return result
}

func (pb *WorkflowPatchBuilder) getModifiedWfTriggerId() string {
	result := ""
	if pb.modified.Trigger != nil && pb.modified.Trigger.Type == "EVENT" && pb.modified.Trigger.GetAttributes().EventAttributes != nil {
		result = pb.modified.Trigger.GetAttributes().EventAttributes.Id
	}
	return result
}

func (pb *WorkflowPatchBuilder) getCurrentWfTriggerId() string {
	result := ""
	if pb.current.Trigger != nil && pb.current.Trigger.Type == "EVENT" && pb.current.Trigger.GetAttributes().EventAttributes != nil {
		result = pb.current.Trigger.GetAttributes().EventAttributes.Id
	}
	return result
}

func (pb *WorkflowPatchBuilder) getModifiedWfTriggerFilter() string {
	result := ""
	if pb.modified.Trigger != nil && pb.modified.Trigger.Type == "EVENT" && pb.modified.Trigger.GetAttributes().EventAttributes != nil && pb.modified.Trigger.GetAttributes().EventAttributes.Filter != nil {
		result = *pb.modified.Trigger.GetAttributes().EventAttributes.Filter
	}
	return result
}

func (pb *WorkflowPatchBuilder) getCurrentWfTriggerFilter() string {
	result := ""
	if pb.current.Trigger != nil && pb.current.Trigger.Type == "EVENT" && pb.current.Trigger.GetAttributes().EventAttributes != nil && pb.current.Trigger.GetAttributes().EventAttributes.Filter != nil {
		result = *pb.current.Trigger.GetAttributes().EventAttributes.Filter
	}
	return result
}
