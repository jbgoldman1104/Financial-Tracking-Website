package models

import (
	"fmt"
	"time"

	database "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/database"
	"gorm.io/gorm"
)

type TransactionReport struct {
	gorm.Model
	DataTransacao  time.Time `json:"dataTransacao"`
	DataImportacao time.Time `json:"dataImportacao"`
	UserID         int       `json:"-"`
	UserImportacao User      `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

type TransactionReportRawQuery struct {
	ID             int
	DataTransacao  string
	DataImportacao string
	NomeUsuario    string
}

func GetAllTransactionReport() []TransactionReport {
	var transactionList []TransactionReport
	database.DB.Find(&transactionList).Order("data_importacao DESC")
	return transactionList
}

func CreateTransactionReport(transactionReport TransactionReport) {
	database.DB.Create(&transactionReport)
}

func GetAllTransactionReportRawQuery() []TransactionReportRawQuery {
	var result []TransactionReportRawQuery
	database.DB.Raw("select tr.ID as ID, to_char(tr.data_transacao, 'DD/MM/YYYY HH24:MI:SS') as data_transacao, to_char(tr.data_importacao, 'DD/MM/YYYY HH24:MI:SS') as data_importacao, u.nome as nome_usuario from transaction_reports tr inner join users u on u.id = tr.user_id order by data_importacao desc").Scan(&result)
	fmt.Println(result)
	return result
}

func GetTransactionById(id string) TransactionReportRawQuery {
	var transactionReport TransactionReportRawQuery
	database.DB.Raw("select tr.ID as ID, to_char(tr.data_transacao, 'DD/MM/YYYY HH24:MI:SS') as data_transacao, to_char(tr.data_importacao, 'DD/MM/YYYY HH24:MI:SS') as data_importacao, u.nome as nome_usuario from transaction_reports tr inner join users u on u.id = tr.user_id where tr.id = ?", id).Scan(&transactionReport)
	return transactionReport
}
