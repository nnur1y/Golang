package models

import (
	"Golang/config"
	"fmt"
	"log"
)

type SendCommentData struct {
	CommentText string
	UserId      string
	RecipeId    string
	Rate        string
}

func (s SendCommentData) SendComment() error {

	{ // Insert a new user

		commentText := s.CommentText
		userId := s.UserId
		recipeId := s.RecipeId
		rate := s.Rate

		db, _ := config.LoadDB()
		result, err := db.Exec(`INSERT INTO feedback( userid, recipeid, comment,rate) VALUES (?, ?, ?,?)`, userId, recipeId, commentText, rate)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)

	}
	return nil

}
