package install

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/RajaSrinivasan/rollpwd/salt"
)

func generate(t time.Time, nm string, pwd string) string {

	layout := "2006-01-02 15"
	ts := t.Format(layout)
	usersalt := salt.Generate(nm)
	salt := string(usersalt)
	h := md5.New()
	io.WriteString(h, salt)
	io.WriteString(h, ts)
	io.WriteString(h, nm)
	io.WriteString(h, pwd)
	pwdbytes := h.Sum(nil)
	pwdstr := hex.EncodeToString(pwdbytes)
	return pwdstr
}

func Password(nm, pwd string, insttime time.Time) string {
	userpwd := generate(insttime, nm, pwd)
	return userpwd
}

func Verify(nm, pwd string, pwdexp string, instime time.Time) bool {

	pwdenc := generate(instime, nm, pwd)
	fmt.Printf("User %s Password supplied %s computed %s\n", nm, pwdexp, pwdenc)
	timebasis := instime.Format(time.ANSIC)
	log.Printf("Time basis %s", timebasis)
	if strings.Compare(pwdenc, pwdexp) != 0 {
		return false
	}
	return true
}
