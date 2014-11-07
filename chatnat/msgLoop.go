package main

import (
	"log"
	"net"

	"code.google.com/p/go.crypto/otr"
	"github.com/tiborvass/uniline"
)

func msgLoop(laddr, raddr *net.UDPAddr) error {
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		return err
	}

	rdy := make(chan struct{})
	nextMessage := make(chan string)
	go func() {
		<-rdy
		log.Println("<OTR> Connection Ready.")
		scanner := uniline.DefaultScanner()
		for scanner.Scan("<OTR> ") {
			line := scanner.Text()
			if len(line) > 0 {
				nextMessage <- line
			}
		}

		checkErr(scanner.Err())
	}()

	for {
		buf := make([]byte, 0, 5000)
		n, err := conn.Read(buf)
		checkErr(err)

		trace("<RAW> %d bytes from Peer From Peer\nMsg:%q\n", n, string(buf))

		out, encrypted, otrSecChange, msgToPeer, err := otrConv.Receive(buf)
		checkErr(err)

		if len(out) > 0 {
			log.Printf("<OTR> %q\n", string(out))
		}

		if !encrypted {
			log.Println("<OTR> Conversation not yet encrypted!!!")
		}

		if len(msgToPeer) > 0 {
			log.Printf("<OTR> Transmitting %d messages.\n", len(msgToPeer))
			for _, msg := range msgToPeer {
				n, err := conn.Write(msg)
				checkErr(err)

				if n < len(msg) {
					log.Fatalln("<OTR> some bytes were not send to peer..")
				}
			}
		}

		switch otrSecChange {
		case otr.NoChange:
			if encrypted {
				msg := <-nextMessage
				msgToPeer, err := otrConv.Send([]byte(msg))
				checkErr(err)

				for _, msg := range msgToPeer {
					n, err := conn.Write(msg)
					checkErr(err)

					if n < len(msg) {
						log.Fatalln("<OTR> some bytes were not send to peer..")
					}
				}
			}

		case otr.NewKeys:
			log.Printf("<OTR> Key exchange completed. SSID:%x\n", otrConv.SSID)
			close(rdy)

		case otr.ConversationEnded:
			log.Println("<OTR> Conversation ended.")
			return nil

		default:
			log.Printf("<OTR> SMPState: %d - not yet implemented!... :(", otrSecChange)
		}
	}
}
