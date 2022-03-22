package helpers

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	globals "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/globals"
)

func GetLoggedUser(c *gin.Context) string {
	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	userStr := fmt.Sprintf("%v", user)
	return userStr
}
