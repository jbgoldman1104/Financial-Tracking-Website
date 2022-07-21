package main

import (
	"github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/database"
	"github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/models"
	"github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/routes"
)

func main() {
	database.OpenConnection()
	Migrate()
	routes.HandleRequests()
}

func Migrate() {
	database.DB.AutoMigrate(&models.FinancialTransaction{})
	database.DB.AutoMigrate(&models.TransactionReport{})
	database.DB.AutoMigrate(&models.User{})
	database.DB.Exec("INSERT INTO users (created_at, updated_at, nome, email, password) select now() as created_at, now() as updated_at, 'Admin' as nome, 'admin@admin.com' as email, '$2a$14$jB541M/dA8PXeHqObePUx.2JfYWFKtEoRhAQjcnvqC4DPqFcGVWUW' as password where not exists (select id from users where nome = 'Admin')")
}
