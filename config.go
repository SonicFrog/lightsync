package main

import (
	"crypto/rsa"
	"fmt"
)

type Configuration struct {
	NodeName string
	NodeKey  rsa.PrivateKey
}

func (c *Configuration) init() {
	fmt.Println("Initiated config")
}