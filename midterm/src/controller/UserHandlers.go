package controller

import (
	"Golang/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// handlerUserRegistration
// handlerUserAuthorization
// logoutHandler
func UserRegistration(c *gin.Context) {
	var u models.User
	u.Username = c.PostForm("Username")
	u.Email = c.PostForm("Email")
	u.Password = c.PostForm("Password")

	exists := u.UsernameExists()
	if exists {
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{
			"message": "Username already taken please try another",
		})
		return
	}

	err := u.Create()
	if err != nil {
		fmt.Println("create new error: ", err)
		err = errors.New("there was an issue creating account, please try again")
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{
			"message":  err,
			"username": u.Username,
		})
		return
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"username": u.Username,
		//"content":  RecipesList,
	})
}

func UserAuthorization(c *gin.Context, store *sessions.CookieStore) {
	var user models.User
	user.Username = c.PostForm("Username")
	password := c.PostForm("Password")
	e := user.Select()
	if e != nil {
		fmt.Println("error selecting hashed password in db by Username, err:", e)
		c.HTML(http.StatusUnauthorized, "authorization.html", gin.H{"message": "Incorrect username and password!"})
		return
	}
	e = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if e == nil {
		session, _ := store.Get(c.Request, "session")
		session.Values["user"] = user
		session.Save(c.Request, c.Writer)
		c.HTML(http.StatusOK, "layout.html", gin.H{"username": user.Username})
		return
	}
	fmt.Println("err:", e)
	c.HTML(http.StatusUnauthorized, "authorization.html", gin.H{"message": "Incorrect username or password!"})
}
