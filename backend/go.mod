module appcraft

go 1.20

require (
	github.com/creasty/defaults v1.6.0
	httpclient v0.0.0
	httpserver v0.0.0
)

require github.com/gorilla/mux v1.8.0 // indirect

replace httpserver => ./httpserver

replace httpclient => ./httpclient
