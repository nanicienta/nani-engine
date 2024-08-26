package main

import (
	"encoding/json"
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
	var workflow map[string]interface{}
	if err := json.Unmarshal(byteValue, &workflow); err != nil {
		log.Fatalf("No se pudo decodificar el archivo: %v", err)
	}

}
