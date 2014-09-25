package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

type Client struct {
	inputCh   chan Message
	outputCh  chan Message
	controlCh chan int
	key       *rsa.PublicKey
	conn      net.Conn
	name      string
}

var Config ConfigurationObject
var Shares map[string]Share
var Clients map[string]*Client
var Running = true

var LogObj *log.Logger

var FileMessageHandlers []MessageHandler
var PeerMessageHandlers []MessageHandler
var ShareMessageHandlers []MessageHandler

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

		name, err := ClientHandshake(conn)

		if err != nil {
			fmt.Printf("Error while handshaking with %s:\n",
				conn.RemoteAddr().String())
			fmt.Print(err)
			continue
		}

		client := NewClientHandler(name, conn)

		Clients[name] = &client
	}
}

func SignalHandler(signals chan os.Signal, ln *net.TCPListener) {
	for {
		sig := <-signals

		switch sig {
		case os.Interrupt:
			LogObj.Println("Shutting down...")
			ln.Close()
		}
	}
}

func NewClientHandler(name string, conn net.Conn) Client {
	input, output := make(chan Message, 10), make(chan Message, 10)

	c := Client{
		inputCh:  input,
		outputCh: output,
		conn:     conn,
		name:     name,
	}

	go c.ClientWriter(input, conn)
	go c.ClientReader(output, conn)

	return c
}

func (c *Client) WriteMessage(msg Message) {
	c.inputCh <- msg
}

func (c *Client) ReadMessage(msg Message) Message {
	return <-c.outputCh
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) Stop() {
	c.controlCh <- 0
}

func ClientHandshake(conn net.Conn) (name string, err error) {
	var msg Message



	msg, err = ReadMessage(conn)

	switch msg.(type) {

	}

	if err != nil {
		fmt.Printf("Error while handshaking with %s:\n",
			conn.RemoteAddr().String())
		fmt.Print(err)
		return
	}

	if err != nil {
		fmt.Printf("Error while handshaking with %s:\n",
			conn.RemoteAddr().String())
		fmt.Print(err)
		conn.Close()
		return
	}

	return
}


func (c *Client) ClientWriter(input <-chan Message, conn net.Conn) {
	for {
		select {
		case msg := <-input:
			err := msg.WriteTo(conn)

			if err != nil {
				fmt.Printf("Error while writing to client %s:\n",
					conn.RemoteAddr().String())
				return
			}

		case <-c.controlCh:
			LogObj.Println("Writer stopping for client ", c.Name)
			c.conn.Close()
		}
	}
}

func (c *Client) ClientReader(output chan<- Message, conn net.Conn) {
	for {
		msg, err := ReadMessage(conn)

		msg.SetSender(c)

		if err != nil {
			LogObj.Printf("Error while reading from client %s\n",
				conn.RemoteAddr().String())
			conn.Close()
			return
		}

		output <- msg
	}
}
