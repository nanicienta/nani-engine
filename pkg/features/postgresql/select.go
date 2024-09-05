package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nanicienta/nani-commons/pkg/constants/features"
	"github.com/nanicienta/nani-commons/pkg/constants/field_id"
	"github.com/nanicienta/nani-commons/pkg/model"
	"github.com/nanicienta/nani-commons/pkg/model/connections"
	"github.com/nanicienta/nani-engine/internal/services"
	"log"
	"sync"
	"time"
)

type PostgresSelectFeature struct{}

var instance *PostgresSelectFeature
var once sync.Once

// InitPostgresSelectFeature initializes the PostgresSelectFeature structure
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

// GetInternalName returns the feature name
func (q *PostgresSelectFeature) GetInternalName() string {
	return features.PostgresSelect
}

// Execute runs the query and processes the result
func (q *PostgresSelectFeature) Execute(
	node model.Node,
	workflow model.Workflow,
	previousPayload map[string]interface{},
) (map[string]interface{}, string, error) {
	start := time.Now()

	queryField, err := q.getQueryField(node)
	if err != nil {
		return nil, "", err
	}

	connection, err := q.getConnection(workflow, node)
	if err != nil {
		return nil, "", err
	}

	connectionString := q.buildConnectionString(connection)
	db, err := q.connectToDatabase(connectionString)
	if err != nil {
		return nil, "", err
	}
	defer db.Close()

	rows, err := q.executeQuery(db, queryField.Value, make(map[string]interface{}))
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	results, err := q.processRows(rows)
	if err != nil {
		return nil, "", err
	}

	elapsed := time.Since(start)
	fmt.Printf("La consulta tomó %s\n", elapsed)

	q.storeResultsInVariable(workflow, node, results)

	return map[string]interface{}{"results": results}, node.Links[0].NodeId, nil
}

func (q *PostgresSelectFeature) getQueryField(node model.Node) (*model.Field, error) {
	found, queryField := node.GetField(field_id.PostgresQueryField)
	if !found {
		return nil, fmt.Errorf("query field not found")
	}
	return &queryField, nil
}

func (q *PostgresSelectFeature) getConnection(
	workflow model.Workflow,
	node model.Node,
) (*connections.BaseConnection, error) {
	found, connection := workflow.GetConnection(node.ConnectionId)
	if !found {
		return nil, fmt.Errorf("connection not found")
	}
	return &connection, nil
}

func (q *PostgresSelectFeature) buildConnectionString(connection *connections.BaseConnection) string {
	return fmt.Sprintf(
		`user=%s dbname=%s sslmode=%s password=%s`,
		connection.User, connection.Database,
		connection.Ssl, connection.Password,
	)
}

func (q *PostgresSelectFeature) connectToDatabase(connectionString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error abriendo la conexión: %v", err)
	}
	return db, nil
}

func (q *PostgresSelectFeature) executeQuery(
	db *sqlx.DB,
	query string,
	arguments map[string]interface{},
) (*sqlx.Rows, error) {
	rows, err := db.NamedQuery(query, arguments)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
	}
	return rows, nil
}

func (q *PostgresSelectFeature) processRows(rows *sqlx.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo las columnas: %v", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		row, err := q.processRow(columns, rows)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error en la lista de filas: %v", err)
	}

	return results, nil
}

func (q *PostgresSelectFeature) processRow(
	columns []string,
	rows *sqlx.Rows,
) (map[string]interface{}, error) {
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	err := rows.Scan(valuePtrs...)
	if err != nil {
		return nil, fmt.Errorf("error escaneando la fila: %v", err)
	}

	row := make(map[string]interface{})
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

	return row, nil
}

func (q *PostgresSelectFeature) storeResultsInVariable(
	workflow model.Workflow,
	node model.Node,
	results []map[string]interface{},
) {
	if node.StoreResultInVariable {
		if node.VariableName == "" {
			log.Fatalf("Error. Hay que especificar el nombre de la variable")
		}
		workflow.Context[node.VariableName] = results
	}
}
