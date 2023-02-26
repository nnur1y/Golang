package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// User ...
type User struct {
	ID        int    `json:"ID,-"`
	Login     string `json:"login"`
	Password  string `json:"Password,-"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Role      string `json:"Role,-"`
}

// Create создание нового пользователя в базе
func (u User) Create() error {
	{ // Insert a new user

		username := "FirstName"
		password := "Password"
		createdAt := time.Now()
		db, err := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
		result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
		if err != nil {
			log.Fatal(err)
		}

		id, err := result.LastInsertId()
		fmt.Println(id)
	}

	fmt.Println("Create new user with ID", u.ID)

	return nil
}
