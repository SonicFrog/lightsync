package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	HardCodedPassword string = "Dummypassword"

	DefaultCertPath string = "/home/ars3nic/.config/lightsync.cert"
	DefaultKeyPath  string = "/home/ars3nic/.config/lightsync.key"

	DefaultKeyLength int = 2048
)

func DefaultTLSConfig() (cfg *tls.Config, err error) {
	return TLSConfig(DefaultCertPath, DefaultKeyPath)
}

func TLSConfig(certpath, keypath string) (cfg *tls.Config, err error) {
	clearcert, err := ioutil.ReadFile(certpath)

	if err != nil {
		return nil, errors.New("Could not load certificate in " + certpath)
	}

	clearkey, err := ioutil.ReadFile(keypath)

	if err != nil {
		return nil, errors.New("Could not load key in " + keypath)
	}

	cert, err := tls.X509KeyPair(clearcert, clearkey)

	if err != nil {
		return nil, err
	}

	cfg = &tls.Config{
		Rand:         rand.Reader,
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,
	}

	return
}

func KeyFingerprint(pub *rsa.PublicKey) (fp string) {
	hash := sha1.Sum(pub.N.Bytes())

	fp = hex.EncodeToString(hash[:])

	return
}

func TLSListener(config *tls.Config) (err error) {
	ln, err := tls.Listen("tcp", "localhost:12000", config)

	if err != nil {
		return
	}

	for {
		var conn net.Conn

		conn, err = ln.Accept()

		if err != nil {
			LogObj.Println("Could not accept connection: ", err)
			return
		}

		go TLSClientAcceptor(conn)
	}

	return
}

func TLSClientAcceptor(conn net.Conn) (err error) {
	tlscon, ok := conn.(*tls.Conn)

	if !ok {
		panic("TLSClientAcceptor has no use for classic connections!")
	}

	state := tlscon.ConnectionState()

	peerKey := state.PeerCertificates[0].PublicKey

	rsaPeerKey, ok := peerKey.(*rsa.PublicKey)

	if !ok {
		//TODO: Change this to clean handling
		panic("Remote peer is not using RSA!!")
	}

	LogObj.Println("Connection from peer", KeyFingerprint(rsaPeerKey))

	//TODO: Validate peer against authorized public keys fingerprints

	return
}

func GenerateAndWriteTLSCertificate() (cert_bytes []byte, err error) {
	var priv *rsa.PrivateKey
	var pub *rsa.PublicKey

	priv, _ = rsa.GenerateKey(rand.Reader, 4096)
	pub = &priv.PublicKey

	sn, err := rand.Int(rand.Reader, big.NewInt(65536))

	cn := make([]byte, 50)
	org := make([]byte, 50)
	orgu := make([]byte, 50)

	_, err = rand.Reader.Read(org)

	if err != nil {
		return
	}

	_, err = rand.Reader.Read(cn)
	if err != nil {
		return
	}

	_, err = rand.Reader.Read(orgu)
	if err != nil {
		return
	}

	cert := &x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Country:            []string{base64.StdEncoding.EncodeToString(cn)},
			Organization:       []string{base64.StdEncoding.EncodeToString(org)},
			OrganizationalUnit: []string{base64.StdEncoding.EncodeToString(orgu)},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA: true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	cert_bytes, err = x509.CreateCertificate(rand.Reader, cert, cert, pub, priv)

	if err != nil {
		LogObj.Println("Unable to create certificate:", err)
		return
	}

	cfile, err := os.OpenFile(DefaultCertPath, os.O_CREATE|os.O_WRONLY, 0750)

	if err != nil {
		return
	}

	defer cfile.Close()

	err = pem.Encode(cfile, &pem.Block{Type: "CERTIFICATE", Bytes: cert_bytes})

	if err != nil {
		LogObj.Println("Could not write certificate to", DefaultCertPath,
			": ", err)
		return
	}

	priv_bytes := x509.MarshalPKCS1PrivateKey(priv)

	kfile, err := os.OpenFile(DefaultKeyPath, os.O_CREATE|os.O_WRONLY, 0750)

	if err != nil {
		LogObj.Println("Unable to write private key: ", err)
		return
	}

	defer kfile.Close()

	err = pem.Encode(kfile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: priv_bytes})

	if err != nil {
		return
	}

	return
}
