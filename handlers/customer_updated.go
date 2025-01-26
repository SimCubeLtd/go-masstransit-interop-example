package handlers

import (
	"fmt"
	"go-masstransit-interop/internal/masstransit"
	"log"
)

type CustomerUpdated struct {
	CustomerID string `json:"customerId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	UpdatedAt  string `json:"updatedAt"`
}

func HandleCustomerUpdated(msg interface{}, envelope masstransit.Envelope) error {
	customer, ok := msg.(*CustomerUpdated)
	if !ok {
		return fmt.Errorf("invalid message type")
	}
	log.Printf("Processed CustomerUpdated: %+v\n", customer)
	return nil
}
