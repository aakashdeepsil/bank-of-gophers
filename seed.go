package main

import (
	"fmt"
	"log"
)

// This function is used to seed the database with some accounts
func seedAccounts(store Storage) error {
	err := seedAccount(store, "John", "Doe", "password")
	if err != nil {
		return fmt.Errorf("error seeding account: %v", err)
	}

	return nil

}

// This function is used to seed the database with an account
func seedAccount(store Storage, firstName string, lastName string, password string) error {
	account, err := NewAccount(firstName, lastName, password)
	if err != nil {
		return fmt.Errorf("error creating account: %v", err)
	}

	if err := store.CreateAccount(account); err != nil {
		return fmt.Errorf("error creating account: %v", err)
	}

	log.Println("Seeded the database successfully!!!")

	log.Printf("account created successfully: %v", account.Number)

	return nil
}
