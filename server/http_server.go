package server

import (
	"authserver/database"
	"authserver/router"
	"fmt"
	"net/http"
)

// HTTPServer is a wrapper for an http server that implements the server interface.
type HTTPServer struct {
	http.Server
}

// CreateHTTPServerRunner creates a server runner using an http server.
func CreateHTTPServerRunner(DBConnection database.DBConnection, routerFactory router.IRouterFactory) Runner {
	server := &HTTPServer{
		Server: http.Server{
			Addr:    ":8080",
			Handler: routerFactory.CreateRouter(),
		},
	}

	return Runner{
		DBConnection: DBConnection,
		Server:       server,
	}
}

// Start starts the http server. Always returns a non-nil error.
func (s *HTTPServer) Start() error {
	fmt.Println("Server is running on port", s.Addr)
	return s.ListenAndServe()
}

// Close does nothing but exists to satisfy the server interface.
func (*HTTPServer) Close() {}
