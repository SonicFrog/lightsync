package main

import (
	"encoding/json"
	"os"
)

const (
	DefaultConfigFile string = "$HOME/.config/lightsync.json"
)

type Configuration struct {
	nodeName string
	keyPath  string
	certPath string
	shares   []ShareConfig
	clients  []ClientConfig
}

type ShareConfig struct {
	name                string
	path                string
	authorizedClientsID []string
}

type ClientConfig struct {
	name string
	id   string
}

type ConfigurationObject interface {
	NodeName() string
	CertPath() string
	KeyPath() string
}

func NewJSONConfiguration(filepath string) (c *Configuration, err error) {
	LogObj.Println("Initializing config...")

	jfile, err := os.Open(filepath)

	if err != nil {
		LogObj.Println("Unable to open config file:", err)
		return
	}

	defer jfile.Close()

	jdec := json.NewDecoder(jfile)

	err = jdec.Decode(c)

	if err != nil {
		LogObj.Println("Unable to parse config file:", err)
		return
	}

	LogObj.Println("Config read from", jfile.Name())

	return
}

func (c *Configuration) CertPath() string {
	return c.certPath
}

func (c *Configuration) KeyPath() string {
	return c.keyPath
}

func (c *Configuration) NodeName() string {
	return c.nodeName
}
