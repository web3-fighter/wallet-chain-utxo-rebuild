package service

import (
	"context"
	"github.com/web3-fighter/wallet-chain-utxo/domain"
)

type WalletUtXoService interface {
	GetSupportChains(ctx context.Context, param domain.SupportChainsParam) (bool, error)
	ConvertAddress(ctx context.Context, param domain.ConvertAddressParam) (string, error)
	ValidAddress(ctx context.Context, param domain.ValidAddressParam) (bool, error)
	GetBlockByNumber(ctx context.Context, param domain.BlockNumberParam) (domain.Block, error)
	GetBlockByHash(ctx context.Context, param domain.BlockHashParam) (domain.Block, error)
	GetBlockHeaderByHash(ctx context.Context, param domain.BlockHeaderHashParam) (domain.BlockHeader, error)
	GetBlockHeaderByNumber(ctx context.Context, param domain.BlockHeaderNumberParam) (domain.BlockHeader, error)
	GetUnspentOutputs(ctx context.Context, req domain.UnspentOutputsParam) ([]domain.UnspentOutput, error)
	GetBalanceByAddress(ctx context.Context, param domain.BalanceByAddressParam) (domain.Balance, error)
	GetFee(ctx context.Context, param domain.FeeParam) (domain.Fee, error)
	SendTx(ctx context.Context, param domain.SendTxParam) (string, error)
	ListTxByAddress(ctx context.Context, param domain.TxAddressParam) ([]domain.TxMessage, error)
	// ----------------
	GetTxByHash(ctx context.Context, param domain.GetTxByHashParam) (domain.TxMessage, error)
	CreateUnSignTransaction(ctx context.Context, param domain.UnSignTransactionParam) (domain.UnSignTransactionResult, error)
	BuildSignedTransaction(ctx context.Context, param domain.SignedTransactionParam) (domain.SignedTransaction, error)
	DecodeTransaction(ctx context.Context, param domain.DecodeTransactionParam) (string, error)
	VerifySignedTransaction(ctx context.Context, param domain.VerifyTransactionParam) (bool, error)
}
