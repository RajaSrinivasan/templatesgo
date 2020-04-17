package serve

import (
	"testing"
)

func TestProvideService(t *testing.T) {
	ProvideService("../../config/certfile", "/Users/rajasrinivasan/.ssh/id_rsa", "localhost:9443", "../../config/html")
}
