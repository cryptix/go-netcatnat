package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"code.google.com/p/go.crypto/otr"
	"github.com/cryptix/goSam"
)

// cumbersome - you need to paste the destination to connect to here....
const addr = "JvRSIkuotlBGEkvD~3s6Yr7nQMCXtaJdsvfejzHNhpXEu7rIX035zyGGOtKee9QhjIdnyMw0Dk7l0jgp289-XZ4lSWo-HYiJg9N67L9lBL0g4M1CA5xoOjkg4fFGqP8XeJoFZ-LzB2-fxH2yIe0gAU-Ye2ZWlW62pRgssW0zzTTZGWbm4Umc7Hf9Em8ZwnoaUCI4NxtiE2faNxKYOLnPd1LYrpurXhovaQruv1-1w~bzkOStXDw7DgG4oPRTwERcOAEzVN6RS2VroIzSewRZEANv1jHfWKmrAEP8JoR1I82sH8oU0pHcz0~oi5eEZAX1O3jMoq6qF-w6DB6tln0uOHNUWcdEO8CMhhLQ-DOCd0nYxr02O3iSCzbe~bV8Tc3dIuSidRYql9YB2F6tvkejP0aS57cZlY~5nlCe-e4wWk~4RaDjtuC~56R~ximCCbYJsGtbhxkg7rhvEi7pCRPiWB~PfFx4QjG7bA58tU2dVJQGNNY-F9dfGGVFgjMqQPHPAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABfLIsBGsXxvLzfTp7F~~zTqfPDayVq2yD7W6~kYSyl-oHevuROA8vCnkhqGsNNdeC"

var (
	keyFile = flag.String("key", "keyfile", "The private keyfile to use.")

	otrPrivKey   otr.PrivateKey
	otrConv      otr.Conversation
	otrSecChange otr.SecurityChange
)

func main() {
	flag.Parse()

	sam, err := goSam.NewDefaultClient()
	checkErr(err)
	defer sam.Close()

	log.Println("Client Created")

	keyBytes, err := ioutil.ReadFile(*keyFile)
	checkErr(err)

	rest, ok := otrPrivKey.Parse(keyBytes)
	if !ok {
		log.Fatalf("ERROR: Failed to parse private key %s\n", *keyFile)
	}
	if len(rest) > 0 {
		log.Fatalln("ERROR: data remaining after parsing private key")
	}

	otrConv.PrivateKey = &otrPrivKey
	otrConv.FragmentSize = 5000

	id, _, err := sam.CreateStreamSession("")
	checkErr(err)

	newC, err := goSam.NewDefaultClient()
	checkErr(err)

	err = newC.StreamConnect(id, addr)
	checkErr(err)

	log.Println("Stream connected. Sending OTR Query")
	fmt.Fprintf(newC.SamConn, "%s.", otr.QueryMessage)

	bufStdin := bufio.NewReader(os.Stdin)

	samReader := bufio.NewReader(newC.SamConn)
	msgLoop(newC.SamConn, samReader, bufStdin)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
