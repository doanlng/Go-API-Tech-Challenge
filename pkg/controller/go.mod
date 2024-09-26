module example.com/controller

replace example.com/dal => ../dal

replace example.com/model => ../model

go 1.23.0

require (
	example.com/dal v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.1.0
)

require example.com/model v0.0.0-00010101000000-000000000000 // indirect
