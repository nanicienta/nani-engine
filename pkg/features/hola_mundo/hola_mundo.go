package hola_mundo

import (
	"fmt"
	"github.com/nanicienta/nani-commons/pkg/constants/features"
	"github.com/nanicienta/nani-commons/pkg/model"
	"github.com/nanicienta/nani-engine/internal/services"
	"sync"
)

type HolaMundoFeature struct {
}

var instance *HolaMundoFeature
var once sync.Once

func InitPostgresSelectFeature() *HolaMundoFeature {
	once.Do(
		func() {
			instance = &HolaMundoFeature{}
			fmt.Println("Initializing hola mundo feature")
		},
	)
	engine := services.GetNaniExecutorInstance()
	engine.RegisterFeature(instance)
	return instance
}

func (q *HolaMundoFeature) Execute(
	node model.Node,
	workflow model.Workflow,
	previousPayload map[string]interface{},
) (map[string]interface{}, string, error) {
	fmt.Print("Hola Mundo")
	return nil, "", nil
}

func (q *HolaMundoFeature) GetInternalName() string {
	return features.HolaMundo
}
