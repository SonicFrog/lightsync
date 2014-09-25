package main

import (
	"lightsync/proto"
)

type ShareHandler struct {
	Share
	requestChannel chan Message
	controlChannel chan int
}

func NewShareHandler(share Share, out chan Message) (sh *ShareHandler) {
	sh = &ShareHandler{
		Share:          share,
		requestChannel: out,
		controlChannel: make(chan int),
	}

	go sh.handleLocal()

	return
}

func (sh *ShareHandler) HandOver(msg Message) {
	sh.requestChannel <- msg
}

func (sh *ShareHandler) handleLocal() {
	defer sh.Share.Watcher.Close()

	for {
		select {
		case evt := <-sh.Events():
			//Event
			LogObj.Println("event: ", evt)
			break

		case err := <-sh.Errors():
			//Error
			LogObj.Println("error: ", err)

		case <-sh.controlChannel:
			LogObj.Println("ShareHandler ", sh.Name, " stopping!")
			return

		case msg := <-sh.requestChannel:
			LogObj.Println("Handling message")
			sh.Handle(msg)
		}
	}
}

func (sh *ShareHandler) HandleFile(msg *FileMessageWrapper) {
	if msg.GetShareName() != sh.Name {
		LogObj.Println("Ignoring message meant for share", msg.GetShareName())
		return
	}

	switch msg.GetAction() {
	case light.FileAction_REMOVED:
		sh.Remove(msg.GetFilename())

	case light.FileAction_CREATED:
		if msg.GetFolder() {
			sh.CreateDir(msg.GetFilename())
		} else {
			sh.CreateFile(msg.GetFilename())
		}

	case light.FileAction_UPDATED:
		sh.CheckHash(msg.GetFilename(), msg.GetHash())

	default:
		panic("Invalid enum value in FileMessage!")
	}
}

func (sh *ShareHandler) HandleShare(msg *ShareMessageWrapper) {
	if msg.GetShareName() != sh.Name {
		LogObj.Println("Ignoring message meant for share", msg.GetShareName())
		return
	}

	switch msg.GetAction() {
	case light.ShareAction_LEAVING:
		sh.RemoveClient(msg.Sender())

	case light.ShareAction_ENTERING:
		sh.AddClient(msg.Sender())
	}
}

func (sh *ShareHandler) HandlePeer(msg *PeerMessageWrapper) {
	for _, m := range msg.GetShares() {
		if m == sh.Name {
			//Request peer to be connected to!!
		}
	}
}

func (sh *ShareHandler) Handle(msg Message) {
	switch msg.(type) {
	case *PeerMessageWrapper:
		sh.HandlePeer(msg.(*PeerMessageWrapper))

	case *ShareMessageWrapper:
		sh.HandleShare(msg.(*ShareMessageWrapper))

	case *FileMessageWrapper:
		sh.HandleFile(msg.(*FileMessageWrapper))

	default:
		panic("Invalid message type!!")
	}
}

func (sh *ShareHandler) CheckHash(name string, hash []byte) bool {
	return true
}

func (sh *ShareHandler) Stop() {
	sh.controlChannel <- 0
}
