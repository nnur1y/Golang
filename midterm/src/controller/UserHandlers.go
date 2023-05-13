package controller

import (
	"Golang/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// handlerUserRegistration
// handlerUserAuthorization
// logoutHandler
func UserRegistration(c *gin.Context, store *sessions.CookieStore) *sessions.CookieStore {
	var u models.User
	u.Username = c.PostForm("Username")
	u.Email = c.PostForm("Email")
	u.Password = c.PostForm("Password")

	exists := u.UsernameExists()
	if exists {
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{
			"message": "Username already taken please try another",
		})
		return store
	}

	u, err := u.Create()
	if err != nil {

		fmt.Println("create new error: ", err)
		err = errors.New("there was an issue creating account, please try again")
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{
			"message": err,
		})
		return store
	} else {
		session, _ := store.Get(c.Request, "session")
		session.Values["user"] = u
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusFound, "/")
		return store

	}

	c.Redirect(http.StatusFound, "/")
	return store
}

func UserAuthorization(c *gin.Context, store *sessions.CookieStore) *sessions.CookieStore {
	//get users name and pass from form
	var user models.User
	user.Username = c.PostForm("Username")
	user.Password = c.PostForm("Password")
	// go to Select function in models.user to select users full info
	user, e := user.Select()

	// e = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if e == "logged" {
		//if there is no problems save it in session as value user
		session, _ := store.Get(c.Request, "session")
		session.Values["user"] = user
		session.Save(c.Request, c.Writer)
		//return to main page with users name
		c.Redirect(http.StatusFound, "/")
		// c.HTML(http.StatusOK, "layout.html", gin.H{"username": user.Username})
		return store
	} else {
		//checking any problems
		fmt.Println("error selecting hashed password in db by Username, err:", e)
		c.HTML(http.StatusUnauthorized, "authorization.html", gin.H{"message": e})
		return store
	}
	fmt.Println("err:", e)
	c.HTML(http.StatusUnauthorized, "authorization.html", gin.H{"message": "Incorrect username or password!"})
	return store
}
