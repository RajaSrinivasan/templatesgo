package serve

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"gitlab.com/projtemplates/go/server/version"
)

var startTime string
var templates *template.Template

type Info struct {
	Title string
}

type Stats struct {
	Title       string
	Hostname    string
	TimeStart   string
	TimeNow     string
	BuilTTime   string
	Repo        string
	Branch      string
	ShortCommit string
	LongCommit  string

	BuildTime string
	Version   string
}

func getStats(c *gin.Context) {
	timenow := time.Now().Format(time.ANSIC)
	hostname, _ := os.Hostname()
	//c.String(http.StatusOK, "Hostname : %s\nStarted %s\nTime Now %s\n", hostname, startTime, timenow)

	c.HTML(
		http.StatusOK,
		"stats.html",
		gin.H{
			"title":        "Server Stats",
			"hostname":     hostname,
			"time_start":   startTime,
			"time_now":     timenow,
			"build_time":   version.BuildTime,
			"major":        version.VersionMajor,
			"minor":        version.VersionMinor,
			"build":        version.VersionBuild,
			"repo":         version.RepoURL,
			"branch":       version.BranchName,
			"short_commit": version.ShortCommitId,
			"long_commit":  version.LongCommitId},
	)
}

func getTop(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": "TOPR LLC.",
		},
	)

}

func ProvideGinService(certfn, pvtkeyfn, hostnport string, htmlpath string) {
	log.Printf("Providing Service using gin HTTPS")
	startTime = time.Now().Format(time.ANSIC)
	r := gin.Default()
	r.LoadHTMLGlob(htmlpath + "/*")
	r.GET("/", getTop)
	r.GET("/stats", getStats)
	privatekeyfile := pvtkeyfn + ".pvt.pem"
	r.RunTLS(hostnport, certfn, privatekeyfile)
}

func getGorillaStats(w http.ResponseWriter, r *http.Request) {
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
		BuilTTime:   version.BuildTime,
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

func getGorillaTop(w http.ResponseWriter, r *http.Request) {
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

func ProvideService(certfn, pvtkeyfn, hostnport string, htmlpath string) {
	log.Printf("Providing Service using gorilla mux HTTPS")
	startTime = time.Now().Format(time.ANSIC)
	var err error
	templates, err = template.ParseGlob(htmlpath + "/*")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", getGorillaTop)
	r.HandleFunc("/stats", getGorillaStats)

	http.Handle("/", r)

	err = http.ListenAndServeTLS(
		hostnport,
		certfn,
		pvtkeyfn+".pvt.pem",
		r)
	log.Fatal(err)
}
