package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "used to setup the cli application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter URL of Server:")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		viper.Set("url", text)
		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
