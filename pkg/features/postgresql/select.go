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

func (q *PostgresSelectFeature) Execute(node model.Node, workflow model.Workflow) {
	found, queryField := node.GetField(field_id.PostgresQueryField)
	if !found {
		//TODO create new error
		_ = fmt.Errorf("query field not found")
	}
	connectionString := "user=postgres dbname=mydatabase sslmode=disable password=example"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error abriendo la conexión: %v", err)
	}
	rows, err := db.Query(queryField.Value)
	if err != nil {
		log.Fatalf("Error ejecutando la consulta: %v", err)
	}
	////TODO check this defer
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Error obteniendo las columnas: %v", err)
	}
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatalf("Error escaneando la fila: %v", err)
		}

		for i, col := range columns {
			workflow.Context[col] = values[i]
			fmt.Printf("%s: %v\n", col, values[i])
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("Error en la lista de filas: %v", err)
	}
}
