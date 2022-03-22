package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	controllers "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/controllers"
	globals "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/globals"
	middleware "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/middleware"
)

var r *gin.Engine

func HandleRequests() {
	r = gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.LoadHTMLGlob("html/*")

	store := cookie.NewStore(globals.Secret)
	store.Options(sessions.Options{MaxAge: 60 * 60 * 1})
	r.Use(sessions.Sessions(globals.SessionName, store))

	handleHtml()
	private := r.Group("/")
	private.Use(middleware.AuthRequired)
	privateRoutes(private)

	r.Run()
}

func handleHtml() {
	r.GET("/", controllers.ShowLoginPage)
	r.POST("/login", controllers.ExecuteLogin)
	r.GET("/logout", controllers.ExecuteLogout)
	r.NoRoute(controllers.RouteNotFound)
}

func privateRoutes(g *gin.RouterGroup) {
	g.GET("/index", controllers.ShowIndexPage)
	g.GET("/user", controllers.ShowUserListPage)
	g.GET("/newuser", controllers.ShowNewUserPage)
	g.GET("/edituser", controllers.ShowEditUserPage)
	g.GET("/importdetail", controllers.ShowDetailImportPage)
	g.GET("/report", controllers.ShowReportPage)

	g.POST("/upload", controllers.UploadFile)
	g.POST("/insertuser", controllers.SaveNewUser)
	g.POST("/updateuser", controllers.UpdateUser)
	g.GET("/deleteuser", controllers.DeleteUser)
	g.POST("/report", controllers.GenerateReport)
}
