package models

import (
	"Golang/config"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type User struct {
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

// Create создание нового пользователя в базе
func (u User) Create() error {
	{ // Insert a new user

		email := u.Email
		password := u.Password
		firstname := u.FirstName
		lastname := u.LastName
		createdAt := time.Now()

		db, _ := config.LoadDB()
		result, err := db.Exec(`INSERT INTO users (email, password, firstname, lastname, created_at) VALUES (?, ?, ?, ?, ?)`, email, password, firstname, lastname, createdAt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)

		fmt.Println(email, password, firstname, lastname, createdAt)
	}
	return nil
}

func (u User) Select() error {
	{ // Insert a new user

		email := u.Email
		password := u.Password

		var (
			getemail    string
			getpassword string
		)
		db, err := sql.Open("mysql", "root:password@(localhost:3306)/world?parseTime=true")
		query := "SELECT email, password FROM users WHERE email = ? and password = ?"
		if err := db.QueryRow(query, email, password).Scan(&getemail, &getpassword); err != nil {
			fmt.Println("Incorrect email or password")
			return err
		}

		fmt.Println("Correct, Logged in")
		return err
	}
}
