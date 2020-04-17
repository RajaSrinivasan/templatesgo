package cmd

import (
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/projtemplates/go/server/install"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the server",
	Long: `
	Installation requires the setup of the backend including creating the
	self signed certificates for the service. `,
	Run: Install,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func Install(cmd *cobra.Command, args []string) {

	serverDir = install.Ask("Server directory", serverDir)
	err := os.MkdirAll(serverDir, os.ModePerm)
	if err != nil {
		log.Fatal("Creating Dir", err)
	}
	viper.Set("server.toplevel", serverDir)

	err = os.MkdirAll(path.Join(serverDir, "etc"), os.ModePerm)
	if err != nil {
		log.Fatal("Creating etc dir", err)
	}

	err = os.MkdirAll(path.Join(serverDir, "log"), os.ModePerm)
	if err != nil {
		log.Fatal("Creating log dir", err)
	}

	err = os.MkdirAll(path.Join(serverDir, "html"), os.ModePerm)
	if err != nil {
		log.Fatal("Creating html dir", err)
	}

	cfgFilename = path.Join(serverDir, "etc", "server.yaml")

	serverCertFileName = path.Join(serverDir, "etc", "certfile")
	viper.Set("server.cert", serverCertFileName)

	pvtKeyFileName = path.Join(serverDir, "etc", "keypair")
	viper.Set("server.pvtkey", pvtKeyFileName)

	serverPort = install.Ask("Server Port", serverPort)
	viper.Set("server.port", serverPort)

	serverURL = install.Ask("Server URL", serverURL)
	viper.Set("server.URL", serverURL)

	htmlPath = path.Join(serverDir, "html")
	viper.Set("server.html", htmlPath)

	logFilesPath = path.Join(serverDir, "log")
	viper.Set("server.logfiles", logFilesPath)

	err = viper.SafeWriteConfigAs(cfgFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Saved Configuration file %s", cfgFilename)
	install.GeneratePrivateKeypair(pvtKeyFileName)
	privatepemname := pvtKeyFileName + ".pvt.pem"
	install.CreateCert(privatepemname, serverCertFileName)
}
