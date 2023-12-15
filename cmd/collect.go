/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"log"
	"math/big"
	"money-distribute-collect/src/utils"
	"time"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect all native coin from bip-44 sequence addresses",
	Long:  `Collect all native coin from bip-44 sequence addresses`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("collect called")
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

		masterIndex, err := cmd.Flags().GetUint("master-index")
		if err != nil {
			log.Panicln(errors.New("master-index is required"))
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

		masterAddress := utils.GetAddressFromPrivateKey(utils.GetPrivateKey(mnemonic, masterIndex))

		client, err := ethclient.Dial(rpc)
		if err != nil {
			log.Panicln(err)
		}

		networkID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Panicln(err)
		}

		// 从start-index开始，到end-index结束，依次检查每个账户的余额，如果大于0，则转账到当前账户的master账户
		for i := startIndex; i <= endIndex; i++ {
			if i == masterIndex {
				continue
			}
			accountPrivateKey := utils.GetPrivateKey(mnemonic, i)
			accountPublicKey := utils.GetPublicKey(accountPrivateKey)
			accountAddress := utils.GetAddressFromPublicKey(accountPublicKey)
			balance, err := client.BalanceAt(context.Background(), accountAddress, nil)
			if err != nil {
				log.Panicln(err)
			}
			if balance.Cmp(decimal.Zero.BigInt()) == 0 {
				continue
			}
			nonce, err := client.PendingNonceAt(context.Background(), accountAddress)
			if err != nil {
				log.Panicln(err)
			}
			gasPrice, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Panicln(err)
			}
			totalGasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(21000)))
			if balance.Cmp(totalGasCost) <= 0 {
				continue
			}
			log.Println("balance:", balance, "nonce:", nonce, "gasPrice:", gasPrice)
			// 构造交易
			tx := types.NewTx(&types.LegacyTx{
				Nonce:    nonce,
				To:       &masterAddress,
				Value:    new(big.Int).Sub(balance, totalGasCost),
				Gas:      21000,
				GasPrice: gasPrice,
			})
			// 签名交易
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(networkID), accountPrivateKey)
			if err != nil {
				log.Panicln("Can not sign transaction ", err)
			}
			// 发送交易
			err = client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				log.Panicln(err)
			}
			txHash := signedTx.Hash()
			txHashString := txHash.Hex()
			log.Println("Send transaction success, nonce: ", nonce, " address: ", accountAddress.Hex(), " index:", i, " txHash: ", txHashString)
			time.Sleep(1 * time.Second)
		}
		log.Println("collect finished")
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.Flags().StringP("mnemonic", "m", "", "Set your mnemonic")
	collectCmd.Flags().StringP("rpc", "r", "", "Set your rpc url")
	collectCmd.Flags().UintP("master-index", "", 0, "Set your master account index to receive collected amounts,default is 0")
	collectCmd.Flags().UintP("start-index", "", 1, "Set your start account index,default is 1")
	collectCmd.Flags().UintP("end-index", "", 20, "Set your end account index,must bigger or equal to start-index")
}
