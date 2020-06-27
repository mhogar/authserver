package sqladapter

type SQLScriptRepository interface {
	GetSQLScript(key string) string
}
