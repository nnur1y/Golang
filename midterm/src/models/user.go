package models

import (
	"database/sql"
	"fmt"

	"Golang/config"
)

type User struct {
	Id       string `json:"Id"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// Create создание нового пользователя в базе
func (u User) Create() (User, error) {
	fmt.Println("start Create")
	db, _ := config.LoadDB()
	var tx *sql.Tx
	var err error
	tx, err = db.Begin()
	if err != nil {
		fmt.Println("failed to begin transaction, err:", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error calling back changes, rollbackErr:", rollbackErr)
		}
		return u, err
	}
	defer tx.Rollback()
	insertQuery := "INSERT INTO users (username, email, password) VALUES (?,?,?);"
	result, err := db.Exec(insertQuery, u.Username, u.Email, u.Password)
	if err != nil {
		fmt.Println("error preparing statement:", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
			return u, err

		}
		defer db.Close()
	}

	// Print the ID of the new user
	newUserID, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("New user ID: %d\n", newUserID)

	// var insertStmt *sql.Stmt
	// insertStmt, err = tx.Prepare("INSERT INTO users (username, email, password) VALUES (?,?,?);")
	// if err != nil {
	// 	fmt.Println("error preparing statement:", err)
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
	// 		return u, err

	// 	}
	// 	defer insertStmt.Close()
	// }

	// var result sql.Result
	// result, err = insertStmt.Exec(u.Username, u.Email, u.Password)
	// fmt.Println(result)
	// if err != nil {
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
	// 		return u, err

	// 	}
	// 	return u, err
	// }

	// if commitErr := tx.Commit(); commitErr != nil {
	// 	fmt.Println("error commiting changes:", err)
	// 	if rollbackErr := tx.Rollback(); rollbackErr != nil {
	// 		fmt.Println("there was an error rolling back changes, rollbackErr:", rollbackErr)
	// 	}
	// 	return u, err
	// }
	return u, nil
}

func (u User) Select() (User, string) {
	db, _ := config.LoadDB()
	var passInForm = u.Password
	query := "SELECT * FROM users WHERE username = ?"
	err := db.QueryRow(query, u.Username).Scan(&u.Id, &u.Username, &u.Email, &u.Password)
	if err != nil {
		fmt.Println("Username  does not found")
		return u, "Username does not found"
	}
	if u.Password != passInForm {
		return u, "Incorrect password"
	} else {
		fmt.Println("Correct, Logged in")
		return u, "logged"
	}

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
