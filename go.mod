module example.com/go-api-tech-challenge

go 1.23.0

replace example.com/controller => ./pkg/controller

replace example.com/dbconn => ./pkg/db_conn

replace example.com/dal => ./pkg/dal

replace example.com/model => ./pkg/model

require (
	example.com/controller v0.0.0-00010101000000-000000000000
	example.com/dbconn v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.1.0
)

require (
	example.com/dal v0.0.0-00010101000000-000000000000 // indirect
	example.com/model v0.0.0-00010101000000-000000000000 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
