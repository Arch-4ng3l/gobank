package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=moritz dbname=gobank password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil

}

func (psql *PostgresStore) Init() error {
	return psql.createAccountTable()
}

func (psql *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id serial PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number serial, 
		balance INT, 
		created_at TIMESTAMP
	)`

	psql.db.Query(query)

	return nil
}

func (psql *PostgresStore) CreateAccount(acc *Account) error {

	query := `INSERT INTO accounts 
						(first_name, last_name, balance, created_at)
						VALUES($1, $2, $3, $4)`

	res, err := psql.db.Query(query,
		acc.FirstName,
		acc.LastName,
		acc.Balance,
		acc.CreatedAt)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", res)

	return nil
}

func (psql *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (psql *PostgresStore) UpdateAccount(acc *Account) error {
	return nil
}

func (psql *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}

func (psql *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := psql.db.Query("SELECT * FROM accounts")

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		acc := &Account{}

		err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance, &acc.CreatedAt)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)

	}

	return accounts, nil
}
