package main

import (
	"crypto/tls"
	"errors"
)

type PeerInfo struct {
	address     string
	port        string
	fingerprint string
}

type TLSPeerConnector struct {
	clients  <-chan Client
	info     chan<- *PeerInfo
	ctrl     chan int "Control channel used to stop the peer finder"
	accepter ClientAccepter
	tlsConf  *tls.Config
}

type PeerConnector interface {
	Clients() <-chan *Client
	Info() chan<- *PeerInfo
}

func NewTLSPeerConnector(cfg *tls.Config, accept ClientAccepter) (*TLSPeerConnector, error) {
	info, clients := make(chan *PeerInfo, 10), make(chan Client, 10)

	pf := &TLSPeerConnector{
		info:     info,
		clients:  clients,
		ctrl:     make(chan int),
		accepter: accept,
		tlsConf:  cfg,
	}

	go pf.internal(info, clients)

	return pf, nil
}

func (pi *PeerInfo) Address() string {
	return pi.address
}

func (pi *PeerInfo) Port() string {
	return pi.port
}

func (pi *PeerInfo) Fingerprint() string {
	return pi.fingerprint
}

func (pf *TLSPeerConnector) internal(info <-chan *PeerInfo, clients chan<- Client) {
	for {
		select {
		case pi := <- info:
			client, err := pf.dial(pi.Address(), pi.Port(), pi.Fingerprint())

			if err != nil {
				return
			}

			clients <- client

		case <-pf.ctrl:
			return
		}
	}
}

func (cn *TLSPeerConnector) Clients() <-chan Client {
	return cn.clients
}

func (cn *TLSPeerConnector) Info() chan<- *PeerInfo {
	return cn.info
}

func (pf *TLSPeerConnector) Stop() {
	pf.ctrl <- 0
}

func (pf *TLSPeerConnector) dial(address, port, fingerprint string) (c Client, err error) {
	conn, err := tls.Dial("tcp", address+":"+port, pf.tlsConf)

	if err != nil {
		return
	}

	c = NewClient(fingerprint, conn)

	k := c.Key()

	if err != nil {
		LogObj.Println(err)
		return
	}

	if fingerprint != KeyFingerprint(k) {
		err = errors.New(conn.RemoteAddr().String() + " not using advertised key!")
		LogObj.Println(err)
	}

	if err != nil {
		LogObj.Println("Error initializing client!")
		return
	}

	return
}
