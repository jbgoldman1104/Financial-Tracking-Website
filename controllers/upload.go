package controllers

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/helpers"
	models "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/models"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("myFile")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	filename := filepath.Base(file.Filename)
	fmt.Println(filename)

	if err := c.SaveUploadedFile(file, "upload/"+filename); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	fmt.Println(strings.Contains("GeeksforGeeks", "for"))

	validationMessage := ""
	var transactionList []models.TransactionReportRawQuery
	var totalRows int

	if strings.Contains(strings.ToLower(filename), "csv") {
		totalRows = readAndPrintUploadedFileCsv(filename)

		validationMessage = readAndLoadUploadedFileCsv(filename, helpers.GetLoggedUser(c))
		transactionList = models.GetAllTransactionReportRawQuery()
	} else {
		totalRows, validationMessage = LoadXml(filename, helpers.GetLoggedUser(c))
	}

	if totalRows == 0 {
		validationMessage = "O arquivo informado está vazio"
	}

	c.HTML(http.StatusOK, "upload.html", gin.H{
		"mensagem":         validationMessage,
		"arquivo":          file.Filename,
		"tamaho":           strconv.Itoa(int(file.Size)),
		"quantidadeLinhas": totalRows,
		"transactionList":  transactionList,
	})
}

func openUploadedFile(filename string) (http.File, [][]string) {
	d := http.Dir("./upload")
	f, err := d.Open(filename)
	if err != nil {
		panic(err)
	}

	filedata, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	return f, filedata
}

func readAndPrintUploadedFileCsv(filename string) int {
	f, filedata := openUploadedFile(filename)
	totalRows := len(filedata)

	for e, value := range filedata {
		fmt.Println(e, value)
	}

	defer f.Close()
	io.Copy(os.Stdout, f)
	return totalRows
}

func readAndLoadUploadedFileCsv(filename string, emailLoggedUser string) string {
	f, filedata := openUploadedFile(filename)
	var dataLote time.Time
	var dataTransacao time.Time

	for e, value := range filedata {
		if e == 0 {
			dataLoteTmp, _ := time.Parse("2006-01-02T15:04:05", value[7])
			dataLote = time.Date(dataLoteTmp.Year(), dataLoteTmp.Month(), dataLoteTmp.Day(), 0, 0, 0, 0, dataLoteTmp.Location())

			if validateDateTransaction(dataLote) {
				fmt.Println("O Lote com essa data já está cadastrado!")
				return "O Lote com essa data já está cadastrado!"
			}
		}

		dataTransacao, _ = time.Parse("2006-01-02T15:04:05", value[7])
		dataTransacao = time.Date(dataTransacao.Year(), dataTransacao.Month(), dataTransacao.Day(), 0, 0, 0, 0, dataTransacao.Location())

		if dataTransacao == dataLote {
			var financialTransaction = new(models.FinancialTransaction)
			financialTransaction.BancoOrigem = value[0]
			financialTransaction.AgenciaOrigem = value[1]
			financialTransaction.ContaOrigem = value[2]
			financialTransaction.BancoDestino = value[3]
			financialTransaction.AgenciaDestino = value[4]
			financialTransaction.ContaDestino = value[5]
			financialTransaction.ValorTransacao, _ = strconv.ParseFloat(value[6], 64)

			dateString := value[7]
			dateConverted, _ := time.Parse("2006-01-02T15:04:05", dateString)
			financialTransaction.DataHoraTransacao = dateConverted
			//fmt.Println(e, financialTransaction.DataHoraTransacao)
			if err := models.ValidateFinancialTransaction(financialTransaction); err == nil {
				models.CreateFinancialTransaction(*financialTransaction)
			}
		}

		if e == len(filedata)-1 {
			user := models.FindUserByEmail(emailLoggedUser)
			transactionReport := models.TransactionReport{DataTransacao: dataLote, DataImportacao: time.Now(), UserImportacao: user}
			models.CreateTransactionReport(transactionReport)
		}
	}

	defer f.Close()
	return ""
}

func validateDateTransaction(dataLote time.Time) bool {
	return models.ExistsFinancialTransactionByDate(dataLote)
}

func ShowDetailImportPage(c *gin.Context) {
	id := c.Query("id")
	fmt.Println("id" + id)
	transactionReport := models.GetTransactionById(id)
	financialTransaction := models.GetAllFinancialTransactionRawQuery(id)
	c.HTML(http.StatusOK, "detailimport.html", gin.H{
		"DataImportacao": transactionReport.DataImportacao,
		"NomeUsuario":    transactionReport.NomeUsuario,
		"DataTransacao":  transactionReport.DataTransacao,
		"transactions":   financialTransaction,
	})
}

func LoadXml(filename string, emailLoggedUser string) (int, string) {
	xmlFile, err := os.Open("./upload/" + filename)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened " + filename)
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var transacoes models.Transacoes
	xml.Unmarshal(byteValue, &transacoes)

	for i := 0; i < len(transacoes.Transacoes); i++ {
		fmt.Println("Origem: " + transacoes.Transacoes[i].Origem.Banco)
		fmt.Println("Destino: " + transacoes.Transacoes[i].Destino.Banco)
		fmt.Println("Valor: " + transacoes.Transacoes[i].Valor)
		fmt.Println("DataTransacao: " + transacoes.Transacoes[i].DataTransacao)
	}

	var dataLote time.Time
	var dataTransacao time.Time

	for e, value := range transacoes.Transacoes {
		if e == 0 {
			dataLoteTmp, _ := time.Parse("2006-01-02T15:04:05", value.DataTransacao)
			dataLote = time.Date(dataLoteTmp.Year(), dataLoteTmp.Month(), dataLoteTmp.Day(), 0, 0, 0, 0, dataLoteTmp.Location())

			if validateDateTransaction(dataLote) {
				fmt.Println("O Lote com essa data já está cadastrado!")
				return 0, "O Lote com essa data já está cadastrado!"
			}
		}

		dataTransacao, _ = time.Parse("2006-01-02T15:04:05", value.DataTransacao)
		dataTransacao = time.Date(dataTransacao.Year(), dataTransacao.Month(), dataTransacao.Day(), 0, 0, 0, 0, dataTransacao.Location())

		if dataTransacao == dataLote {
			var financialTransaction = new(models.FinancialTransaction)
			financialTransaction.BancoOrigem = value.Origem.Banco
			financialTransaction.AgenciaOrigem = value.Origem.Agencia
			financialTransaction.ContaOrigem = value.Origem.Conta
			financialTransaction.BancoDestino = value.Destino.Banco
			financialTransaction.AgenciaDestino = value.Destino.Agencia
			financialTransaction.ContaDestino = value.Destino.Conta
			financialTransaction.ValorTransacao, _ = strconv.ParseFloat(value.Valor, 64)

			dateString := value.DataTransacao
			dateConverted, _ := time.Parse("2006-01-02T15:04:05", dateString)
			financialTransaction.DataHoraTransacao = dateConverted
			if err := models.ValidateFinancialTransaction(financialTransaction); err == nil {
				models.CreateFinancialTransaction(*financialTransaction)
			}
		}

		if e == len(transacoes.Transacoes)-1 {
			user := models.FindUserByEmail(emailLoggedUser)
			transactionReport := models.TransactionReport{DataTransacao: dataLote, DataImportacao: time.Now(), UserImportacao: user}
			models.CreateTransactionReport(transactionReport)
		}
	}

	quantidadeLinhas := len(transacoes.Transacoes)
	return quantidadeLinhas, ""
}
