package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/nanicienta/nani-commons/pkg/features"
	"log"
)

type PostgresSelectFeature struct {
	features.BaseFeature
	DB     *sql.DB
	Query  string
	Params []interface{}
}

func (q *PostgresSelectFeature) Execute() {
	rows, err := q.DB.Query(q.Query, q.Params...)
	if err != nil {
		log.Fatalf("Error ejecutando la consulta: %v", err)
	}
	//TODO check this defer
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
			q.BaseFeature.Context[col] = values[i]
			fmt.Printf("%s: %v\n", col, values[i])
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatalf("Error en la lista de filas: %v", err)
	}
}
