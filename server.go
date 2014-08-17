package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"log"
	"os/signal"
)

type Client struct {
	InputChannel  chan Message
	OutputChannel chan Message
	Name          string
	Key           rsa.PublicKey
}

var Config Configuration
var Shares map[string]Share
var Clients map[string]*Client
var Running = true

var LogObj *log.Logger

func main() {
	LogObj = log.New(os.Stdout, "lightsync", log.Ltime)
	LogObj.SetPrefix("lightsync ")
	LogObj.Printf("starting...\n")

	address, port := "localhost", "12000"

	addr, err := net.ResolveTCPAddr("tcp", address+":"+port)

	signalChannel := make(chan os.Signal, 10)

	signal.Notify(signalChannel, os.Kill, os.Interrupt)

	if err != nil {
		return
	}

	ln, err := net.ListenTCP("tcp", addr)

	go SignalHandler(signalChannel, ln)

	if err != nil {
		return
	}

	for {
		conn, err := ln.AcceptTCP()

		if err != nil {
			break
		}

		name, key, err := NetworkClientHandshake(conn)

		if err != nil {
			fmt.Printf("Error while handshaking with %s:\n",
				conn.RemoteAddr().String())
			fmt.Print(err)
			continue
		}

		client := NewClientHandler(name, conn, key)

		Clients[name] = &client
	}

	for _, c := range Clients {
		c.WriteMessage(CloseConnectionMessage{})
	}
}

func SignalHandler(signals chan os.Signal, ln *net.TCPListener) {
	for {
		sig := <-signals

		switch sig {
		case os.Interrupt:
			fmt.Println("Shutting down...")
			ln.Close()
		}
	}
}

func NewClientHandler(name string, conn *net.TCPConn, key rsa.PublicKey) Client {
	input, output := make(chan Message, 10), make(chan Message, 10)

	go NetworkClientWriter(input, conn)
	go NetworkClientReader(output, conn)

	return Client{input, output, name, key}
}

func (c *Client) WriteMessage(msg Message) {
	c.InputChannel <- msg
}

func (c *Client) ReadMessage(msg Message) Message {
	return <-c.OutputChannel
}

func NetworkClientHandshake(conn *net.TCPConn) (name string, key rsa.PublicKey, err error) {
	var msg Message

	key = rsa.PublicKey{N: big.NewInt(0), E: 0}

	msg = BaseAttribs{sender: Config.NodeName, entityName: Config.NodeName}

	msg, err = ReadMessage(conn)

	switch msg.(type) {
	case ClientHelloMessage:
		break

	default:
		msg = ProtocolViolationMessage{}
		msg.WriteTo(conn)
		err = errors.New("Protocol violation from " + conn.RemoteAddr().String())
		conn.Close()
		return
	}

	if err != nil {
		fmt.Printf("Error while handshaking with %s:\n",
			conn.RemoteAddr().String())
		fmt.Print(err)
		return
	}

	msg, err = ReadMessage(conn)

	if err != nil {
		fmt.Printf("Error while handshaking with %s:\n",
			conn.RemoteAddr().String())
		fmt.Print(err)
		conn.Close()
		return
	}

	return
}

func NetworkClientWriter(input chan Message, conn *net.TCPConn) {
	for Running {
		msg := <-input
		err := msg.WriteTo(conn)

		if err != nil {
			fmt.Printf("Error while writing to client %s:\n",
				conn.RemoteAddr().String())
			return
		}
	}
}

func NetworkClientReader(output chan Message, conn *net.TCPConn) {
	for Running {

		msg, err := ReadMessage(conn)

		if err != nil {
			fmt.Printf("Error while reading from client %s\n",
				conn.RemoteAddr().String())
			conn.Close()
			return
		}

		output <- msg
	}
}
