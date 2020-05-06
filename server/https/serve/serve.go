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
	exptime, ok := sess.Values["expirytime"]
	if ok {
		log.Printf("Session Expiry %s", exptime)
		exptimeval, err := time.Parse(time.RFC822, exptime.(string))
		if err != nil {
			log.Printf("Error converting expiry time %s", err)
		}
		//exptimeval = exptimeval.Local()
		reftime := time.Now().Local()
		expired := reftime.After(exptimeval)
		log.Printf("%s %s", reftime.Format(time.RFC822), exptimeval.Format(time.RFC822))
		if expired {
			log.Printf("Session Validity has expired. Will force a relogin")
			log.Printf("%s %s", reftime.Format(time.RFC822), exptimeval.Format(time.RFC822))
		}
	} else {
		log.Printf("Unable to determine session expiry time")
	}
	return true
}

func validateLogin(un, pw string, w http.ResponseWriter, r *http.Request) bool {
	viper.SetConfigFile(sysconfigFilename)
	usrkey := fmt.Sprintf("users.%s", un)
	exppwd := viper.GetString(usrkey)
	log.Printf("User %s Password %s Encrypted %s.", un, pw, exppwd)
	sess, err := store.Get(r, "topr")
	if err != nil {
		log.Printf("Unable to create a new session from store.\n%s", err)
		sess, err = store.New(r, "topr")
		if err != nil {
			log.Printf("Even recreating session failed")
		}
		return false
	}

	gotpwd := install.Verify(un, pw, exppwd)
	if !gotpwd {
		log.Printf("Password did not verify. Invalidating session")
		sess.Values["validated"] = false
		err = sess.Save(r, w)
		return false
	}
	log.Printf("Password verified. Validating session")
	timenow := time.Now().Local()
	exptime := timenow.Add(time.Hour * 1)
	exptimestr := exptime.Format(time.RFC822)
	sess.Values["validated"] = true
	sess.Values["expirytime"] = exptimestr
	err = sess.Save(r, w)
	return true
}

func renewSession(w http.ResponseWriter, r *http.Request) {
	sess, err := store.Get(r, "topr")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	timenow := time.Now().Local()
	exptime := timenow.Add(time.Hour * 1)
	exptimestr := exptime.Format(time.RFC822)
	sess.Values["expirytime"] = exptimestr
	sess.Values["validated"] = true
	err = sess.Save(r, w)
}

func getStats(w http.ResponseWriter, r *http.Request) {

	v := validUser(w, r)
	if !v {
		log.Printf("Not a valid user")
		return
	}
	renewSession(w, r)
	stats := templates.Lookup("stats.html")
	if stats == nil {
		log.Printf("Unable to locate stats.html")
		w.Write([]byte("Stats"))
		return
	}
	timenow := time.Now().Format(time.RFC822)
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
	renewSession(w, r)

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

	getIndex(w, r)
	log.Printf("Served login")

}

func ProvideService(certfn, pvtkeyfn, hostnport string, htmlpath string, cfgfilename string) {

	sysconfigFilename = cfgfilename

	startTime = time.Now().Format(time.RFC822)
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

func timeCompare() {
	year2000 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	year3000 := time.Date(3000, 1, 1, 0, 0, 0, 0, time.Local)

	isYear2000BeforeYear3000 := year2000.Before(year3000) // True
	isYear3000BeforeYear2000 := year3000.Before(year2000) // False

	timenow := time.Now().Local()
	b1 := timenow.After(year2000)
	b2 := timenow.After(year3000)
	fmt.Printf("b1 %v b2 %v\n", b1, b2)

	reftime, _ := time.Parse(time.RFC822, "Tue May  5 18:19:04 2020")
	b1 = timenow.After(reftime.Local())
	b2 = timenow.Before(reftime.Local())
	fmt.Printf("After %v Before %v\n", b1, b2)
	fmt.Printf("year2000.Before(year3000) = %v\n", isYear2000BeforeYear3000)
	fmt.Printf("year3000.Before(year2000) = %v\n", isYear3000BeforeYear2000)

}
