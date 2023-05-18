package controller

import (
	"Golang/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func Logout(c *gin.Context, store *sessions.CookieStore) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Print(err)
	}

	session.Options.MaxAge = -1

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		fmt.Println("Failed to drop session:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{"message": "Logged out"})
}

func Profile(c *gin.Context, store *sessions.CookieStore) {
	var user models.User
	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("session:", session)
	userVal, ok := session.Values["user"]

	if !ok {
		fmt.Println("no user", userVal)
		user.Username = ""
	} else {
		user = *userVal.(*models.User)
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{"user": user})

}
