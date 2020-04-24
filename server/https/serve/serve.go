package serve

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/projtemplates/go/server/version"
)

var startTime string

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

func ProvideService(certfn, pvtkeyfn, hostnport string, htmlpath string) {
	log.Printf("Providing Service HTTPS")
	startTime = time.Now().Format(time.ANSIC)
	r := gin.Default()
	r.LoadHTMLGlob(htmlpath + "/*")
	r.GET("/", getTop)
	r.GET("/stats", getStats)
	privatekeyfile := pvtkeyfn + ".pvt.pem"
	r.RunTLS(hostnport, certfn, privatekeyfile)

}
