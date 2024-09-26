module example.com/go-api-tech-challenge

go 1.23.0

replace example.com/controller => ./pkg/controller

replace example.com/dbConn => ./pkg/dbConn

replace example.com/dal => ./pkg/dal

replace example.com/model => ./pkg/model

require (
	example.com/controller v0.0.0-00010101000000-000000000000
	example.com/dbConn v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.1.0
)

require (
	example.com/dal v0.0.0-00010101000000-000000000000 // indirect
	example.com/model v0.0.0-00010101000000-000000000000 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)
