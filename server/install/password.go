package install

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"
	"time"

	"github.com/RajaSrinivasan/rollpwd/salt"
)

func generate(t time.Time, nm string, pwd string, s []byte) string {
	layout := "2006-01-02 15"
	ts := t.Format(layout)
	salt := string(s)
	h := md5.New()
	io.WriteString(h, salt)
	io.WriteString(h, ts)
	io.WriteString(h, nm)
	io.WriteString(h, pwd)
	pwdbytes := h.Sum(nil)
	pwdstr := hex.EncodeToString(pwdbytes)
	return pwdstr
}

func Password(nm, pwd string) string {
	usersalt := salt.Generate(nm)
	userpwd := generate(time.Now(), nm, pwd, usersalt)
	return userpwd
}

func Verify(nm, pwd string, pwdexp string, instime time.Time) bool {
	usersalt := salt.Generate(nm)
	pwdenc := generate(instime, nm, pwd, usersalt)
	// fmt.Printf("User %s Password %s\n", nm, pwdenc)
	if strings.Compare(pwdenc, pwdexp) != 0 {
		return false
	}
	return true
}
