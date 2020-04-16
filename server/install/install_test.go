package install

import (
	"log"
	"path"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
)

func TestCreateCert(t *testing.T) {

	log.SetPrefix(t.Name())

	home, err := homedir.Dir()
	if err != nil {
		t.Error(err)
	}
	pvtkeyfn := path.Join(home, ".ssh", "id_rsa")
	err = CreateCert(pvtkeyfn, "certfile")
	if err != nil {
		t.Error(err)
	}

}
