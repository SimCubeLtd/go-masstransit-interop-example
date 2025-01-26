package handlers

import (
	"fmt"
	"go-masstransit-interop/internal/masstransit"
	"log"
)

type OrderSubmitted struct {
	OrderID   string `json:"orderId"`
	Timestamp string `json:"timestamp"`
}

func HandleOrderSubmitted(msg interface{}, envelope masstransit.Envelope) error {
	order, ok := msg.(*OrderSubmitted)
	if !ok {
		return fmt.Errorf("invalid message type")
	}
	log.Printf("Processed OrderSubmitted: %+v\n", order)
	return nil
}
