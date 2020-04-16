package install

import (
	"log"
	"testing"
)

func TestCreateCert(t *testing.T) {
	err := CreateCert("server.pvt.pem", "certfile")
	if err != nil {
		t.Error(err)
	}
	log.Printf("Created cert file")
}
