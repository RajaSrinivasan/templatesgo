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
	"gitlab.com/projtemplates/go/server/version"
)

var startTime string
var templates *template.Template
var store *sessions.CookieStore

var StoreKey []byte

const StoreKeyLength = 128

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

func getGorillaStats(w http.ResponseWriter, r *http.Request) {

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

func getGorillaIndex(w http.ResponseWriter, r *http.Request) {

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
func getGorillaTop(w http.ResponseWriter, r *http.Request) {

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

func ProvideService(certfn, pvtkeyfn, hostnport string, htmlpath string) {

	startTime = time.Now().Format(time.ANSIC)
	var err error
	store = sessions.NewCookieStore(StoreKey)
	templates, err = template.ParseGlob(htmlpath + "/*")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", getGorillaTop)
	r.HandleFunc("/index", getGorillaIndex)
	r.HandleFunc("/stats", getGorillaStats)

	http.Handle("/", r)

	err = http.ListenAndServeTLS(
		hostnport,
		certfn,
		pvtkeyfn+".pvt.pem",
		r)
	log.Fatal(err)
}
