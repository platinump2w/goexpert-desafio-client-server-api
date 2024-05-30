package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		if err := ctx.Err(); errors.Is(err, context.DeadlineExceeded) {
			log.Println("[CLIENT] Execution time exceeded")
		}
		return
	default:
	}

	exchangeRate, err := getExchangeRate(ctx)
	if err != nil {
		log.Printf("[CLIENT] Failed to get exchange rate: %s\n", err.Error())
		return
	}

	err = createExchangeRateFile(exchangeRate)
	if err != nil {
		log.Printf("[CLIENT] Failed to create exchange rate file: %s\n", err.Error())
	}
}

func getExchangeRate(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func createExchangeRateFile(content []byte) error {
	_, err := os.Stat("./cotacao.txt")

	if os.IsNotExist(err) {
		_, err = os.Create("./cotacao.txt")
	}
	if err != nil {
		return err
	}

	err = os.WriteFile("./cotacao.txt",
		[]byte(fmt.Sprintf("Dólar: %s", content)),
		0644,
	)
	if err != nil {
		return err
	}

	return nil
}
