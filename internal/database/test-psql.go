package database

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func NewTestPostgresConnection(t *testing.T, info ConnectionInfo) (
	*sql.DB,
	func(*sql.DB, *testing.T, string),
	func(*sql.DB, *testing.T, string),
) {
	t.Helper()

	// db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
	// 	info.Host, info.Username, info.DBName, info.SSLMode, info.Password))
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=postgres dbname=%s sslmode=%s",
		info.Host, info.DBName, info.SSLMode))
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, CreateTestTables, DeleteTestTables
}

func DeleteTestTables(db *sql.DB, t *testing.T, table string) {
	t.Helper()

	var res sql.Result
	var err error

	switch table {

	case "test_users":
		res, err = db.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", "test_users"))

	case "test_refreshsessions":
		res, err = db.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", "test_users, test_refreshsessions"))

	case "test_groups":
		res, err = db.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", "test_users, test_groups"))

	case "test_notes":
		res, err = db.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", "test_users, test_groups, test_notes"))

	default:
		t.Fatal("Invalid table name: " + table)
	}

	if err != nil {
		t.Fatal(err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		t.Fatal(err)
	}

	db.Close()
}

func CreateTestTables(db *sql.DB, t *testing.T, table string) {
	t.Helper()

	switch table {

	case "test_users":
		CreateTestTableUsers(t, db)

	case "test_refreshsessions":
		CreateTestTableUsers(t, db)
		CreateTestTableRefreshsessions(t, db)

	case "test_groups":
		CreateTestTableUsers(t, db)
		CreateTestTableGroups(t, db)

	case "test_notes":
		CreateTestTableUsers(t, db)
		CreateTestTableGroups(t, db)
		CreateTestTableNotes(t, db)

	default:
		t.Fatal("Invalid table name: " + table)
	}
}

func CreateTestTableGroups(t *testing.T, db *sql.DB) {
	t.Helper()

	res, err := db.Exec(
		`CREATE TABLE test_groups(
			id SERIAL PRIMARY KEY,
			user_login VARCHAR(100) NOT NULL REFERENCES test_users(login) ON UPDATE CASCADE ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL
		)`,
	)
	if err != nil {
		t.Fatal(err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		t.Fatal(err)
	}
}

func CreateTestTableUsers(t *testing.T, db *sql.DB) {
	t.Helper()

	res, err := db.Exec(
		`CREATE TABLE test_users(
			id SERIAL PRIMARY KEY,
			login VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(500) NOT NULL
		)`,
	)
	if err != nil {
		t.Fatal(err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		t.Fatal(err)
	}
}

func CreateTestTableRefreshsessions(t *testing.T, db *sql.DB) {
	t.Helper()

	res, err := db.Exec(
		`CREATE TABLE test_refreshsessions(
			id SERIAL PRIMARY KEY,
			user_login VARCHAR(100) NOT NULL REFERENCES test_users(login) ON UPDATE CASCADE ON DELETE CASCADE,
			fingerprint VARCHAR(300) NOT NULL,
			refreshtoken VARCHAR(300) NOT NULL,
			exp TIMESTAMP NOT NULL,
			iat TIMESTAMP NOT NULL
		)`,
	)
	if err != nil {
		t.Fatal(err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		t.Fatal(err)
	}
}

func CreateTestTableNotes(t *testing.T, db *sql.DB) {
	t.Helper()

	res, err := db.Exec(
		`CREATE TABLE test_notes(
			id SERIAL PRIMARY KEY,
			user_login VARCHAR(100) NOT NULL REFERENCES test_users(login) ON UPDATE CASCADE ON DELETE CASCADE,
			title VARCHAR(100) NOT NULL,
			text VARCHAR(10485760),
			group_id INT REFERENCES test_groups(id) ON UPDATE CASCADE ON DELETE SET NULL
		)`,
	)
	if err != nil {
		t.Fatal(err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		t.Fatal(err)
	}
}
