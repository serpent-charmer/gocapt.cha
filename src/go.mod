module gocapt.cha/main

go 1.25.3

require (
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	gocapt.cha/captcha/dummycaptcha v0.0.0-00010101000000-000000000000
)

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	gocapt.cha/captcha/dummycaptcha/mask v0.0.0-00010101000000-000000000000 // indirect
)

replace gocapt.cha/captcha/dummycaptcha => ./captcha/dummycaptcha/

replace gocapt.cha/captcha/dummycaptcha/mask => ./captcha/dummycaptcha/mask
