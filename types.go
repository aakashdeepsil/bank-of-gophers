package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type CreateTransferRequest struct {
	ToAccountNumber int64 `json:"to_account_number"`
	Amount          int64 `json:"amount"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Number            int64     `json:"number"`
	Balance           int64     `json:"balance"`
	EncryptedPassword string    `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewAccount(firstName string, lastName string, password string) (*Account, error) {
	encrypted_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	log.Printf("encrypted_password generated: %v", encrypted_password)
	if err != nil {
		return nil, fmt.Errorf("error generating password hash: %v", err)
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(rand.Intn(1000000)),
		EncryptedPassword: string(encrypted_password),
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func (a *Account) ValidatePassword(password string) error {
	if a.EncryptedPassword == "" {
		return fmt.Errorf("account password cannot be empty")
	}

	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password))
}

type APIServer struct {
	listenAddress string
	store         Storage
}

type APIError struct {
	Message string `json:"message"`
}

type APIHandleFunc func(w http.ResponseWriter, r *http.Request) error

type Storage interface {
	// CreateAccount creates a new account in the storage.
	CreateAccount(account *Account) error

	// GetAccount returns an account by its ID.
	GetAccountByID(id int) (*Account, error)

	// GetAccounts returns all accounts.
	GetAccount() ([]*Account, error)

	// DeleteAccount deletes all accounts from the storage.
	DeleteAccount() error

	// DeleteAccountByID deletes an account by its ID.
	DeleteAccountByID(id int) error

	// TransferMoney transfers money from one account to another.
	TransferMoney(fromID int64, toID int64, amount int64) error

	// GetAccountByNumber returns an account by its number.
	GetAccountByNumber(number int64) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Number      int64  `json:"number"`
}
