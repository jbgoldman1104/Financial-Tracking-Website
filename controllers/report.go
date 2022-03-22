package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/models"
)

func ShowReportPage(c *gin.Context) {
	c.HTML(http.StatusOK, "report.html", nil)
}

func GenerateReport(c *gin.Context) {
	month := c.PostForm("month")
	fmt.Println("month" + month)
	transactionReport := models.GetSuspectedFinancialTransactionRawQuery(month)
	accountReport := models.GetSuspectedAccountTransactionRawQuery(month)
	agencyReport := models.GetSuspectedAgencyTransactionRawQuery(month)
	c.HTML(http.StatusOK, "report.html", gin.H{
		"content":           0,
		"transacaoSuspeita": transactionReport,
		"contaSuspeita":     accountReport,
		"agenciaSuspeita":   agencyReport,
	})
}
