package events

import (
	"errors"
	"sync"
)

var errHandlerAlreadyRegistered = errors.New("handler already registered in this event")

type EventDispatcher struct {
	//um evento pode ter varios handlers, executar diversas tarefas
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher{
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers,ok := ed.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if handlers,ok := ed.handlers[eventName]; ok {
		for i,h := range handlers {
			if h == handler{
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func(ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	//se existir evento registrado com o nome de evento inserido
	if _,ok := ed.handlers[eventName]; ok {
		//percorrer todos handlers do evento
		for _,h := range ed.handlers[eventName] {
			//verificar se o handler inserido ja existe 
			if h == handler {
				return errHandlerAlreadyRegistered
			}
		}
	}
	//adicionar evento no map
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}



func(ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	//verificar se o evento esta registrado
	if _,ok := ed.handlers[eventName]; ok {
		//loop nos handlers do evento
		for _,h := range ed.handlers[eventName] {
			if h == handler{
				return true
			}
		}
	}
	return false
}

func(ed *EventDispatcher) Clear() error {
	ed.handlers = make(map[string][]EventHandlerInterface)
	return nil
}