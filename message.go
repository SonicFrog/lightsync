package main

import (
	"errors"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"io"
	"lightsync/proto"
)

const (
	FileMessageOP  byte = 0x1
	ShareMessageOP      = 0x2
	PeerMessageOP       = 0x3
)

type Message interface {
	SetSender(client Client)
	Sender() Client
	WriteTo(writer io.Writer) error
}

type MessageHandler interface {
	HandOver(msg Message)
}

type MessageWrapper struct {
	op     byte
	sender Client
}

type FileMessageWrapper struct {
	MessageWrapper
	*light.FileMessage
}

type PeerMessageWrapper struct {
	MessageWrapper
	*light.PeerMessage
}

type ShareMessageWrapper struct {
	MessageWrapper
	*light.ShareMessage
}

func (w *MessageWrapper) SetSender(sender Client) {
	w.sender = sender
}

func (w *MessageWrapper) Sender() (c Client) {
	return w.sender
}

func (w *FileMessageWrapper) WriteTo(writer io.Writer) (err error) {
	data, err := proto.Marshal(w)

	if err != nil {
		return
	}

	_, err = writer.Write(data)

	return
}

func (w *PeerMessageWrapper) WriteTo(writer io.Writer) (err error) {
	var t byte = PeerMessageOP
	data, err := proto.Marshal(w)
	n := len(data)

	if err != nil {
		return
	}

	err = binary.Write(writer, binary.BigEndian, n)

	if err != nil {
		return
	}

	err = binary.Write(writer, binary.BigEndian, t)

	if err != nil {
		return
	}

	_, err = writer.Write(data)

	return
}

func (w *ShareMessageWrapper) WriteTo(writer io.Writer) (err error) {
	data, err := proto.Marshal(w)

	if err != nil {
		return
	}

	_, err = writer.Write(data)

	return
}

func ReadMessage(reader io.Reader) (msg Message, err error) {
	var length int32
	var mtype byte

	err = binary.Read(reader, binary.BigEndian, &mtype)

	if err != nil {
		return
	}

	err = binary.Read(reader, binary.BigEndian, &length)

	if err != nil {
		return
	}

	data := make([]byte, length)

	n, err := reader.Read(data)

	if err != nil || int32(n) != length {
		LogObj.Println("Message reading error:", err)
		return
	}

	switch mtype {
	case FileMessageOP:
		pb := &light.FileMessage{}
		err = proto.Unmarshal(data, pb)
		msg = &FileMessageWrapper{MessageWrapper{FileMessageOP, nil}, pb}

	case ShareMessageOP:
		pb := &light.ShareMessage{}
		err = proto.Unmarshal(data, pb)
		msg = &ShareMessageWrapper{MessageWrapper{ShareMessageOP, nil}, pb}

	case PeerMessageOP:
		pb := &light.PeerMessage{}
		err = proto.Unmarshal(data, pb)
		msg = &PeerMessageWrapper{MessageWrapper{PeerMessageOP, nil}, pb}

	default:
		return nil, errors.New("Invalid message type received!" + string(mtype))
	}

	return
}
