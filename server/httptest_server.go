package server

import (
	"authserver/database"
	"authserver/router"
	"net/http/httptest"
)

// HTTPTestServer is a wrapper for an httptest server that implements the server interface.
type HTTPTestServer struct {
	*httptest.Server
}

// CreateHTTPTestServerRunner creates a server runner using an httptest server.
func CreateHTTPTestServerRunner(DBConnection database.DBConnection, routerFactory router.IRouterFactory) Runner {
	server := &HTTPTestServer{
		Server: httptest.NewUnstartedServer(routerFactory.CreateRouter()),
	}

	return Runner{
		DBConnection: DBConnection,
		Server:       server,
	}
}

// Start start the server. Always returns a nil error.
func (s *HTTPTestServer) Start() error {
	s.Server.Start()
	return nil
}
