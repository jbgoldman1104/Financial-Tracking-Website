package main

import (
	"github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/database"
	"github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/models"
	"github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/routes"
)

func main() {
	database.OpenConnection()
