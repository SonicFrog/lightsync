package main

import (
	"crypto/rsa"
	"log"
	"os"
	"testing"
	"time"
)

func TestGenerator(t *testing.T) {
	LogObj = log.New(os.Stdout, "lightsync: ", log.LstdFlags)

	_, err := GenerateAndWriteTLSCertificate()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestFingerprint(t *testing.T) {
	cfg, err := DefaultTLSConfig()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	rsaServerKey := cfg.Certificates[0].PrivateKey.(*rsa.PrivateKey)

	t.Log(KeyFingerprint(&rsaServerKey.PublicKey))
}

func TestTLSListener(t *testing.T) {
	cfg, err := DefaultTLSConfig()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	go TLSListener(cfg)

	time.Sleep(1)
}
