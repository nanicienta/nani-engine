package services

import (
	"fmt"
	"github.com/nanicienta/nani-commons/pkg/features"
	"github.com/nanicienta/nani-commons/pkg/model"
	"sync"
)

type NaniExecutorService struct {
	features map[string]features.Feature
}

var instance *NaniExecutorService
var once sync.Once

func GetNaniExecutorInstance() *NaniExecutorService {
	once.Do(
		func() {
			instance = &NaniExecutorService{
				features: make(map[string]features.Feature),
			}
			fmt.Println("Initializing nani engine singleton")
		},
	)
	return instance
}
func (nes *NaniExecutorService) RegisterFeature(feature features.Feature) {
	nes.features[feature.GetInternalName()] = feature
}

func (nes *NaniExecutorService) ExecuteWorkflow(workflow model.Workflow) {
	nextStep := workflow.Init
	nextStepNode := workflow.GetNode(nextStep)
	found, feature := nes.getFeature(nextStepNode.Type)
	if !found {
		//TODO create new error
		_ = fmt.Errorf("feature not found for node type: %s", nextStepNode.Type)
	}
	feature.Execute(
		nextStepNode,
		workflow,
	) //TODO here should I create a context to send all things necesary
}

func (nes *NaniExecutorService) getFeature(featureType string) (bool, features.Feature) {
	feature, exists := nes.features[featureType]
	return exists, feature
}
