package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// This function is used to create a new Postgres storage
func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=bank-of-gophers sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

// This function is used to initialize the Postgres table
func (s *PostgresStorage) Init() error {
	return s.CreateAccountTable()
}

// This function is used to create the accounts table
func (s *PostgresStorage) CreateAccountTable() error {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS accounts (id serial primary key, first_name varchar(255), last_name varchar(255), encrypted_password varchar(255), number serial, balance bigint, created_at timestamp default current_timestamp)")
	return err
}

// This function is used to insert a new account into the database
func (s *PostgresStorage) CreateAccount(account *Account) error {
	response, err := s.db.Exec("INSERT INTO accounts (first_name, last_name, encrypted_password, number, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		account.FirstName, account.LastName, account.EncryptedPassword, account.Number, account.Balance, account.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Account created:", response)

	return err
}

// This function is used to get all accounts from the database
func (s *PostgresStorage) GetAccount() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}

	var accounts []*Account
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.EncryptedPassword, &account.Number, &account.Balance, &account.CreatedAt); err != nil {
			return nil, err
		}

		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// This function is used to get an account by its ID
func (s *PostgresStorage) GetAccountByID(id int) (*Account, error) {
	var account Account
	err := s.db.QueryRow("SELECT id, first_name, last_name, number, balance,  created_at FROM accounts WHERE id = $1", id).
		Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// This function is used to delete all accounts from the database
func (s *PostgresStorage) DeleteAccount() error {
	_, err := s.db.Exec("DELETE FROM accounts")
	return err
}

// This function is used to delete an account by its ID
func (s *PostgresStorage) DeleteAccountByID(id int) error {
	_, err := s.db.Exec("DELETE FROM accounts WHERE id = $1", id)
	return err
}

// This function is used to transfer money from one account to another
func (s *PostgresStorage) TransferMoney(fromAccountNumber int64, toAccountNumber int64, amount int64) error {
	return nil
}

// This function is used to get an account by its number and password
func (s *PostgresStorage) GetAccountByNumber(number int64) (*Account, error) {
	var account Account
	err := s.db.QueryRow("SELECT id, first_name, last_name, encrypted_password, number, balance, created_at FROM accounts WHERE number = $1", number).
		Scan(&account.ID, &account.FirstName, &account.LastName, &account.EncryptedPassword, &account.Number, &account.Balance, &account.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %v", err)
	}

	return &account, nil
}
