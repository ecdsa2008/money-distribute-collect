/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"log"
	"money-distribute-collect/src/utils"

	"github.com/spf13/cobra"
)

// displayCmd represents the display command
var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Display bip-44 sequence addresses info",
	Long:  `Display bip-44 sequence addresses info with address, private key, balance, nonce`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("display called")
		mnemonic, err := cmd.Flags().GetString("mnemonic")
		if err != nil {
			log.Panicln(errors.New("mnemonic is required"))
		}
		if mnemonic == "" {
			log.Panicln(errors.New("mnemonic is required"))
		}
		rpc, err := cmd.Flags().GetString("rpc")
		if err != nil {
			log.Panicln(errors.New("rpc is required"))
		}
		if rpc == "" {
			log.Panicln(errors.New("rpc is required"))
		}
		startIndex, err := cmd.Flags().GetUint("start-index")
		if err != nil {
			log.Panicln(errors.New("start-index is required"))
		}
		endIndex, err := cmd.Flags().GetUint("end-index")
		if err != nil {
			log.Panicln(errors.New("end-index is required"))
		}
		if endIndex < startIndex {
			log.Panicln(errors.New("end-index must bigger or equal to start-index"))
		}
		client, err := ethclient.Dial(rpc)
		if err != nil {
			log.Panicln(err)
		}
		for i := startIndex; i <= endIndex; i++ {
			accountPrivateKey := utils.GetPrivateKey(mnemonic, i)
			accountPublicKey := utils.GetPublicKey(accountPrivateKey)
			accountAddress := utils.GetAddressFromPublicKey(accountPublicKey)
			balance, err := client.BalanceAt(context.Background(), accountAddress, nil)
			if err != nil {
				log.Panicln(err)
			}
			humanReadableBalance := decimal.NewFromBigInt(balance, 0).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))).String()
			nonce, err := client.PendingNonceAt(context.Background(), accountAddress)
			if err != nil {
				log.Panicln(err)
			}

			//获取私钥的字节表示
			privateKeyBytes := crypto.FromECDSA(accountPrivateKey)
			//私钥Hex表示
			privateKeyHexString := fmt.Sprintf("%x", privateKeyBytes)
			fmt.Println("privateKeyHexString: ", "0x"+privateKeyHexString)
			log.Println("Account index: ", i, " Address: ", accountAddress.Hex(), " Private key: ", privateKeyHexString, " Balance: ", humanReadableBalance, " Nonce: ", nonce)
		}
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)
	displayCmd.Flags().StringP("mnemonic", "m", "", "Set your mnemonic")
	displayCmd.Flags().StringP("rpc", "r", "", "Set your rpc url")
	displayCmd.Flags().UintP("start-index", "", 1, "Set your start account index,default is 1")
	displayCmd.Flags().UintP("end-index", "", 20, "Set your end account index,must bigger or equal to start-index")
}
