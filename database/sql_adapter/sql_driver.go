package sqladapter

// SQLDriver is an interface for encapsulating methods specific to each sql driver.
type SQLDriver interface {
	SQLScriptRepository

	// GetDriverName returns the name for the driver.
	GetDriverName() string
}
