package sqlstore

import (
	"database/sql"
	"noteapp/pkg/logger"
)

type Sqlstore struct {
	Db             *sql.DB
	UserRepository *UserRepository
}

// Constructor
func NewStore(db_connection string) (*Sqlstore, error) {
	db, err := sql.Open("postgres", db_connection)
	if err != nil {
		logger.NewLog("sqlstore", "NewStore", err, nil, 2, "failed to connect to database")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Sqlstore{Db: db}, nil
}
