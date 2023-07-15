package main

import (
	"math/rand"
	"time"
    "crypto/sha256"
    "encoding/hex"
)

type LoginRequest struct {
    Number int64 `josn:"number"`
    Password string `json:"password"`
}

type TransferRequest struct {

    ToAccount int `json:"toAccount"` 
    Amount int `json:"amount"` 
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
    Password  string `json:"password"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
    Password string     `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName, passwd string) *Account {
    encPasswd := CreateHash(passwd) 
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
        Password: encPasswd,
		Number:    int64(rand.Intn(10000000000000000)),
		CreatedAt: time.Now().UTC(),
	}
}

func CreateHash(s string) string {
    hash := sha256.New()
    encS := hex.EncodeToString(hash.Sum([]byte(s)))
    return encS

}
