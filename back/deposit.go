package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type PixKeyResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
	Type   string `json:"type"`
	Key    string `json:"key"`
}

/*func pix_key() (string, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	url := "https://sandbox.asaas.com/api/v3/pix/addressKeys"

	payload := strings.NewReader("{\"type\":\"EVP\"}")
	token_key := os.Getenv("TOKEN")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("access_token", token_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response PixKeyResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Key, nil
}*/

type QRCodeResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
	QRCode string `json:"qrCode"`
}

func create_qr_code() (string, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var value float64
	var description string
	var pix_key string = os.Getenv("PIX_KEY")

	fmt.Print("Value (R$): ")
	fmt.Scan(&value)
	fmt.Print("Add a description: ")
	fmt.Scan(&description)
	if description == "" {
		description = "Sem descrição"
	}

	url := "https://sandbox.asaas.com/api/v3/pix/qrCodes/static"

	payload := fmt.Sprintf(`{"addressKey":"%s","description":"%s","value":%.2f,"format":"PAYLOAD","allowsMultiplePayments":false}`, pix_key, description, float64(value))

	token_key := os.Getenv("TOKEN")

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("access_token", token_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var response QRCodeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Error unmarshalling response: %v", err)
	}
	fmt.Println("PIX QR Code:", response.QRCode)
	return response.QRCode, nil
}
