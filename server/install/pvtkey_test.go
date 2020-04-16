package install

import (
	"testing"
)

func TestGeneratePrivateKeypair(t *testing.T) {
	GeneratePrivateKeypair("server")
}
