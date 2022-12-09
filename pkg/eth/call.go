package eth

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

type MintData struct {
	Pubkey    string         `json:"pubkey"`
	Receiver  common.Address `json:"receiver"`
	TokenId   int64          `json:"tokenId"`
	Signature string         `json:"signature"`
	Tx        *types.Transaction
}

const RpcUrl = "https://rpc.ankr.com/polygon_mumbai"

const ControllerContract = "0x4Dfbb17704fFAeadEab2d26797989254c956d018"
const MembershipContract = "0x8cb8327774735Eb02b5F50CfC6ac0224048f18Aa"
const PrivateKey = "564df7aedfe0fa9eb710bf6351d0f36ad770b83b46ff5352ffd9dd729e7596fa"

func AddressFromPubkey(pubkey string) (common.Address, error) {
	pk, err := crypto.DecompressPubkey(common.FromHex(pubkey))
	if err != nil {
		return common.BytesToAddress([]byte{}), err
	}
	return crypto.PubkeyToAddress(*pk), nil
}

func SignMint(mintRequest MintData, signerFn func(common.Address, *types.Transaction) (*types.Transaction, error)) (*types.Transaction, error) {
	privKey, verifier := WalletFromPrivate(PrivateKey)
	membershipABI, _ := MembershipMetaData.GetAbi()

	tokenId := big.NewInt(time.Now().Unix())
	packMint, err := membershipABI.Pack("mint", mintRequest.Receiver, tokenId)
	if err != nil {
		return nil, err
	}

	hashPackMint := crypto.Keccak256Hash(packMint)
	//log.Println(common.Bytes2Hex(hashPackMint.Bytes()))

	nonce := big.NewInt(0)
	pack := crypto.Keccak256Hash(EncodePacked(
		hashPackMint.Bytes(), common.HexToAddress(MembershipContract).Bytes(), EncodeUint256Int(nonce.Int64())))
	//log.Println(common.Bytes2Hex(pack.Bytes()))

	hash := SignMessage(pack.Bytes())
	signature := Sign(hash, privKey)
	//log.Println(common.Bytes2Hex(signature))

	rClient, _ := ethclient.Dial(RpcUrl)
	controller, err := NewController(common.HexToAddress(ControllerContract), rClient)
	if err != nil {
		return nil, err
	}

	address, _ := AddressFromPubkey(mintRequest.Pubkey)
	return controller.Execute(&bind.TransactOpts{
		From:   address,
		Signer: signerFn,
	}, IControllerControllerRequest{
		Verifier: *verifier,
		To:       common.HexToAddress(MembershipContract),
		Nonce:    nonce,
		TokenId:  tokenId,
		Data:     packMint,
	}, signature)
}
