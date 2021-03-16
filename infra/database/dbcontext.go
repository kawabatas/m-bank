package database

import (
	"context"
	"database/sql"
)

type dbContext interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}
