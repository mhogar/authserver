package server

import (
	"authserver/common"
	"authserver/database"
)

// Server is an interface for starting and closing a server.
type Server interface {
	// Start starts the server and returns any errors encountered while it is running.
	Start() error

	// Close closes the server.
	Close()
}

// Runner encapsulates dependences and runs the server.
type Runner struct {
	DBConnection database.DBConnection
	Server       Server
}

// Run runs the server and returns any errors.
func (s Runner) Run() error {
	//connect to the database
	err := s.DBConnection.OpenConnection()
	if err != nil {
		return common.ChainError("error opening database connection", err)
	}

	err = s.DBConnection.Ping()
	if err != nil {
		return common.ChainError("error reaching database", err)
	}

	//start the server
	return s.Server.Start()
}
