package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	featureConstants "github.com/nanicienta/nani-commons/pkg/constants/features"
	"github.com/nanicienta/nani-commons/pkg/constants/field_id"
	"github.com/nanicienta/nani-commons/pkg/model"
	"github.com/nanicienta/nani-engine/internal/services"
	"log"
	"sync"
	"time"
)

type PostgresSelectFeature struct {
}

var instance *PostgresSelectFeature
var once sync.Once

func InitPostgresSelectFeature() *PostgresSelectFeature {
	once.Do(
		func() {
			instance = &PostgresSelectFeature{}
			fmt.Println("Initializing postgres select feature")
		},
	)
	engine := services.GetNaniExecutorInstance()
	engine.RegisterFeature(instance)
	return instance
}

func (q *PostgresSelectFeature) GetInternalName() string {
	return featureConstants.PostgresSelect
}

func (q *PostgresSelectFeature) Execute(
	node model.Node,
	workflow model.Workflow,
	previousPayload map[string]interface{},
) (map[string]interface{}, string, error) {
	start := time.Now()

	found, queryField := node.GetField(field_id.PostgresQueryField)
	if !found {
		return nil, "", fmt.Errorf("query field not found")
	}

	found, connection := workflow.GetConnection(node.ConnectionId)
	if !found {
		return nil, "", fmt.Errorf("connection not found")
	}

	connectionString := fmt.Sprintf(
		`user=%s dbname=%s sslmode=%s password=%s`,
		connection.User, connection.Database,
		connection.Ssl, connection.Password,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error abriendo la conexión: %v", err)
	}
	defer db.Close()

	rows, err := db.Query(queryField.Value)
	if err != nil {
		log.Fatalf("Error ejecutando la consulta: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Error obteniendo las columnas: %v", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Crear un contenedor para cada fila
		row := make(map[string]interface{})

		// Crear un slice de interfaces para las columnas
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		// Asignar los punteros
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Escanear la fila actual en los punteros
		err = rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatalf("Error escaneando la fila: %v", err)
		}

		// Copiar los valores en el mapa de la fila
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}

		// Añadir la fila a los resultados
		results = append(results, row)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error en la lista de filas: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("La consulta tomó %s\n", elapsed)

	if node.StoreResultInVariable {
		if node.VariableName == "" {
			log.Fatalf("Error. Hay que especificar el nombre de la variable: %v")
		}
		workflow.Context[node.VariableName] = results
	}

	return map[string]interface{}{"results": results}, node.Links[0].NodeId, nil
}
