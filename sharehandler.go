package main

import (
	"fmt"
)

type ShareHandler struct {
	Share
	RequestChannel chan Message
	isActive bool
}

type MessageHandler interface {
	Handle(Message) error
}

func NewShareHandler(share Share, out chan Message) (sh *ShareHandler) {
	sh = &ShareHandler{share, out, true}
	go sh.HandleLocal()
	return
}

func (sh *ShareHandler) Active() bool {
	return sh.isActive
}

func (sh *ShareHandler) HandleLocal() {
	defer sh.Share.Watcher.Close()
	for {
		select {
		case evt := <- sh.Events() :
			//Event
			fmt.Println("event: ", evt)
			break

		case err := <- sh.Errors() :
			//Error
			fmt.Println("error: ", err)


		case <- Control :
			LogObj.Println("ShareHandler ", sh.Name, " stopping!")
			break
		}

	}
}

func (sh *ShareHandler) Handle(msg Message) (err error) {
	switch msg.(type) {
	case FileRemoveMessage:
		msg, ok := msg.(FileRemoveMessage)

		if !ok {
			panic("Invalid message type while handling message in share")
		}

		sh.Remove(msg.Name())
		break

	case FileCreatedMessage:
		msg, ok := msg.(FileCreatedMessage)

		if !ok {
			panic("Invalid message type while handling message in share")
		}

		sh.CreateFile(msg.FileName())
		break

	case FileHashMessage:
		msg, ok := msg.(FileHashMessage)
		if !ok {
			panic("Invalid message type while handling message!")
		}

		if msg.Share() != sh.Name {

		}
		sh.CheckHash(msg.FileName(), msg.Hash())


	case ShareACKMessage:
		msg, ok := msg.(ShareACKMessage)
		if !ok {
			panic("Fatal while handling message in share " + sh.Name)
		}
		if msg.Share() != sh.Name {
			panic("Share "+sh.Name+" had to handle a message for "+msg.Share())
		}
		sh.AddClient(msg.Sender())
		break

	case ShareLeaveMessage:
		msg, ok := msg.(ShareLeaveMessage)
		if !ok {
			panic("Wow!")
		}

		if msg.Share() != sh.Name {
			panic("Share "+sh.Name+ " had to handle a message for "+msg.Share())
		}
		sh.RemoveClient(msg.Sender())

	case DirectoryCreateMessage:
		msg, ok := msg.(DirectoryCreateMessage)
		if !ok {
			panic("Invalid message type while handling message!")
		}
		if msg.Share() != sh.Name {
			panic("ShareHandler received a message for another share!")
		}
		sh.CreateDir(msg.Name())


	default:
		panic("Invalid message received in ShareHandler!!")
	}

	return
}

func (sh *ShareHandler) CheckHash(name string, hash []byte) bool {
	return true
}

func (sh *ShareHandler) Stop() {
	sh.isActive = false
}
