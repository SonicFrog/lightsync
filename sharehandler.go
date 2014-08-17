package main

import (
	"errors"
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
	for sh.Active() {
		select {
		case evt := <- sh.Events() :
			//Event
			fmt.Println("event: ", evt)
			break

		case err := <- sh.Errors() :
			//Error
			fmt.Println("error: ", err)
			break
		}
	}
}

func (sh *ShareHandler) Handle(msg Message) (err error) {

	switch msg.(type) {
	case FileRemoveMessage:
	case DirectoryRemoveMessage:
		sh.Remove(msg.Name())
		break

	case FileUpdatedMessage:
		//Send FilePartRequest for each part of the file
		break

	case FileCreatedMessage:
		sh.CreateFile(msg.Name())
		break

	case FileHashMessage:
		sh.CheckHash(msg.Name(), msg.Payload())
		break

	case ShareACKMessage:

		break

	case ShareLeaveMessage:
		break

	case DirectoryCreateMessage:
		break

	default:
		err = errors.New("Invalid message type handed to ShareHandler!")
	}

	return
}

func (sh *ShareHandler) CheckHash(name string, hash []byte) bool {
	return true
}

func (sh *ShareHandler) Stop() {
	sh.isActive = false
}
