package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

const RpcUrl = "https://rpc.ankr.com/polygon_mumbai"

const ControllerContract = "0x4Dfbb17704fFAeadEab2d26797989254c956d018"
const MembershipContract = "0x8cb8327774735Eb02b5F50CfC6ac0224048f18Aa"
const PrivateKey = "564df7aedfe0fa9eb710bf6351d0f36ad770b83b46ff5352ffd9dd729e7596fa"

func TestChainInfo(ctx context.Context) {
	rClient, _ := ethclient.Dial(RpcUrl)
	chainId, err := rClient.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("chain id:", chainId.Int64())
}

func TestSignMint(ctx context.Context) {
	privKey, address := WalletFromPrivate(PrivateKey)
	log.Println(address.String())
	membershipABI, _ := MembershipMetaData.GetAbi()

	tokenId := big.NewInt(time.Now().Unix())
	packMint, err := membershipABI.Pack("mint", address, tokenId)
	if err != nil {
		log.Fatal(err)
	}

	hashPackMint := crypto.Keccak256Hash(packMint)
	log.Println(common.Bytes2Hex(hashPackMint.Bytes()))

	nonce := big.NewInt(0)
	pack := crypto.Keccak256Hash(EncodePacked(
		hashPackMint.Bytes(), common.HexToAddress(MembershipContract).Bytes(), EncodeUint256Int(nonce.Int64())))
	log.Println(common.Bytes2Hex(pack.Bytes()))

	hash := SignMessage(pack.Bytes())
	signature := Sign(hash, privKey)
	log.Println(common.Bytes2Hex(signature))

	rClient, _ := ethclient.Dial(RpcUrl)
	controller, err := NewController(common.HexToAddress(ControllerContract), rClient)
	if err != nil {
		log.Fatal(err)
	}
	chainId, _ := rClient.ChainID(ctx)
	auth, _ := bind.NewKeyedTransactorWithChainID(privKey, chainId)
	transaction, err := controller.Execute(&bind.TransactOpts{
		From:   *address,
		Signer: auth.Signer,
	}, IControllerControllerRequest{
		Verifier: *address,
		To:       common.HexToAddress(MembershipContract),
		Nonce:    nonce,
		TokenId:  tokenId,
		Data:     packMint,
	}, signature)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tx: " + transaction.Hash().String())
}
