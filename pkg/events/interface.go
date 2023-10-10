package events

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() any
}

//operacoes dos eventos, quando executado
type EventHandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup)
}

//gerenciador de eventos
type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	//disparo dos handlers registrados
	Dispatch(event EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	//verificar se o evento esta registrado com o handler
	Has(eventName string, handler EventHandlerInterface) bool
	//remover todos eventos
	Clear() error
}