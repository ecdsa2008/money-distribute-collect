/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"money-distribute-collect/src/utils"

	"github.com/spf13/cobra"
)

// walletCmd represents the wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "HDWallet Tools",
	Long:  `HDWallet Tools to generate mnemonic, bip-44 sequence addresses, private key, public key, address, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("wallet called")
		randomMnemonicFlag := cmd.Flag("random-mnemonic")
		if randomMnemonicFlag != nil {
			mnemonic := utils.GenMnemonic()
			log.Println("mnemonic:", mnemonic)
			for i := 0; i < 20; i++ {
				log.Println("index:", i, "address:", utils.GetAddressFromPrivateKey(utils.GetPrivateKey(mnemonic, uint(i))).Hex())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.Flags().BoolP("random-mnemonic", "r", false, "Generate random mnemonic")

	walletCmd.AddCommand(&cobra.Command{
		Use: "cancel",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("wallet cancel called")
		},
	})
}
