package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	globals "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/globals"
	models "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/models"
)

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func ExecuteLogin(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Por favor, faça logout primeiro"})
		return
	}

	username := c.PostForm("email")
	password := c.PostForm("password")

	if emptyUserPass(username, password) {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "As informações email/password não podem ser nulas"})
		return
	}

	if !checkUserPass(username, password) {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Usuário/Senha inválidos"})
		return
	}

	session.Set(globals.Userkey, username)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Falha ao salvar a sessão"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/index")
}

func ExecuteLogout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token de sessão inválido"})
		return
	}
	session.Delete(globals.Userkey)
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar sessão"})
		return
	}

	c.HTML(http.StatusInternalServerError, "login.html", nil)
}

func emptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func checkUserPass(username, password string) bool {
	userInfo := models.FindUserByEmail(username)
	return models.CheckPasswordHash(password, userInfo.Password)
}
