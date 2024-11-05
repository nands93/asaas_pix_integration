package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func create_webhook() (string, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var authToken string = os.Getenv("AUTHTOKEN")
	var email string = os.Getenv("EMAIL")
	var ngrok string = os.Getenv("URL")
	token_key := os.Getenv("TOKEN")

	url := "https://sandbox.asaas.com/api/v3/webhooks"

	events := []string{"PAYMENT_CREATED", "PAYMENT_RECEIVED"}
	name := "cashin"
	enabled := true
	interrupted := false
	apiVersion := 3
	sendType := "SEQUENTIALLY"

	events_string := fmt.Sprintf(`"%s"`, strings.Join(events, `","`))

	jsonPayload := fmt.Sprintf(
		`{"events":[%s],"name":"%s","url":"%s","email":"%s","enabled":%t,"interrupted":%t,"apiVersion":%d,"authToken":"%s","sendType":"%s"}`,
		events_string, name, ngrok, email, enabled, interrupted, apiVersion, authToken, sendType,
	)

	payload := strings.NewReader(jsonPayload)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("access_token", token_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(body))
	return string(body), nil
}
