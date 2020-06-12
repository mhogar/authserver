package migrations

import (
	sqladapter "authserver/database/sql_adapter"
)

type m20200611224101 struct {
	Adapter *sqladapter.SQLAdapter
}

func (m m20200611224101) GetTimestamp() string {
	return "20200611224101"
}

func (m m20200611224101) Up() error {
	return nil
}

func (m m20200611224101) Down() error {
	return nil
}
