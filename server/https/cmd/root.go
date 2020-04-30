package cmd

import (
	"encoding/hex"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
	"gitlab.com/projtemplates/go/server/https/serve"
	"gitlab.com/projtemplates/go/server/install"
)

var verbosityLevel int
var cfgfile string

var serverURL = "localhost"
var serverPort = "443"

var serverDir string
var cfgFilename string

var serverCertFileName string
var pvtKeyFileName string
var htmlPath string
var logFilesPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Server Template",
	Long: `
	Generic TLS server `,
	Version: "0.1-0",
	Run:     Server,
}

func showConfiguration() {
	log.Printf("Config filename %s", cfgFilename)
	log.Printf("Server URL %s port %s", serverURL, serverPort)
	log.Printf("Cert file %s. Private Key file %s", serverCertFileName, pvtKeyFileName)
	log.Printf("HTML path %s", htmlPath)
	log.Printf("Log files path %s", logFilesPath)
}

// Server provides the service ie runs as a daemon.
func Server(cmd *cobra.Command, args []string) {
	log.Println("Starting the service")
	if verbosityLevel > 0 {
		showConfiguration()
	}
	serve.ProvideService(serverCertFileName, pvtKeyFileName, serverURL, htmlPath)
}

func loadConfigurations() {

	viper.SetConfigFile(cfgFilename)
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
		serverURL = viper.GetString("server.url")
		serverPort := viper.GetString("server.port")
		log.Printf("Server URL set to %s", serverURL)

		serverURL = serverURL + ":" + serverPort
		serverCertFileName = viper.GetString("server.cert")
		pvtKeyFileName = viper.GetString("server.pvtkey")
		htmlPath = viper.GetString("server.html")
		logFilesPath = viper.GetString("server.logfiles")
		install.InstallDate = viper.GetString("server.installed")
		log.Printf("Installed date set to %s", install.InstallDate)
		storekey := viper.Get("store.key")
		if storekey == nil {
			log.Fatal("No store key found")
		}
		storekeystr := viper.GetString("store.key")
		install.StoreKey, err = hex.DecodeString(storekeystr)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Restored store key")
	}
	if verbosityLevel > 0 {
		showConfiguration()
	}
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
	home, err := homedir.Dir()
	if err == nil {

		serverDir = path.Join(home, "server")

		cfgFilename = path.Join(serverDir, "etc", "server.yaml")
		htmlPath = path.Join(serverDir, "html")
		logFilesPath = path.Join(serverDir, "log")

	}

	if cfgfile != "" {
		cfgFilename = cfgfile
	}
	loadConfigurations()
}
