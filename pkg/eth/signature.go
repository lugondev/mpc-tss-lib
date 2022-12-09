package eth

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
)

func SignMessage(data []byte) []byte {
	msg := fmt.Sprintf("\u0019Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func Sign(hash []byte, privateKey *ecdsa.PrivateKey) []byte {
	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	sig[64] += 27

	return sig
}

func VerifySig(from, sigHex string, msg []byte) bool {
	fromAddr := common.HexToAddress(from)
	return fromAddr == RecoverSig(sigHex, msg)
}

func RecoverSig(sigHex string, msg []byte) common.Address {
	sig := hexutil.MustDecode(sigHex)
	if sig[64] != 27 && sig[64] != 28 {
		return common.HexToAddress("0x")
	}
	sig[64] -= 27

	pubKey, err := crypto.SigToPub(SignMessage(msg), sig)
	if err != nil {
		return common.HexToAddress("0x")
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr
}

func EncodePacked(input ...[]byte) []byte {
	return bytes.Join(input, nil)
}

func EncodeBytesString(v string) []byte {
	decoded, err := hex.DecodeString(v)
	if err != nil {
		panic(err)
	}
	return decoded
}

func EncodeUint256String(v string) []byte {
	bn := new(big.Int)
	bn.SetString(v, 10)
	return math.U256Bytes(bn)
}

func EncodeUint256Int(v int64) []byte {
	bn := new(big.Int)
	bn.SetInt64(v)
	return math.U256Bytes(bn)
}

func EncodeUint256Array(arr []int) []byte {
	var res [][]byte
	for _, v := range arr {
		b := EncodeUint256Int(int64(v))
		res = append(res, b)
	}

	return bytes.Join(res, nil)
}

func WalletFromPrivate(privateKeyStr string) (*ecdsa.PrivateKey, *common.Address) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, nil
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return privateKey, &address
}

func GetAddressFromECDSAPrivate(ec *ecdsa.PrivateKey) *common.Address {
	publicKeyECDSA, ok := ec.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &address
}
