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
	"money-distribute-collect/src/utils"
	"strings"
	"time"
)

// distributeCmd represents the distribute command
var distributeCmd = &cobra.Command{
	Use:   "distribute",
	Short: "Distribute money to bip-44 sequence addresses",
	Long:  `Distribute money to bip-44 sequence addresses`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("distribute called")
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

		eachAccountDistributed, err := cmd.Flags().GetString("each-account-distributed")
		if err != nil {
			log.Panicln(errors.New("each-account-distributed is required"))
		}
		eachAccountDistributed = strings.TrimSpace(eachAccountDistributed)
		if eachAccountDistributed == "" {
			log.Panicln(errors.New("each-account-distributed is required"))
		}
		eachAccountDistributedBigDecimal, err := decimal.NewFromString(eachAccountDistributed)
		if err != nil {
			log.Panicln(errors.New("each-account-distributed not a valid human-readable number, eg:0.001"))
		}
		eachAccountDistributedNativeBigDecimal := eachAccountDistributedBigDecimal.Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)))
		if eachAccountDistributedNativeBigDecimal.Equal(decimal.Zero) {
			log.Panicln("each-account-distributed must bigger than 0")
		}

		distributorIndex, err := cmd.Flags().GetUint("distributor-index")
		if err != nil {
			log.Panicln(errors.New("distributor-index is required"))
		}
		startIndex, err := cmd.Flags().GetUint("start-index")
		if err != nil {
			log.Panicln(errors.New("start-index is required"))
		}
		endIndex, err := cmd.Flags().GetUint("end-index")
		if err != nil {
			log.Panicln(errors.New("end-index is required"))
		}

		if startIndex > endIndex {
			log.Panicln(errors.New("start-index must less than or equal to end-index"))
		}

		// 生成distributor账户
		distributorPrivateKey := utils.GetPrivateKey(mnemonic, distributorIndex)
		distributorAddress := utils.GetAddressFromPrivateKey(distributorPrivateKey)
		log.Println("Distributor account address: ", distributorAddress.Hex())

		client, err := ethclient.Dial(rpc)
		if err != nil {
			log.Panicln(err)
		}

		// 获取网络 ID（Chain ID）
		networkID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Panicf("Failed to get network ID: %v\n", err)
		}

		// 获取distributor账户nonce
		distributorNonce, err := client.PendingNonceAt(context.Background(), distributorAddress)
		if err != nil {
			log.Panicln("Can not get nonce ", err)
		}

		// 获取distributor账户余额
		distributorBalance, err := client.BalanceAt(context.Background(), distributorAddress, nil)
		if err != nil {
			log.Panicln("Can not get balance ", err)
		}
		log.Println("Distributor account balance: ", distributorBalance)
		for i := startIndex; i <= endIndex; i++ {
			if i == distributorIndex {
				continue
			}
			// 获取当前gas price
			txGasPrice, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Panicln("Can not get gas price ", err)
			}
			// 获取地址
			targetAddress := utils.GetAddressFromPrivateKey(utils.GetPrivateKey(mnemonic, i))

			// 构造交易
			tx := types.NewTx(&types.LegacyTx{
				Nonce:    distributorNonce,
				To:       &targetAddress,
				Value:    eachAccountDistributedNativeBigDecimal.BigInt(),
				Gas:      21000,
				GasPrice: txGasPrice,
			})
			// 签名交易
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(networkID), distributorPrivateKey)
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
			log.Println("Send transaction success, nonce: ", distributorNonce, " txHash: ", txHashString, " to: ", targetAddress.Hex(), " index: ", i, " amount: ", eachAccountDistributed)
			distributorNonce++
			time.Sleep(5 * time.Second)
		}
		log.Println("Distribute success")
	},
}

func init() {
	rootCmd.AddCommand(distributeCmd)
	distributeCmd.Flags().StringP("mnemonic", "m", "", "Set your mnemonic")
	distributeCmd.Flags().StringP("rpc", "r", "", "Set your rpc url")
	distributeCmd.Flags().StringP("each-account-distributed", "", "", "Set your each account distributed human-readable number, default is 0，eg:0.1")
	distributeCmd.Flags().UintP("distributor-index", "", 0, "Set your distributor account index,default is 0")
	distributeCmd.Flags().UintP("start-index", "", 1, "Set your start account index,default is 1, means only distribute to start account")
	distributeCmd.Flags().UintP("end-index", "", 20, "Set your end account index,must bigger or equal to start-index")
}
