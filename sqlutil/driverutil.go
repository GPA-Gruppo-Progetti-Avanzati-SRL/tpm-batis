package sqlutil

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
)

var sqlDriverNamesByType map[reflect.Type]string

// SQLDriverToDriverName The database/sql API doesn't provide a way to get the registry name for
// a driver from the driver type.
func SQLDriverToDriverName(driver driver.Driver) string {

	if sqlDriverNamesByType == nil {
		sqlDriverNamesByType = map[reflect.Type]string{}

		for _, driverName := range sql.Drivers() {
			// Tested empty string DSN with MySQL, PostgreSQL, and SQLite3 drivers.
			db, _ := sql.Open(driverName, "")

			if db != nil {
				driverType := reflect.TypeOf(db.Driver())
				sqlDriverNamesByType[driverType] = driverName
			}
		}
	}

	driverType := reflect.TypeOf(driver)
	if driverName, found := sqlDriverNamesByType[driverType]; found {
		return driverName
	}

	return ""
}
