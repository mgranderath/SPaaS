// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a stopped application",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverDefined()
		tokenDefined()
		token := viper.GetString("token")
		client := &http.Client{Transport: tr}
		url := viper.GetString("url") + "/api/app/" + args[0] + "/start"
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		isLoggedIn(res)
		scanner := bufio.NewScanner(res.Body)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := Application{}
			err := json.Unmarshal([]byte(scanner.Text()), &line)
			if err != nil {
				color.Red(err.Error())
			}
			switch line.Type {
			case "info":
				color.Yellow(fmt.Sprintf("%-14v", "Info:") + line.Message)
			case "error":
				color.Red(fmt.Sprintf("%-14v", "Error:") + line.Message)
			case "success":
				color.Green(fmt.Sprintf("%-14v", "Success:") + line.Message)
				if line.Extended != nil {
					for _, value := range line.Extended {
						color.Green("    " + fmt.Sprintf("%-14v", value.Key) + value.Value.(string))
					}
				}
			default:
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
