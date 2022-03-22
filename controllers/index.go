package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	models "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/models"
)

func ShowIndexPage(c *gin.Context) {
	var transactionList = models.GetAllTransactionReportRawQuery()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"transactionList": transactionList,
	})
}

func RouteNotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}
