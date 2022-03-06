package database

import "github.com/DATA-DOG/go-sqlmock"

func InitializePostgresMock() sqlmock.Sqlmock{
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	connection = db
	return mock
}
