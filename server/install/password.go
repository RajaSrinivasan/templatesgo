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

var reftime time.Time

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

func SetInstallDate(insdate string) {
	reftime, _ = time.Parse(time.ANSIC, insdate)
}

func Password(nm, pwd string) string {
	userpwd := generate(reftime, nm, pwd)
	return userpwd
}

func Verify(nm, pwd string, pwdexp string) bool {

	pwdenc := generate(reftime, nm, pwd)
	fmt.Printf("User %s Password supplied %s computed %s\n", nm, pwdexp, pwdenc)
	timebasis := reftime.Format(time.ANSIC)
	log.Printf("Time basis %s", timebasis)
	if strings.Compare(pwdenc, pwdexp) != 0 {
		return false
	}
	return true
}
