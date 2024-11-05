package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

type PixKeyResponse struct {
	Object string `json:"object"`
	ID     string `json:"id"`
	Type   string `json:"type"`
	Key    string `json:"key"`
}

type QRCodeResponse struct {
	ID                     string `json:"id"`
	Payload                string `json:"payload"`
	AllowsMultiplePayments bool   `json:"allowsMultiplePayments"`
	ExpirationDate         string `json:"expirationDate"`
}

type WebhookNotification struct {
	Event   string `json:"event"`
	Payment struct {
		ID     string  `json:"id"`
		Value  float64 `json:"value"`
		Status string  `json:"status"`
	} `json:"payment"`
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

func create_qr_code() string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var aws_region string = os.Getenv("AWS_REGION")
	tableName := "TransacoesAsaas"

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(aws_region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	var value float64
	var pix_key string = os.Getenv("PIX_KEY")

	fmt.Print("Value (R$): ")
	fmt.Scan(&value)

	url := "https://sandbox.asaas.com/api/v3/pix/qrCodes/static"

	payload := fmt.Sprintf(`{"addressKey":"%s","value":%.2f,"format":"PAYLOAD","allowsMultiplePayments":false, "expirationDate":"2024-12-31 23:59:59"}`, pix_key, float64(value))

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

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		log.Fatalf("Erro ao desserializar JSON: %v", err)
	}
	if err != nil {
		log.Fatalf("Error unmarshalling response: %v", err)
	}

	item, err := attributevalue.MarshalMap(response)
	if err != nil {
		log.Fatalf("failed to marshal QRCodeResponse, %v", err)
	}

	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Fatalf("failed to put item, %v", err)
	}
	return response.Payload
}
