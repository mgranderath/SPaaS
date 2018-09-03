package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func serverDefined() bool {
	if !viper.InConfig("url") {
		fmt.Println("URL not set in config. Run \"paas setup\"")
		os.Exit(1)
		return false
	}
	return true
}

func tokenDefined() bool {
	if !viper.InConfig("token") {
		fmt.Println("You're not logged in. Run \"paas login\"")
		os.Exit(1)
		return false
	}
	return true
}

func isLoggedIn(resp *http.Response) bool {
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("You're not logged in. Run \"paas login\"")
		os.Exit(1)
		return false
	}
	return true
}
