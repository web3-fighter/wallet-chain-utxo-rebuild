package bitcoin

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/log"
	"github.com/web3-fighter/wallet-chain-utxo/domain"
	"github.com/web3-fighter/wallet-chain-utxo/service"
	"github.com/web3-fighter/wallet-chain-utxo/service/base"
	"github.com/web3-fighter/wallet-chain-utxo/service/unimplemente"
)

const ChainName = "Bitcoin"

const (
	p2pkhFormat  = "p2pkh"
	p2wpkhFormat = "p2wpkh"
	p2shFormat   = "p2sh"
	p2trFormat   = "p2tr"
)

var _ service.WalletUtXoService = (*BitcoinNodeService)(nil)

type BitcoinNodeService struct {
	btcClient       *base.BtcClient
	btcDataClient   *base.BaseDataClient
	thirdPartClient *BlockChainClient
	unimplemente.UnimplementedService
}

func (s *BitcoinNodeService) ConvertAddress(ctx context.Context, param domain.ConvertAddressParam) (string, error) {
	var address string
	compressedPubKeyBytes, err := hex.DecodeString(param.PublicKey)
	if err != nil {
		log.Error("decode public key fail", "err", err)
		return address, fmt.Errorf("decode public key fail: %w", err)
	}
	pubKeyHash := btcutil.Hash160(compressedPubKeyBytes)
	switch param.Format {
	case p2pkhFormat:
		p2pkhAddr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
		if err != nil {
			log.Error("create p2pkh address fail", "err", err)
			return address, err
		}
		address = p2pkhAddr.EncodeAddress()
		break
	case p2wpkhFormat:
		witnessAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
		if err != nil {
			log.Error("create p2wpkh fail", "err", err)
		}
		address = witnessAddr.EncodeAddress()
		break
	case p2shFormat:
		witnessAddr, _ := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
		script, err := txscript.PayToAddrScript(witnessAddr)
		if err != nil {
			log.Error("create p2sh address script fail", "err", err)
			return address, err
		}
		p2shAddr, err := btcutil.NewAddressScriptHash(script, &chaincfg.MainNetParams)
		if err != nil {
			log.Error("create p2sh address fail", "err", err)
			return address, err
		}
		address = p2shAddr.EncodeAddress()
		break
	case p2trFormat:
		pubKey, err := btcec.ParsePubKey(compressedPubKeyBytes)
		if err != nil {
			log.Error("parse public key fail", "err", err)
			return address, err
		}
		taprootPubKey := schnorr.SerializePubKey(pubKey)
		taprootAddr, err := btcutil.NewAddressTaproot(taprootPubKey, &chaincfg.MainNetParams)
		if err != nil {
			log.Error("create taproot address fail", "err", err)
			return address, err
		}
		address = taprootAddr.EncodeAddress()
	default:
		return address, errors.New("Do not support address type")
	}
	return address, nil
}

func (s *BitcoinNodeService) ValidAddress(ctx context.Context, param domain.ValidAddressParam) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) GetBlockByNumber(ctx context.Context, param domain.BlockNumberParam) (domain.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) GetBlockByHash(ctx context.Context, param domain.BlockHashParam) (domain.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) GetBlockHeaderByHash(ctx context.Context, param domain.BlockHeaderHashParam) (domain.BlockHeader, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) GetBlockHeaderByNumber(ctx context.Context, param domain.BlockHeaderNumberParam) (domain.BlockHeader, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) ListBlockHeaderByRange(ctx context.Context, param domain.BlockHeaderByRangeParam) ([]domain.BlockHeader, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) GetAccount(ctx context.Context, param domain.AccountParam) (domain.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) GetFee(ctx context.Context, param domain.FeeParam) (domain.Fee, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) SendTx(ctx context.Context, param domain.SendTxParam) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) ListTxByAddress(ctx context.Context, param domain.TxAddressParam) ([]domain.TxMessage, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) GetTxByHash(ctx context.Context, param domain.GetTxByHashParam) (domain.TxMessage, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) CreateUnSignTransaction(ctx context.Context, param domain.UnSignTransactionParam) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) BuildSignedTransaction(ctx context.Context, param domain.SignedTransactionParam) (domain.SignedTransaction, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) DecodeTransaction(ctx context.Context, param domain.DecodeTransactionParam) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) VerifySignedTransaction(ctx context.Context, param domain.VerifyTransactionParam) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) GetExtraData(ctx context.Context, param domain.ExtraDataParam) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewBitcoinNodeService(
	btcClient *base.BtcClient,
	btcDataClient *base.BaseDataClient,
	thirdPartClient *BlockChainClient,
) *BitcoinNodeService {
	return &BitcoinNodeService{
		btcClient:       btcClient,
		btcDataClient:   btcDataClient,
		thirdPartClient: thirdPartClient,
	}
}
