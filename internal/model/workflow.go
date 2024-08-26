package model

import "github.com/nanicienta/nani-commons/pkg/constants"

type Step struct {
	StepId      string                `json:"stepId"`
	StepName    string                `json:"stepName"`
	FeatureName constants.FeatureName `json:"stepType"`
}

type Workflow struct {
	WorkflowId          string `json:"workflowId"`
	WorkflowName        string `json:"workflowName"`
	WorkflowDescription string `json:"workflowDescription"`
	Steps               []Step `json:"steps"`
}
