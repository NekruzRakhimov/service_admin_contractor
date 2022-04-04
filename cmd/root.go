package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// RootCmd является командой, вызываемой по-умолчанию
var RootCmd = &cobra.Command{
	Use:   "service_admin_contractor",
	Short: "`service_admin_contractor` microservice provides admin contractor functionality",
	Long:  "`service_admin_contractor` microservice provides admin contractor functionality",
}

// Execute вызывается единожды из main.main().
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
	})
}
