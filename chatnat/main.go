package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"code.google.com/p/go.crypto/otr"
)

func init() {
	TRACE = os.Getenv("TRACE") == "true"
	flag.StringVar(&waddellAddr, "waddell", "", "waddell signaling service address (required)")
	flag.StringVar(&serverID, "id", "", "id of the peer to connect to (optional)")
}

var (
	keyFile  = flag.String("key", "keyfile", "The private keyfile to use.")
	genKey   = flag.Bool("gen", false, "generate a new private key")
	keyBytes []byte

	otrPrivKey   otr.PrivateKey
	otrConv      otr.Conversation
	otrSecChange otr.SecurityChange

	waddellAddr string
	serverID    string
	VERBOSE     bool
	TRACE       bool
)

func main() {
	flag.Parse()

	if *genKey {
		GenKey()
	}

	if waddellAddr == "" {
		usage()
	}

	var err error
	keyBytes, err = ioutil.ReadFile(*keyFile)
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

	log.Printf("using waddell server: %s\n", waddellAddr)

	if len(serverID) > 1 {
		log.Printf("attempting to connect to: %s\n", serverID)
		Connect(serverID)
	} else {
		Serve()
	}

	// wait until we exit.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	<-sigc
}

func usage() {
	t := `usage: %s --waddell <address> [--id <id>]

Connects via waddell signaling server at <address>.
No --id puts %s in listening mode, it will output its <id>.
Specify --id <id> to connect to the corresponsing process.
`
	procname := os.Args[0]
	fmt.Fprintf(os.Stderr, t, procname, procname)
	os.Exit(1)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error: %q\n", err)
		os.Exit(1)
	}
}

func trace(s string, vals ...interface{}) {
	if TRACE {
		log.Printf(s, vals...)
	}
}
