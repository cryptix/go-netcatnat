package main

import (
	"log"
	"net"
	"time"

	nw "github.com/getlantern/nattywad"
)

func Connect(id string) {
	nclient := nw.Client{
		DialWaddell:       dialWaddel,
		OnSuccess:         onClientSuccess,
		OnFailure:         onClientFailure,
		KeepAliveInterval: time.Second * 10,
	}

	sp := &nw.ServerPeer{ID: id, WaddellAddr: waddellAddr}
	nclient.Configure([]*nw.ServerPeer{sp})
}

func onClientSuccess(info *nw.TraversalInfo) {
	trace("Traversal Succeeded: %s\n", info)
	trace("Peer Country: %s\n", info.Peer.Extras["country"])
	trace("Peer ID: %s\n", info.Peer.ID)

	laddr := info.LocalAddr
	raddr := info.RemoteAddr
	log.Printf("connected %s to %s\n", laddr, raddr)
	if err := msgLoop(laddr, raddr); err != nil {
		log.Printf("failed to netcat: %s\n", err)
	}
}

func onClientFailure(info *nw.TraversalInfo) {
	trace("Traversal Failed: %s\n", info)
	trace("Peer Country: %s\n", info.Peer.Extras["country"])
	trace("Peer ID: %s\n", info.Peer.ID)
	log.Printf("failed to connect to %s\n", info.Peer.ID)
}

func dialWaddel(addr string) (net.Conn, error) {
	return net.Dial("tcp", waddellAddr)
}
