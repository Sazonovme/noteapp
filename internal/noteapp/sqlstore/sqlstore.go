package sqlstore

import (
	"database/sql"
	"noteapp/pkg/logger"
)

type Sqlstore struct {
	Db             *sql.DB
	UserRepository *UserRepository
	AuthRepository *AuthRepository
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
	store := &Sqlstore{Db: db}
	store.ConfigureRepositories()
	return store, nil
}

func (s *Sqlstore) ConfigureRepositories() {
	s.AuthRepository = NewAuthRepository(s)
	s.UserRepository = NewUserRepository(s)
}
