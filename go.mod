module github.com/nanicienta/nani-engine

go 1.23.0

require (
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/nanicienta/nani-commons v1.0.0
)

replace github.com/nanicienta/nani-commons => ../nani-commons
