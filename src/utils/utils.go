package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func GetPrivateKey(mnemonic string, accountIndex uint) *ecdsa.PrivateKey {
	seed := bip39.NewSeed(mnemonic, "") // 可以提供密码短语
	masterKey, _ := bip32.NewMasterKey(seed)
	// 44'
	purposeKey, _ := masterKey.NewChildKey(0x8000002C)
	// 60'
	coinTypeKey, _ := purposeKey.NewChildKey(0x8000003C)
	// 0'
	accountKey, _ := coinTypeKey.NewChildKey(0x80000000)
	// 0
	changeKey, _ := accountKey.NewChildKey(0)
	// addressIndex
	addressKey, _ := changeKey.NewChildKey(uint32(accountIndex))
	privateKey, _ := crypto.ToECDSA(addressKey.Key)
	return privateKey
}

func GetPublicKey(privateKey *ecdsa.PrivateKey) *ecdsa.PublicKey {
	return &privateKey.PublicKey
}

func GetAddressFromPublicKey(publicKey *ecdsa.PublicKey) common.Address {
	return crypto.PubkeyToAddress(*publicKey)
}

func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(*GetPublicKey(privateKey))
}

func GenMnemonic() string {
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

func PrivateKey2String(privateKey *ecdsa.PrivateKey) string {
	privateKeyBytes := privateKey.D.Bytes()              // 将 *big.Int 转换为字节数组
	privateKeyHex := hex.EncodeToString(privateKeyBytes) // 将字节数组转换为十六进制字符串
	return "0x" + privateKeyHex
}

func PublicKey2String(publicKey *ecdsa.PublicKey) string {
	// 获取椭圆曲线参数
	curve := publicKey.Curve.Params()

	// 获取公钥的X和Y坐标
	x := publicKey.X
	y := publicKey.Y

	// 分配一个足够长度的字节切片
	compressed := make([]byte, (curve.BitSize+7)>>3+1)

	// 填充X坐标
	xBytes := x.Bytes()
	copy(compressed[1:], xBytes)

	// 设置前缀
	if y.Bit(0) == 0 { // Y坐标是偶数
		compressed[0] = 0x02
	} else { // Y坐标是奇数
		compressed[0] = 0x03
	}

	publicKeyHex := hex.EncodeToString(compressed)
	return "0x" + publicKeyHex
}
