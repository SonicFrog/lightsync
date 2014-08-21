package main

import (
	"lightsync/proto"
)

type ShareHandler struct {
	Share
	RequestChannel chan Message
	ControlChannel chan int
	isActive       bool
}

func NewShareHandler(share Share, out chan Message) (sh *ShareHandler) {
	sh = &ShareHandler{
		Share:          share,
		RequestChannel: out,
		isActive:       true,
		ControlChannel: make(chan int),
	}

	go sh.HandleLocal()

	return
}

func (sh *ShareHandler) Active() bool {
	return sh.isActive
}

func (sh *ShareHandler) HandOver(msg Message) {
	sh.RequestChannel <- msg
}

func (sh *ShareHandler) HandleLocal() {
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

		case <-sh.ControlChannel:
			LogObj.Println("ShareHandler ", sh.Name, " stopping!")
			break
		}

	}
}

func (sh *ShareHandler) HandleFile(msg *FileMessageWrapper)  {
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
	sh.isActive = false
}
