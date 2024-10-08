module example.com/controller

replace example.com/dal => ../dal

replace example.com/model => ../model

replace example.com/dbconn => ../db_conn

go 1.23.0

require (
	example.com/dal v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.1.0
)

require (
	example.com/model v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
