package serve

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getTop(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{
			"title": "Home Page",
		},
	)

}

func ProvideService(certfn, pvtkeyfn, hostnport string, htmlpath string) {
	log.Printf("Providing Service HTTPS")

	r := gin.Default()
	r.LoadHTMLGlob(htmlpath + "/*")
	r.GET("/", getTop)
	privatekeyfile := pvtkeyfn + ".pvt.pem"
	r.RunTLS(hostnport, certfn, privatekeyfile)

}
