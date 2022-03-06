package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var configPath string

var rootCMD = &cobra.Command{
	Use:   "ip-locator",
	Short: "ip-locator is a tool to lookup IP addresses",
}

func init() {
	rootCMD.PersistentFlags().StringVarP(&configPath, "configs-path", "c", "env.yaml", "path to configs directory")
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
