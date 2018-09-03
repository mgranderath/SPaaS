package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications created",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		serverDefined()
		tokenDefined()
		token := viper.GetString("token")
		client := &http.Client{Transport: tr}
		url := viper.GetString("url") + "/api/app"
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		scanner := bufio.NewScanner(res.Body)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := Application{}
			err := json.Unmarshal([]byte(scanner.Text()), &line)
			if err != nil {
				color.Red(err.Error())
			}
			color.Yellow(line.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
