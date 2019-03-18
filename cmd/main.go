package main

import (
	"fmt"
	//"github.com/cielo24/cbt24/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Runs the login API call",
	Long: `
        Required Arguments: username, password
        Optional Arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		host, _ := cmd.Flags().GetString("host")
		version, _ := cmd.Flags().GetString("version")
		username, err := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		fmt.Printf("SOMETHING: %s", username)
		if err != nil {
			fmt.Printf("err: %s", err)
		}

		if url == "" {
			url = "api.cielo24.com"
		}
		if host == "" {
			host = "api.cielo24.com"
		}
		if version == "" {

		}
		if viper.GetString("username") != "" {
			username = viper.GetString("username")
		fmt.Printf("USERNAME: %s", username)
		}
		if viper.GetString("password") != "" {
			password = viper.GetString("password")
		}
		loginSuccess := util.Login(url, version, host, username, password)

		if loginSuccess != "" {
			fmt.Printf("Login successful. API Token: %s", loginSuccess)
		} else {
			fmt.Println("Login attempt failed")
		}
	},
}
