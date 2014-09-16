package main

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

/**
 * Crappy way to find new peers:
 * - Each peer should run a http server to share peers with other peers
 * - HTTP request to /peerfingerprint to get peer address and port
 **/
type AnnouncePeerFinder struct {
	announce  string
	Peers     chan *PeerInfo "The peer finder outputs PeerInfo in this channel"
	peerNames map[int]string
	run       bool
	mut       sync.Mutex
}

const (
	DefaultAnnounceTime    = 5
	DefaultAnnounceTimeOut = 3
)

func NewPeerFinder(announce string, peers []string) (out *AnnouncePeerFinder) {
	npeers := make(map[int]string)

	for index, peer := range peers {
		npeers[index] = peer
	}

	out = &AnnouncePeerFinder{
		announce:  announce,
		peerNames: npeers,
		Peers:     make(chan *PeerInfo, 10),
		run:       true,
	}

	go out.internal()

	return
}

//This functions could probably be removed considering we have to re-read config file anyway
func (pf *AnnouncePeerFinder) AddPeer(fingerprint string) {
	pf.mut.Lock()
	defer pf.mut.Unlock()

	for _, peer := range pf.peerNames {
		if peer == fingerprint {
			return
		}
	}

	pf.peerNames[len(pf.peerNames)] = fingerprint
}

func (pf *AnnouncePeerFinder) RemovePeer(fingerprint string) {
	pf.mut.Lock()

	defer pf.mut.Unlock()

	for index, peer := range pf.peerNames {
		if peer == fingerprint {
			delete(pf.peerNames, index)
			return
		}
	}
}

func (pf *AnnouncePeerFinder) internal() {
	client := &http.Client{
		Timeout: DefaultAnnounceTimeOut * time.Second,
	}

	defer close(pf.Peers)

	for pf.run {
		for _, name := range pf.peerNames {
			resp, err := client.Get(pf.announce + "/" + name)

			if err != nil {
				LogObj.Println("PeerFinder:", err)
				continue
			}

			data, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				LogObj.Println("PeerFinder:", err)
				continue
			}

			addr := string(data)

			LogObj.Println("Found peer at " + addr)
		}
		time.Sleep(DefaultAnnounceTime * time.Second)
	}
}

func (pf *AnnouncePeerFinder) Stop() {
	pf.run = false
}
