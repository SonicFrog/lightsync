package main

import (
	"crypto"
	"io"
)

type Client interface {
	Name() string
	Key() crypto.PublicKey
	WriteMessage(msg Message) error
	ReadMessage() (Message, error)
}

type ClientAttribs struct {
	Fingerprint string
	Key         crypto.PublicKey
	Reader      io.Reader
	Writer      io.Writer
	Input       chan Message
	Output      chan Message
}

type TLSClient struct {
	ClientAttribs
}

type SSHClient struct {
	ClientAttribs
}

func NewClient(fingerprint string, conn io.ReadWriter) (c Client) {
	return
}