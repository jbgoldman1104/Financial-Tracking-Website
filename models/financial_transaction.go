package models

import (
	"encoding/xml"
	"fmt"
	"time"

	database "github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/database"
	"github.com/marcio-olmedo-cavalini/financial-transactions-go-webapp/globals"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type FinancialTransaction struct {
	gorm.Model
	BancoOrigem       string    `json:"bancoOrigem" validate:"nonzero"`
	AgenciaOrigem     string    `json:"agenciaOrigem" validate:"nonzero"`
	ContaOrigem       string    `json:"contaOrigem" validate:"nonzero"`
	BancoDestino      string    `json:"bancoDestino" validate:"nonzero"`
	AgenciaDestino    string    `json:"agenciaDestino" validate:"nonzero"`
	ContaDestino      string    `json:"contaDestino" validate:"nonzero"`
	ValorTransacao    float64   `json:"valorTransacao" validate:"nonzero"`
	DataHoraTransacao time.Time `json:"dataHoraTransacao" validate:"nonzero"`
}

type FinancialTransactionRawQuery struct {
	DataHoraTransacao string
	BancoOrigem       string
	AgenciaOrigem     string
	ContaOrigem       string
	BancoDestino      string
	AgenciaDestino    string
	ContaDestino      string
	ValorTransacao    string
}

type AccountTransactionRawQuery struct {
	Banco             string
	Agencia           string
	Conta             string
	TipoTransacao     string
	ValorMovimentacao float64
}

type AgencyTransactionRawQuery struct {
	Banco             string
	Agencia           string
	TipoTransacao     string
	ValorMovimentacao float64
}

type Transacoes struct {
	XMLName    xml.Name    `xml:"transacoes"`
	Transacoes []Transacao `xml:"transacao"`
}

type Transacao struct {
	XMLName       xml.Name `xml:"transacao"`
	Origem        Origem   `xml:"origem"`
	Destino       Destino  `xml:"destino"`
	Valor         string   `xml:"valor"`
	DataTransacao string   `xml:"data"`
}

type Origem struct {
	XMLName xml.Name `xml:"origem"`
	Banco   string   `xml:"banco"`
	Agencia string   `xml:"agencia"`
	Conta   string   `xml:"conta"`
}

type Destino struct {
	XMLName xml.Name `xml:"destino"`
	Banco   string   `xml:"banco"`
	Agencia string   `xml:"agencia"`
	Conta   string   `xml:"conta"`
}

func ValidateFinancialTransaction(financialTransaction *FinancialTransaction) error {
	if err := validator.Validate(financialTransaction); err != nil {
		return err
	}

	return nil
}

func CreateFinancialTransaction(financialTransaction FinancialTransaction) {
	database.DB.Create(&financialTransaction)
}

func ExistsFinancialTransactionByDate(dataFiltro time.Time) bool {
	var ft FinancialTransaction
	result := database.DB.Where("data_hora_transacao >= ?", dataFiltro.Format("2006-01-02")).First(&ft)
	return result.RowsAffected != 0
}

func GetAllFinancialTransactionRawQuery(id string) []FinancialTransactionRawQuery {
	var result []FinancialTransactionRawQuery
	database.DB.Raw("select to_char(ft.data_hora_transacao, 'DD/MM/YYYY') as data_hora_transacao, ft.banco_origem, ft.agencia_origem, ft.conta_origem, ft.banco_destino, ft.agencia_destino, ft.conta_destino, ft.valor_transacao::text as valor_transacao from financial_transactions ft WHERE to_char(ft.data_hora_transacao, 'DD/MM/YYYY') = (select to_char(data_transacao, 'DD/MM/YYYY') from transaction_reports where id = ?)", id).Scan(&result)
	return result
}

func GetSuspectedFinancialTransactionRawQuery(month string) []FinancialTransactionRawQuery {
	var result []FinancialTransactionRawQuery
	database.DB.Raw("select to_char(ft.data_hora_transacao, 'DD/MM/YYYY') as data_hora_transacao, ft.banco_origem, ft.agencia_origem, ft.conta_origem, ft.banco_destino, ft.agencia_destino, ft.conta_destino, ft.valor_transacao::text as valor_transacao from financial_transactions ft WHERE ft.valor_transacao >= ? and to_char(ft.data_hora_transacao, 'YYYYMM') = ?", globals.SuspectedTransactionValue, month).Scan(&result)
	return result
}

func GetSuspectedAccountTransactionRawQuery(month string) []AccountTransactionRawQuery {
	var result []AccountTransactionRawQuery
	database.DB.Raw("select ft.banco_origem as banco, ft.agencia_origem as agencia, ft.conta_origem as conta, 'SAIDA' as tipo_transacao, sum(ft.valor_transacao) as valor_transacao from financial_transactions ft WHERE ft.valor_transacao >= ? and to_char(ft.data_hora_transacao, 'YYYYMM') = ? group by ft.banco_origem, ft.agencia_origem, ft.conta_origem union select ft.banco_destino as banco, ft.agencia_destino as agencia, ft.conta_destino as conta, 'ENTRADA' as tipo_transacao, sum(ft.valor_transacao) as valor_transacao from financial_transactions ft WHERE ft.valor_transacao >= ? and to_char(ft.data_hora_transacao, 'YYYYMM') = ? group by ft.banco_destino, ft.agencia_destino, ft.conta_destino", globals.SuspectedBankAccountValue, month, globals.SuspectedBankAccountValue, month).Scan(&result)
	fmt.Println(result)
	return result
}

func GetSuspectedAgencyTransactionRawQuery(month string) []AgencyTransactionRawQuery {
	var result []AgencyTransactionRawQuery
	database.DB.Raw("select ft.banco_origem as banco, ft.agencia_origem as agencia, 'SAIDA' as tipo_transacao, sum(ft.valor_transacao) as valor_transacao from financial_transactions ft WHERE ft.valor_transacao >= ? and to_char(ft.data_hora_transacao, 'YYYYMM') = ? group by ft.banco_origem, ft.agencia_origem union select ft.banco_destino as banco, ft.agencia_destino as agencia, 'ENTRADA' as tipo_transacao, sum(ft.valor_transacao) as valor_transacao from financial_transactions ft WHERE ft.valor_transacao >= ? and to_char(ft.data_hora_transacao, 'YYYYMM') = ? group by ft.banco_destino, ft.agencia_destino", globals.SuspectedBankAgencyValue, month, globals.SuspectedBankAgencyValue, month).Scan(&result)
	fmt.Println(result)
	return result
}
