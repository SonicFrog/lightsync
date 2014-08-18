package main

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

type Message interface {
	Name() string
	Sender() *Client
	Share() string
	Payload() []byte
	WriteTo(write io.Writer) error
}

type FileMessage interface {
	FileName() string
	Hash() []byte
	Time() time.Time
}

type ShareMessage interface {
	ShareName() string
	Leaving() bool
}

type Handler interface {
	Handle(msg Message) error
}

type BaseAttribs struct {
	senderRef  *Client
	entityName string
	sender     string
}

type ShareAttribs struct {
	BaseAttribs
}

type FileAttribs struct {
	BaseAttribs
	fName string
	hash  []byte
	time  time.Time
}

type PeerAttribs struct {
	cert    []byte
	name    string
	address string
	version string
}

type ClientHelloMessage struct {
	BaseAttribs
}

type ServerHelloMessage struct {
	BaseAttribs
}

type ClientKeyMessage struct {
	BaseAttribs
}

type ServerKeyMessage struct {
	BaseAttribs
}

type ShareDiscoverMessage struct {
	ShareAttribs
}

type ShareACKMessage struct {
	ShareAttribs
}

type ShareLeaveMessage struct {
	ShareAttribs
}

type ShareLastModMessage struct {
	ShareAttribs
}

type FileRemoveMessage struct {
	FileAttribs
}

type FileUpdatedMessage struct {
	FileAttribs
}

type FileCreatedMessage struct {
	FileAttribs
}

type FileHashMessage struct {
	FileAttribs
}

type FilePartRequest struct {
	FileAttribs
}

type FilePartTransfer struct {
	FileAttribs
}

type DirectoryCreateMessage struct {
	FileAttribs
}

type DirectoryDiffMessage struct {
	FileAttribs
}

type PeerAnnounceMessage struct {
	BaseAttribs
}

type PeerRequestMessage struct {
	BaseAttribs
}

type CloseConnectionMessage struct {
	BaseAttribs
}

type KeepAliveMessage struct {
	BaseAttribs
}

type UnknownShareMessage struct {
	BaseAttribs
}

type UnauthorizedClientMessage struct {
	BaseAttribs
}

type ProtocolViolationMessage struct {
	BaseAttribs
}

const (
	//Handshake messages
	ClientHelloMessageOP byte = 0x00
	ServerHelloMessageOP byte = 0x01
	ClientKeyMessageOP   byte = 0x02
	ServerKeyMessageOP   byte = 0x03

	//Share related messages
	ShareDiscoverMessageOP = 0x10
	ShareACKMessageOP      = 0x11
	ShareLeaveMessageOP    = 0x12
	ShareLastModMessageOP  = 0x13

	//File state messages
	FileRemovedMessageOP = 0x20
	FileUpdatedMessageOP = 0x21
	FileCreatedMessageOP = 0x22
	FileHashMessageOP    = 0x23

	//File transfer messages
	FilePartRequestOP  = 0x30
	FilePartTransferOP = 0x31

	//Directory transfer messages
	DirectoryCreateMessageOP = 0x40
	DirectoryRemoveMessageOP = 0x41
	DirectoryDiffMessageOP   = 0x42

	//Peer exchange messages
	PeerAnnounceMessageOP = 0x60
	PeerRequestMessageOP  = 0x61

	//Connection control
	CloseConnectionMessageOP = 0x80

	//Error messages
	UnknownShareMessageOP       = 0x90
	UnauthorizedClientMessageOP = 0x91
	ProtocolViolationMessageOP  = 0x99
)

func ReadMessage(input io.Reader) (msg Message, err error) {
	var MType byte
	var CNameSize, PayloadLength int
	var ClientName string

	err = binary.Read(input, binary.LittleEndian, &MType)

	if err != nil {
		return
	}

	err = binary.Read(input, binary.LittleEndian, &CNameSize)

	if err != nil {
		return
	}

	bytes := make([]byte, CNameSize)
	read, err := input.Read(bytes)

	if err != nil || read != CNameSize {
		return
	}

	ClientName = string(bytes[:])

	err = binary.Read(input, binary.LittleEndian, &PayloadLength)

	if err != nil {
		return
	}

	bytes = make([]byte, PayloadLength)

	read, err = input.Read(bytes)

	if err != nil || read != PayloadLength {
		if read != PayloadLength {
			err = errors.New("Invalid message received!")
		}
		return
	}

	msg = BaseAttribs{entityName: ClientName, sender: ClientName}

	return
}

func (msg BaseAttribs) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg BaseAttribs) Sender() *Client {
	return msg.senderRef
}

func (msg BaseAttribs) Name() string {
	return msg.entityName
}

func (msg BaseAttribs) Payload() []byte {
	return []byte{}
}

func (msg BaseAttribs) Share() string {
	return ""
}

func (msg ClientHelloMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg ServerHelloMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg ClientKeyMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg ServerKeyMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg ShareDiscoverMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg ShareACKMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg ShareLeaveMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg ShareLastModMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg FileAttribs) Hash() []byte {
	return make([]byte, 0)
}

func (msg FileAttribs) Time() time.Time {
	return msg.time
}

func (msg FileAttribs) FileName() string {
	return msg.fName
}

func (msg FileRemoveMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg FileUpdatedMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg FileCreatedMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg FileHashMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg FilePartRequest) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg FilePartTransfer) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg DirectoryCreateMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg DirectoryDiffMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg PeerAnnounceMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg PeerRequestMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg CloseConnectionMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg UnknownShareMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg UnauthorizedClientMessage) WriteTo(writer io.Writer) (err error) {
	return
}

func (msg ProtocolViolationMessage) WriteTo(writer io.Writer) (err error) {
	return
}
