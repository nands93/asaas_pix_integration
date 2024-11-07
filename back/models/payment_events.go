package models

type Payment struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type PaymentEvent struct {
	ID               string `json:"id"`
	EventType        string `json:"event_type"`
	PAYMENT_CREATED  bool   `json:"payment_created"`
	PAYMENT_RECEIVED bool   `json:"payment_received"`
}
