package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

type Stoppable interface {
	Stop()
}

const (
	BuiltinConfiguration string = "$HOME/.config/lightsync.conf"
)

var ServerConfig *tls.Config
var Config ConfigurationObject
var Shares map[string]Share
var Clients map[string]Client
var ClientSync sync.Mutex

var LogObj *log.Logger

func init() {
	var err error

	LogObj = log.New(os.Stdout, "lightsync", log.Ltime)
	LogObj.SetPrefix("lightsync ")

	defer func() {
		if err := recover(); err != nil {
			LogObj.Println("Error starting up lightsync:", err)
			os.Exit(1)
		}
	}()

	Shares = make(map[string]Share)
	Clients = make(map[string]Client)

	Config, err = NewJSONConfiguration(os.ExpandEnv(BuiltinConfiguration))

	if err != nil {
		panic("Could not load configuration!")
	}

	ServerConfig, err = DefaultTLSConfig()

	if err != nil {
		panic(err)
	}
}

func main() {
	LogObj.Printf("starting...\n")

	signalChannel := make(chan os.Signal, 10)

	signal.Notify(signalChannel, os.Kill, os.Interrupt)

	ln, err := NewTLSClientAccepter(ServerConfig, nil, AddClient)

	if err != nil {
		LogObj.Println(err)
		return
	}

	go SignalHandler(signalChannel, ln)
}

func SignalHandler(signals chan os.Signal, ln net.Listener) {
	for {
		sig := <-signals

		switch sig {
		case os.Interrupt:
			LogObj.Println("Shutting down...")
			ln.Close()
		}
	}
}

func AddClient(c Client) {
	ClientSync.Lock()
	if _, contains := Clients[c.Name()]; !contains {
		Clients[c.Name()] = c
	}
	ClientSync.Unlock()
}
