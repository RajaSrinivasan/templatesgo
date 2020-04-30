package install

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

var InstallDate string
var StoreKey []byte

const StoreKeyLength = 128

func makeTemplate() x509.Certificate {

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	randomNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatal(err)
	}
	nameinfo := pkix.Name{
		Organization:       []string{"TOPR LLC."},
		CommonName:         "localhost",
		OrganizationalUnit: []string{"Engineering"},
		Country:            []string{"US"},
		Province:           []string{"Pennsylvania"},
		Locality:           []string{"Downingtown"},
	}
	certTemplate := x509.Certificate{
		SerialNumber:          randomNumber,
		Subject:               nameinfo,
		EmailAddresses:        []string{"admin@toprllc.com", "rs@toprllc.com"},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certTemplate.DNSNames = []string{"localhost", "localhost.local", "server.toprllc.com"}
	return certTemplate
}

func loadPrivate(fn string) *rsa.PrivateKey {
	filedata, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(filedata)
	if block == nil {
		log.Fatal("Not a valid PEM file")
	}
	if block.Type != "RSA PRIVATE KEY" {
		log.Fatal("Not an RSA Privatre Key")
	}

	pvtkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %s and created a private key", fn)
	return pvtkey
}

func saveCert(der []byte, fn string) {
	cert := &pem.Block{Type: "CERTIFICATE", Bytes: der}
	of, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	pem.Encode(of, cert)
	of.Close()
	log.Printf("Certificate saved to %s", fn)
}

func CreateCert(pvtkeyfn string, certfn string) error {

	log.Println("Creating Cert")
	makeTemplate()
	pvtkey := loadPrivate(pvtkeyfn)
	pubkey := pvtkey.PublicKey
	template := makeTemplate()

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &pubkey, pvtkey)
	if err != nil {
		log.Fatal(err)
	}
	saveCert(cert, certfn)
	return nil
}
