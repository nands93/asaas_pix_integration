package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

type Payment struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type PaymentEvent struct {
	ID               string  `json:"id"`
	EventType        string  `json:"event_type"`
	Payment          Payment `json:"payment"`
	PAYMENT_CREATED  bool    `json:"payment_created"`
	PAYMENT_RECEIVED bool    `json:"payment_received"`
}

/*func create_webhook() (string, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
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
}*/

func SQS_handler(sqsClient *sqs.Client, queueUrl string, payload PaymentEvent) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("fail to serialize payload: %w", err)
	}

	_, err = sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueUrl),
		MessageBody: aws.String(string(payloadBytes)),
	})
	if err != nil {
		return fmt.Errorf("message to SQS failed %w", err)
	}
	log.Printf("Message to SQS successful %s", payload.ID)
	return nil
}

func webhook_handler() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var aws_region string = os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(aws_region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	queueUrl := os.Getenv("SQS_QUEUE_URL")
	if queueUrl == "" {
		log.Fatal("SQS_QUEUE_URL not defined in .env")
	}

	http.HandleFunc("/payments-webhook", func(w http.ResponseWriter, r *http.Request) {
		var paymentEvent PaymentEvent

		err := json.NewDecoder(r.Body).Decode(&paymentEvent)
		if err != nil {
			http.Error(w, "Fail to decode JSON payload", http.StatusBadRequest)
			return
		}

		log.Printf("Received payment event: %+v", paymentEvent)

		err = SQS_handler(sqsClient, queueUrl, paymentEvent)
		if err != nil {
			http.Error(w, "Fail to send paylong to SQS Queue", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"received": true}`))
	})

	port := "8000"
	log.Printf("Running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func sendPaymentWebhook(payment PaymentEvent) error {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	payloadBytes, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("fail to serialize payload: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error on creating a request %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on sending webhook %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send webhook, status: %v", resp.Status)
	}

	log.Printf("Payment webhook sent successfully to %s", webhookURL)
	return nil
}
