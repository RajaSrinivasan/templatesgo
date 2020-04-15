package cmd

import (
	"github.com/spf13/cobra"
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
	serverPort = install.Ask("Server Port", serverPort)
	serverURL = install.Ask("Server URL", serverURL)
	serverURL = serverURL + ":" + serverPort
	serverCertFileName = install.Ask("Server Certificate Filename", serverCertFileName)
	pvtKeyFileName = install.Ask("Private Key filename", pvtKeyFileName)
	htmlPath = install.Ask("HTML Path", htmlPath)
}
