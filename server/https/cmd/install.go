package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.com/projtemplates/go/server/https/serve"
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

func copyFile(fromdir, filename string, destdir string) {

	inpfile, err := os.Open(path.Join(fromdir, filename))
	if err != nil {
		log.Printf("%s", err)
		return
	}
	defer inpfile.Close()
	outfile, err := os.Create(path.Join(destdir, filename))
	if err != nil {
		log.Printf("%s", err)
		return
	}
	defer outfile.Close()
	_, err = io.Copy(outfile, inpfile)
	if err != nil {
		log.Printf("%s", err)
	}
	log.Printf("Installed %s in %s", filename, destdir)
}

func copyDir(from, to string) {
	log.Printf("Installing files of %s to %s", from, to)
	files, err := ioutil.ReadDir(from)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		copyFile(from, file.Name(), to)
	}
}

func installHTMLFiles(htmlDir string) {
	exedir, _ := osext.ExecutableFolder()
	templdir := path.Join(exedir, "..", "templates")
	copyDir(templdir, htmlDir)
}

func Install(cmd *cobra.Command, args []string) {

	serverDir = install.Ask("Server directory", serverDir)
	err := os.MkdirAll(serverDir, os.ModePerm)
	if err != nil {
		log.Fatal("Creating Dir", err)
	}
	insttime := time.Now()
	reftime := insttime.Format(time.ANSIC)
	viper.Set("server.installed", reftime)
	install.SetInstallDate(reftime)
	viper.Set("server.toplevel", serverDir)

	err = os.MkdirAll(path.Join(serverDir, "etc"), os.ModePerm)
	if err != nil {
		log.Fatal("Creating etc dir", err)
	}

	err = os.MkdirAll(path.Join(serverDir, "log"), os.ModePerm)
	if err != nil {
		log.Fatal("Creating log dir", err)
	}

	htmlDir := path.Join(serverDir, "html")
	_, err = os.Stat(htmlDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(htmlDir, os.ModePerm)
		if err != nil {
			log.Fatal("Creating html dir ", err)
		}
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

	adminpwd := install.Ask("Admin Password (username: admin)", "admin")
	userpwd := install.Ask("User Password (username: user)", "user")
	adminpwdenc := install.Password("admin", adminpwd)
	viper.Set("users.admin", adminpwdenc)
	userpwdenc := install.Password("user", userpwd)
	viper.Set("users.user", userpwdenc)

	storekey := make([]byte, serve.StoreKeyLength)
	_, err = rand.Read(storekey)
	if err != nil {
		log.Fatal(err)
	}
	viper.Set("store.key", hex.EncodeToString(storekey))
	err = viper.SafeWriteConfigAs(cfgFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Saved Configuration file %s", cfgFilename)

	log.Printf("Generating certificates")
	install.GeneratePrivateKeypair(pvtKeyFileName)
	privatepemname := pvtKeyFileName + ".pvt.pem"
	install.CreateCert(privatepemname, serverCertFileName)
	log.Printf("Setting up website")
	installHTMLFiles(htmlPath)
}
