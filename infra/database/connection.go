package database

import (
	"fmt"
	"strings"
)

// DSN create MySQL Data Source Name.
func DSN(host, user, password, dbname string) string {
	if strings.HasPrefix(host, "/") {
		// unix socket
		return fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=true&interpolateParams=true", user, password, host, dbname)
	}
	// tcp
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&interpolateParams=true", user, password, host, dbname)
}
