package main

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"

	"code.google.com/p/go.crypto/otr"
)

func GenKey() {
	newKey := new(otr.PrivateKey)

	newKey.Generate(rand.Reader)

	keyBytes := newKey.Serialize(nil)

	err := ioutil.WriteFile(*keyFile, keyBytes, 0700)
	checkErr(err)

	log.Printf("Done! Fingerprint: %v", newKey.Fingerprint())
	os.Exit(0)
}
