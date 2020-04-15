package cmd

import (
	"log"

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
	cfgFilename = install.Ask("Configuration File Name", cfgFilename)

	serverCertFileName = install.Ask("Certificate file Name", serverCertFileName)
	viper.Set("server.cert", serverCertFileName)

	pvtKeyFileName = install.Ask("Private Key filename", pvtKeyFileName)
	viper.Set("server.pvtkey", pvtKeyFileName)

	serverPort = install.Ask("Server Port", serverPort)
	viper.Set("server.port", serverPort)

	serverURL = install.Ask("Server URL", serverURL)
	viper.Set("server.URL", serverURL)

	serverURL = serverURL + ":" + serverPort

	htmlPath = install.Ask("HTML Path", htmlPath)
	viper.Set("server.html", htmlPath)

	logFilesPath = install.Ask("Log Files Path", logFilesPath)
	viper.Set("server.logfiles", htmlPath)

	err := viper.SafeWriteConfigAs(cfgFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Saved Configuration file %s", cfgFilename)
}
