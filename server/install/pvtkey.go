package install

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

const keySize = 4096

func GeneratePrivateKeypair(pairname string) {
	privatekey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		log.Fatal("Generating private key. ", err)
	}
	encodedpvtkey := x509.MarshalPKCS1PrivateKey(privatekey)
	var privatePem = pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: encodedpvtkey,
	}

	pvtpemfile, err := os.Create(pairname + ".pvt.pem")
	if err != nil {
		log.Fatal("Creating private pem file", pairname)
	}
	defer pvtpemfile.Close()
	err = pem.Encode(pvtpemfile, &privatePem)
	if err != nil {
		log.Fatal("Encoding private key ", err)
	}
	log.Printf("Created private key file %s.pvt.pem", pairname)

	encodedpublickey, err := x509.MarshalPKIXPublicKey(&privatekey.PublicKey)
	if err != nil {
		log.Fatal("Creating public key from private ", err)
	}
	publicpem := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: encodedpublickey,
	}

	pubpemfile, err := os.Create(pairname + ".pub.pem")
	if err != nil {
		log.Fatal("Creating public pem file", err)
	}
	defer pubpemfile.Close()

	err = pem.Encode(pubpemfile, &publicpem)
	if err != nil {
		log.Fatal("Encoding public pem", err)
	}
	log.Printf("Created public key %s.pub.pem", pairname)
}
