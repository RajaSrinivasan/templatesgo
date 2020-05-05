package install

import (
	"fmt"
	"testing"
	"time"
)

const timebasis = "Tue Apr 28 15:14:54 2020"

func TestGenerate(t *testing.T) {
	basistime := time.Now().Format(time.ANSIC)
	fmt.Printf("Time Basis : %s\n", basistime)

	ap := Password("admin", "admin")
	fmt.Printf("Admin password %s\n", ap)
	ap = Password("user", "user")
	fmt.Printf("User password %s\n", ap)

}

func TestVerify(t *testing.T) {

	SetInstallDate(timebasis)

	adminpwd := "0df596012473bf23af0d059121033013"
	status := Verify("admin", "admin", adminpwd)
	fmt.Printf("Verification Admin password %v\n", status)

	userpwd := "025acc827008345b90a93d3704276b1b"
	status = Verify("user", "user", userpwd)
	fmt.Printf("Verification User password %v\n", status)

}
