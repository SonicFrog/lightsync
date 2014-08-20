package main

import (
	"errors"
)

/**
 * Interface type to structs that can dispatch message
 **/
type Dispatcher interface {
	StartDispatch(chan<- Message, chan<- int)
	StopDispatch()

	RegisterHandler(name string, handler MessageHandler) error
	RemoveHandler(name string) error
}

type DefaultDispatcher struct {
	ctrl     chan int "Control channel: posting to this channel stops the dispatcher"
	handlers map[string]MessageHandler
}

func (d *DefaultDispatcher) StartDispatcher() chan<- Message {
	input := make(chan Message)
	d.ctrl = make(chan int)

	go d.dispatchRoutine(input)

	return input
}

func (d *DefaultDispatcher) StopDispatch() {
	d.ctrl <- 0
}

/**
 * This registers a handler to which all received messages will be handed
 * Multiple calls with different names will result in each message being
 * handed over multiple times!
 **/
func (d *DefaultDispatcher) RegisterHandler(name string, handler MessageHandler) error {
	if _, contains := d.handlers[name]; contains {
		LogObj.Printf("Tried to register handler %s which was already registered\n",
			name)
		return errors.New("Already registered " + name)
	}

	d.handlers[name] = handler
	return nil
}

func (d *DefaultDispatcher) RemoveHandler(name string) error {
	if _, contains := d.handlers[name]; !contains {
		LogObj.Printf("Tried to remove unregistered handler %s\n", name)
		return errors.New("Unregistered handler!")
	}

	delete(d.handlers, name)
	return nil
}

func (d *DefaultDispatcher) dispatchRoutine(input <-chan Message) {
	for {
		select {
		case msg := <-input:
			for _, h := range d.handlers {
				h.HandOver(msg)
			}

		case <-d.ctrl:
			LogObj.Println("Dispatcher exiting...")
			return
		}
	}
}
