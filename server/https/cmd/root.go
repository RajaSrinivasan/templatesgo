package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verbosityLevel int
var cfgfile string

var serverURL = "localhost"

var cfgFilename = "/val/server/etc/server.yaml"
var serverCertFileName = "/var/server/etc/servercert"
var pvtKeyFileName = "/var/server/etc/privatekey"
var htmlPath = "/var/server/etc/html"
var logFilesPath = "/var/server/log"
var serverPort = "443"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Server Template",
	Long: `
	Generic TLS server `,
	Run: Server,
}

// Server provides the service ie runs as a daemon.
func Server(cmd *cobra.Command, args []string) {
	log.Println("Starting the service")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgfile, "config", "", "config file name. Default is "+cfgFilename)
	rootCmd.PersistentFlags().IntVarP(&verbosityLevel, "verbose", "v", 0, "verbosity level 1 .. 16")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgfile != "" {
		cfgFilename = cfgfile
	}
	viper.SetConfigFile(cfgFilename)
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
		serverURL = viper.GetString("server.url")
		serverPort := viper.GetString("server.port")
		log.Printf("Server URL set to %s", serverURL)

		serverURL = serverURL + ":" + serverPort
		serverCertFileName = viper.GetString("server.certfile")
		pvtKeyFileName = viper.GetString("server.privatekey")
		htmlPath = viper.GetString("server.htmlpath")

	}
}
