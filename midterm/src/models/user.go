package models

import (
	"database/sql"
	"fmt"

	"Golang/config"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `json:"Id"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// Create создание нового пользователя в базе
func (u *User) Create() error {
	db, _ := config.LoadDB()
	var hash []byte
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		return err
	}

	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		fmt.Println("failed to begin transaction, err:", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error calling back changes, rollbackErr:", rollbackErr)
		}
		return err
	}
	defer tx.Rollback()
	var insertStmt *sql.Stmt
	insertStmt, err = tx.Prepare("INSERT INTO users (username, email, password) VALUES (?,?,?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
			return err

		}
		defer insertStmt.Close()
	}

	var result sql.Result
	result, err = insertStmt.Exec(u.Username, u.Email, hash)
	aff, _ := result.RowsAffected()
	if aff == 0 {
		return err
	}
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
			return err

		}
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		fmt.Println("error commiting changes:", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
		}
		return err
	}
	return nil
}

func (u *User) Select() error {
	db, _ := config.LoadDB()
	query := "SELECT * FROM users WHERE username = ?"
	err := db.QueryRow(query, u.Username).Scan(&u.Id, &u.Username, &u.Email, &u.Password)
	if err != nil {
		fmt.Println("getUser() error selecting User, err:", err)
		return err
	}

	fmt.Println("Correct, Logged in")
	return nil
}

func (u *User) UsernameExists() (exists bool) {
	exists = true
	db, _ := config.LoadDB()
	stmt := "SELECT Id FROM users WHERE Username = ?"
	row := db.QueryRow(stmt, u.Username)
	var uId string
	err := row.Scan(&uId)
	if err == sql.ErrNoRows {
		return false
	}
	return exists
}
