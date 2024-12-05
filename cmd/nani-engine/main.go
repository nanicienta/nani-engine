package main

import (
	"encoding/json"
	"github.com/nanicienta/nani-commons/pkg/model"
	"github.com/nanicienta/nani-engine/internal/services"
	"github.com/nanicienta/nani-engine/pkg/features/hola_mundo"
	"github.com/nanicienta/nani-engine/pkg/features/postgresql"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	file, err := os.Open("/Users/sebastianbogado/nanicienta/nani-engine/workflows/workflow_1.json")
	if err != nil {
		log.Fatalf("No se pudo abrir el archivo: %v", err)
	}
	defer file.Close()

	// Leer el contenido del archivo
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("No se pudo leer el archivo: %v", err)
	}
	var workflow model.Workflow
	if err := json.Unmarshal(byteValue, &workflow); err != nil {
		log.Fatalf("No se pudo decodificar el archivo: %v", err)
	}
	workflow.InitWorkflow()
	initFeatures()
	executor := services.GetNaniExecutorInstance()
	executor.ExecuteWorkflow(workflow)
}

func initFeatures() {
	postgresql.InitPostgresSelectFeature()
	hola_mundo.InitPostgresSelectFeature()

}
