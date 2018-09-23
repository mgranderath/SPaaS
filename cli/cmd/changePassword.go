package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// changePasswordCmd represents the changePassword command
var changePasswordCmd = &cobra.Command{
	Use:   "change-password",
	Short: "change the password of the SPaaS Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		serverDefined()
		tokenDefined()
		fmt.Println("Enter new password:")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password := string(bytePassword)
		fmt.Println()
		v := url.Values{}
		v.Add("password", password)
		token := viper.GetString("token")
		client := &http.Client{Transport: tr}
		url := viper.GetString("url") + "/change-password"
		req, _ := http.NewRequest("POST", url, strings.NewReader(v.Encode()))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		isLoggedIn(res)
		if res.StatusCode == http.StatusBadRequest {
			fmt.Println("Password has to be at least 8 characters")
			return
		}
		defer res.Body.Close()
		fmt.Println("Success changing password")
	},
}

func init() {
	rootCmd.AddCommand(changePasswordCmd)
}
