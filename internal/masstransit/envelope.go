package masstransit

import (
	"encoding/json"
	"fmt"
	"go-masstransit-interop/internal/worker"
	"reflect"
	"strings"
)

type Envelope struct {
	MessageId     string                 `json:"messageId"`
	CorrelationId *string                `json:"correlationId"`
	MessageType   []string               `json:"messageType"`
	Message       json.RawMessage        `json:"message"`
	Headers       map[string]interface{} `json:"headers"`
	Host          HostInfo               `json:"host"`
}

type HostInfo struct {
	MachineName            string `json:"machineName"`
	ProcessName            string `json:"processName"`
	ProcessId              int    `json:"processId"`
	Assembly               string `json:"assembly"`
	AssemblyVersion        string `json:"assemblyVersion"`
	FrameworkVersion       string `json:"frameworkVersion"`
	MassTransitVersion     string `json:"massTransitVersion"`
	OperatingSystemVersion string `json:"operatingSystemVersion"`
}

func ExtractTypeName(urn string) (string, error) {
	parts := strings.Split(urn, ":")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid URN format")
	}
	return parts[2] + ":" + parts[3], nil
}

func ProcessMessage(envelope Envelope, workerPool *worker.Pool) error {
	if len(envelope.MessageType) == 0 {
		return fmt.Errorf("missing messageType")
	}

	typeName, err := ExtractTypeName(envelope.MessageType[0])
	if err != nil {
		return fmt.Errorf("invalid messageType: %w", err)
	}

	messageHandler, found := GetHandler(typeName)
	if !found {
		return fmt.Errorf("handler not found for messageType: %s", typeName)
	}

	msgType, found := GetType(typeName)
	if !found {
		return fmt.Errorf("unknown message type: %s", typeName)
	}

	msgInstance := reflect.New(msgType).Interface()

	if err := json.Unmarshal(envelope.Message, msgInstance); err != nil {
		return fmt.Errorf("error unmarshaling message: %w", err)
	}

	workerPool.Submit(func() {
		if err := messageHandler(msgInstance, envelope); err != nil {
			fmt.Printf("Error processing message: %v\n", err)
		}
	})

	return nil
}
