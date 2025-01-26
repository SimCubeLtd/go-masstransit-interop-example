package main

import (
	"go-masstransit-interop/handlers"
	"go-masstransit-interop/internal/masstransit"
	"reflect"
)

func registerTypes() {
	masstransit.RegisterType("SomeNamespace:OrderSubmitted", reflect.TypeOf(handlers.OrderSubmitted{}))
	masstransit.RegisterType("SomeNamespace:CustomerUpdated", reflect.TypeOf(handlers.CustomerUpdated{}))
}

func registerHandlers() {
	masstransit.RegisterHandler("SomeNamespace:OrderSubmitted", handlers.HandleOrderSubmitted)
	masstransit.RegisterHandler("SomeNamespace:CustomerUpdated", handlers.HandleCustomerUpdated)
}
