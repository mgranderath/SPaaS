package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"docker.io/go-docker/api/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a running application",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("token")
		client := &http.Client{}
		url := "http://" + viper.GetString("url") + ":" + viper.GetString("port") + "/api/app/" + args[0]
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		containerInfo := types.ContainerJSON{}
		err = json.Unmarshal(body, &containerInfo)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(strings.ToUpper(args[0]) + ":")
		fmt.Println("Container Name: " + containerInfo.Name)
		fmt.Println("Created: " + containerInfo.Created + "     " + "State: " + containerInfo.State.Status)
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
