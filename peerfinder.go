package main

import (
	"errors"
	"crypto/rsa"
	"crypto/tls"
)

type PeerFinder interface {
	FindRoutine() chan<- *PeerInfo
	ClientOutput() <-chan *Client
}

type PeerInfo struct {
	address     string
	port        string
	fingerprint string
}

type TLSPeerFinder struct {
	clientOutput chan *Client
	infoInput    chan *PeerInfo
	ctrl         chan int
	tlsConf      *tls.Config
}

func NewTLSPeerFinder(cfg *tls.Config, accept ClientAccepter) (pf PeerFinder, err error) {
	pf = &TLSPeerFinder{
		infoInput:    make(chan *PeerInfo, 10),
		clientOutput: make(chan *Client, 10),
	}
	return
}

func StartPeerFinder(pf PeerFinder) chan<- *PeerInfo {
	return pf.FindRoutine()
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

func (pf *TLSPeerFinder) FindRoutine() chan<- *PeerInfo {
	go pf.internal()

	return pf.infoInput
}

func (pf *TLSPeerFinder) internal() {
	for {
		select {
		case pi := <-pf.infoInput:
			client, err := pf.Dial(pi.Address(), pi.Port(), pi.Fingerprint())

			if err != nil {
				return
			}

			AddClient(client)

		case <-pf.ctrl:
			return
		}
	}
}

func (pf *TLSPeerFinder) Stop() {
	pf.ctrl <- 0
}

func (pf *TLSPeerFinder) Dial(address, port, fingerprint string) (c *Client, err error) {
	conn, err := tls.Dial("tcp", address+":"+port, pf.tlsConf)

	if err != nil {
		return
	}

	state := conn.ConnectionState()

	peerKey, ok := state.PeerCertificates[0].PublicKey.(rsa.PublicKey)

	if !ok {
		err = errors.New("Remote peer " + conn.RemoteAddr().String() +" not using RSA")
		return
	}

	if fingerprint != KeyFingerprint(&peerKey) {
		err = errors.New(conn.RemoteAddr().String() + " not using advertised key!")
		LogObj.Println(err)
	}

	c = NewClient(fingerprint, conn, &peerKey)

	return
}

func (pf *TLSPeerFinder) ClientOutput() <-chan *Client {
	return pf.clientOutput
}
