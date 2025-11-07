module gocapt.cha/main

go 1.25.3

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	gocapt.cha/canvas v0.0.0-00010101000000-000000000000
)

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	gocapt.cha/mask v0.0.0-00010101000000-000000000000 // indirect
)

replace gocapt.cha/canvas => ./canvas
replace gocapt.cha/mask => ./mask
