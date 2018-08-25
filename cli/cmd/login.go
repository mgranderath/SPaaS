package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

type tokenWrap struct {
	Token string `json:"token"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter username:")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSuffix(username, "\n")
		fmt.Println("Enter password:")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password := string(bytePassword)
		fmt.Println()
		if !viper.InConfig("url") || !viper.InConfig("port") {
			fmt.Println("URL and Port not set in config. Run \"paas setup\"")
			return
		}
		v := url.Values{}
		v.Add("username", username)
		v.Add("password", password)
		url := "http://" + viper.GetString("url") + ":" + viper.GetString("port") + "/login"
		resp, err := http.PostForm(url, v)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Wrong login credentials")
			return
		}
		defer resp.Body.Close()
		token := &tokenWrap{}
		err = json.NewDecoder(resp.Body).Decode(token)
		if err != nil {
			fmt.Println(err.Error())
		}
		viper.Set("token", token.Token)
		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
