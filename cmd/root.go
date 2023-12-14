package cmd

import (
	"log"
)
import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "money-distribute-collect",
	Short: "Ethereum-Like native coin auto distribute and collect",
	Long:  `Ethereum-Like native coin auto distribute and collect, see https://github.com/ecdsa2008/money-distribute-collect`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("run mdc...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panicln(err)
	}
}
