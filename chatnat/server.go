package main

import (
	"log"
	"net"

	nw "github.com/getlantern/nattywad"
)

func Serve() {
	nserver := nw.Server{OnSuccess: onServerSuccess}
	nserver.Configure(waddellAddr)
}

func onServerSuccess(laddr, raddr *net.UDPAddr) bool {
	log.Printf("connected %s to %s\n", laddr, raddr)
	if err := msgLoop(laddr, raddr); err != nil {
		log.Printf("failed to netcat: %s\n", err)
	}
	return true
}
