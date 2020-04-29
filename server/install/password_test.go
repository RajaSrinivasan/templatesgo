package install

import (
	"fmt"
	"testing"
	"time"
)

const timebasis = "Tue Apr 28 15:14:54 2020"

func TestGenerate(t *testing.T) {
	fmt.Printf("Time Basis : %s\n", time.Now().Format(time.ANSIC))
	ap := Password("admin", "admin", time.Now())
	fmt.Printf("Admin password %s\n", ap)
	ap = Password("user", "user", time.Now())
	fmt.Printf("User password %s\n", ap)
}

func TestVerify(t *testing.T) {

	basis, _ := time.Parse(time.ANSIC, timebasis)
	adminpwd := "0df596012473bf23af0d059121033013"
	//fmt.Printf("Admin Password expected: %s\n", adminpwd)
	status := Verify("admin", "admin", adminpwd, basis)
	fmt.Printf("Verification Admin password %v\n", status)

	userpwd := "025acc827008345b90a93d3704276b1b"
	//fmt.Printf("User Password expected: %s\n", userpwd)
	status = Verify("user", "user", userpwd, basis)
	fmt.Printf("Verification User password %v\n", status)
}
