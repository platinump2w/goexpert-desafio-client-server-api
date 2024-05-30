package main

import (
	"context"
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	databaseName          = "exchange_rates.db"
	exchangeRateTableName = "exchange_rate"

	exchangeRateApiUrl = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	exchangeRateApiTimeout = 200 * time.Millisecond
	databaseTimeout        = 10 * time.Millisecond
)

func main() {
	log.Println("[SERVER] starting server")

	http.HandleFunc("/cotacao", handleExchangeRatesRequest)

	log.Println("[SERVER] Server started on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panicln("[SERVER] There was an error starting the server", err)
		return
	}
	log.Println("[SERVER] Server stopped")
}

func handleExchangeRatesRequest(responseWriter http.ResponseWriter, request *http.Request) {
	select {
	case <-request.Context().Done():
		log.Println("[SERVER] The request context is done")
		return
	default:
	}

	exchangeRates, err := getExchangeRates(responseWriter)
	if err != nil {
		log.Println("[SERVER] There was an error getting Exchange Rates", err)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	err = saveExchangeRate(*exchangeRates)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	responseWriter.Write([]byte(exchangeRates.UsdBrl.Bid))
}

func getExchangeRates(responseWriter http.ResponseWriter) (*ExchangeRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), exchangeRateApiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, exchangeRateApiUrl, nil)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	apiResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	defer apiResponse.Body.Close()

	body, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	var exchangeRate ExchangeRateResponse
	err = json.Unmarshal(body, &exchangeRate)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return &exchangeRate, nil
}

func setupDatabase() (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(&ExchangeRateEntity{})
	if err != nil {
		return nil, err
	}

	return database, nil
}

func saveExchangeRate(exchangeRate ExchangeRateResponse) error {
	database, err := setupDatabase()
	if err != nil {
		log.Println("[SERVER] error connecting to database", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	database.WithContext(ctx).Create(&ExchangeRateEntity{
		Code:       exchangeRate.UsdBrl.Code,
		Codein:     exchangeRate.UsdBrl.Codein,
		Name:       exchangeRate.UsdBrl.Name,
		High:       exchangeRate.UsdBrl.High,
		Low:        exchangeRate.UsdBrl.Low,
		VarBid:     exchangeRate.UsdBrl.VarBid,
		PctChange:  exchangeRate.UsdBrl.PctChange,
		Bid:        exchangeRate.UsdBrl.Bid,
		Ask:        exchangeRate.UsdBrl.Ask,
		Timestamp:  exchangeRate.UsdBrl.Timestamp,
		CreateDate: exchangeRate.UsdBrl.CreateDate,
	})

	log.Println("[SERVER] Exchange rate saved to database")

	return nil
}

type ExchangeRateResponse struct {
	UsdBrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type Tabler interface {
	TableName() string
}

func (ExchangeRateEntity) TableName() string {
	return exchangeRateTableName
}

type ExchangeRateEntity struct {
	ID         int    `gorm:"primaryKey"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}
