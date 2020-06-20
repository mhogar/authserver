package migrations

import (
	postgresadapter "authserver/database/postgres_adapter"
)

type m20200611224101 struct {
	Adapter *postgresadapter.PostgresAdapter
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
