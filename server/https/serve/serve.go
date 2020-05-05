package serve

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"gitlab.com/projtemplates/go/server/install"
	"gitlab.com/projtemplates/go/server/version"
)

var startTime string
var templates *template.Template
var store *sessions.CookieStore
var sysconfigFilename string

var StoreKey []byte

const StoreKeyLength = 32

type Info struct {
	Title string
}

type Stats struct {
	Title       string
	Hostname    string
	TimeStart   string
	TimeNow     string
	BuildTime   string
	Repo        string
	Branch      string
	ShortCommit string
	LongCommit  string
	Version     string
}

func validUser(w http.ResponseWriter, r *http.Request) bool {
	sess, err := store.Get(r, "topr")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	_, ok := sess.Values["validated"]
	if !ok {
		log.Printf("Not a validated session")

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return false
	}
	return true
}

func validateLogin(un, pw string, w http.ResponseWriter, r *http.Request) bool {
	viper.SetConfigFile(sysconfigFilename)
	usrkey := fmt.Sprintf("users.%s", un)
	exppwd := viper.GetString(usrkey)
	tv, err := time.Parse(time.ANSIC, install.InstallDate)
	if err != nil {
		log.Printf("Install Time %s cannot be converted \n%s", install.InstallDate, err)
	}
	log.Printf("User %s Password %s Encrypted %s. Install Time %s", un, pw, exppwd, install.InstallDate)
	sess, err := store.Get(r, "topr")
	if err != nil {
		log.Printf("Unable to create a new session from store.\n%s", err)
		sess, err = store.New(r, "topr")
		if err != nil {
			log.Printf("Even recreating session failed")
		}
		return false
	}

	gotpwd := install.Verify(un, pw, exppwd, tv)
	if !gotpwd {
		log.Printf("Password did not verify. Invalidating session")
		sess.Values["validated"] = false
		err = sess.Save(r, w)
		return false
	}
	log.Printf("Password verified. Validating session")

	sess.Values["validated"] = true
	err = sess.Save(r, w)
	return true
}

func getStats(w http.ResponseWriter, r *http.Request) {

	v := validUser(w, r)
	if !v {
		log.Printf("Not a valid user")
		return
	}

	stats := templates.Lookup("stats.html")
	if stats == nil {
		log.Printf("Unable to locate stats.html")
		w.Write([]byte("Stats"))
		return
	}
	timenow := time.Now().Format(time.ANSIC)
	hostname, _ := os.Hostname()
	var st = Stats{
		TimeStart:   startTime,
		TimeNow:     timenow,
		Hostname:    hostname,
		BuildTime:   version.BuildTime,
		Branch:      version.BranchName,
		Version:     fmt.Sprintf("%d.%d.%d", version.VersionMajor, version.VersionMinor, version.VersionBuild),
		Repo:        version.RepoURL,
		ShortCommit: version.ShortCommitId,
		LongCommit:  version.LongCommitId,
	}

	err := stats.Execute(w, st)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Served stats")
}

func getIndex(w http.ResponseWriter, r *http.Request) {

	v := validUser(w, r)
	if !v {
		log.Printf("Not a valid user")
		return
	}
	index := templates.Lookup("index.html")
	if index == nil {
		log.Printf("Cannot find index.html")
		w.Write([]byte("index"))
		return
	}
	var info = Info{Title: "TOPR LLC"}
	err := index.Execute(w, info)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Served index")

}
func getTop(w http.ResponseWriter, r *http.Request) {

	index := templates.Lookup("login.html")
	if index == nil {
		log.Printf("Cannot find login.html")
		w.Write([]byte("login"))
		return
	}
	var info = Info{Title: "TOPR LLC"}
	err := index.Execute(w, info)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Served login")

}

func doLogin(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Printf("Form parse error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var un, pw string
	un = r.FormValue("username")
	pw = r.FormValue("password")
	log.Printf("Attempts to login from %s with password %s", un, pw)
	vlstat := validateLogin(un, pw, w, r)
	if !vlstat {
		log.Printf("Validation of Login failed")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//http.Redirect(w, r, "/index", http.StatusAccepted)
	getIndex(w, r)
	log.Printf("Served login")

}

func ProvideService(certfn, pvtkeyfn, hostnport string, htmlpath string, cfgfilename string) {

	sysconfigFilename = cfgfilename

	startTime = time.Now().Format(time.ANSIC)
	var err error
	store = sessions.NewCookieStore(StoreKey)
	store.MaxAge(0)
	templates, err = template.ParseGlob(htmlpath + "/*")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", getTop)
	r.HandleFunc("/index", getIndex)
	r.HandleFunc("/stats", getStats)

	r.HandleFunc("/a/login", doLogin)

	http.Handle("/", r)

	err = http.ListenAndServeTLS(
		hostnport,
		certfn,
		pvtkeyfn+".pvt.pem",
		r)
	log.Fatal(err)
}
