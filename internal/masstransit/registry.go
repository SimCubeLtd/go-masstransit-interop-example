package masstransit

import (
	"reflect"
	"sync"
)

type MessageHandler func(msg interface{}, envelope Envelope) error

var (
	handlerRegistry = make(map[string]MessageHandler)
	typeRegistry    = make(map[string]reflect.Type)
	registryMutex   = sync.RWMutex{}
)

func RegisterHandler(typeName string, handler MessageHandler) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	handlerRegistry[typeName] = handler
}

func RegisterType(typeName string, msgType reflect.Type) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	typeRegistry[typeName] = msgType
}

func GetHandler(typeName string) (MessageHandler, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	handler, found := handlerRegistry[typeName]
	return handler, found
}

func GetType(typeName string) (reflect.Type, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	msgType, found := typeRegistry[typeName]
	return msgType, found
}
