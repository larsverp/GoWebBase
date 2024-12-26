package application

import (
	"database/sql"
)

type Application struct {
	DB  *sql.DB
	Log Logger
}
