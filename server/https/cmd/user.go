package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/projtemplates/go/server/install"
)

// showCmd represents the show command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User admin ",
	Long: `
	Create a new user, or delete or change password for a user.
	
	With no flags - create a new user.
	`,
	Run:  User,
	Args: cobra.MinimumNArgs(1),
}

const minLengthPassword = 4

var modify_opt bool
var delete_opt bool

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.PersistentFlags().BoolVarP(&modify_opt, "modify", "m", false, "Modify password of user")
	userCmd.PersistentFlags().BoolVarP(&delete_opt, "delete", "d", false, "Modify password of user")
}

func User(cmd *cobra.Command, args []string) {

	modviper := viper.New()
	modviper.SetConfigFile(cfgFilename)
	if err := modviper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var pwd string
	if !delete_opt {
		prompt := fmt.Sprintf("Password for user %s", args[0])
		pwd = install.Ask(prompt, "?")
		if len(pwd) < minLengthPassword {
			log.Printf("Invalid password")
			return
		}
	}
	nowkey := fmt.Sprintf("users.%s", args[0])
	insttime, _ := time.Parse(time.ANSIC, install.InstallDate)

	switch {
	case modify_opt:
		log.Printf("Modify password for user %s", args[0])
		nowval := modviper.Get(nowkey)
		if nowval == nil {
			log.Printf("User %s is not defined", args[0])
			return
		}
		newpwd := install.Password(args[0], pwd, insttime)
		modviper.Set(nowkey, newpwd)
	case delete_opt:
		log.Printf("Deleting user %s", args[0])
		nowval := modviper.Get(nowkey)
		if nowval == nil {
			log.Printf("User %s is not defined", args[0])
			return
		}
		modviper.Set(nowkey, "invalid")
	default:
		log.Printf("Create a new user %s", args[0])
		nowval := modviper.Get(nowkey)
		if nowval != nil {
			if strings.Compare(modviper.GetString(nowkey), "invalid") != 0 {
				log.Printf("User %s is already present with a password. Not adding", args[0])
				return
			}
		}
		newpwd := install.Password(args[0], pwd, insttime)
		modviper.Set(nowkey, newpwd)
	}
	err := modviper.WriteConfigAs(cfgFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Updated Configuration file %s", cfgFilename)

}
