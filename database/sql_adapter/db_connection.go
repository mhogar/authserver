package sqladapter

import "errors"

func (db *SQLAdapter) OpenConnection() error {
	return errors.New("not implemented yet")
}

func (db *SQLAdapter) CloseConnection() error {
	return errors.New("not implemented yet")
}

func (db *SQLAdapter) Ping() error {
	return errors.New("not implemented yet")
}
