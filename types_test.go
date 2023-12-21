package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	account, err := NewAccount("x", "y", "password")

	assert.Nil(t, err)

	assert.Equal(t, "x", account.FirstName)
	assert.Equal(t, "y", account.LastName)
	assert.NotEmpty(t, account.EncryptedPassword)
	assert.NotZero(t, account.Number)
	assert.NotZero(t, account.CreatedAt)
}
