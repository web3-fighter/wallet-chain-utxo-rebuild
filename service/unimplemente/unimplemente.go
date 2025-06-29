package unimplemente

import (
	"context"
	"github.com/web3-fighter/wallet-chain-utxo/domain"
	"github.com/web3-fighter/wallet-chain-utxo/service"
)

var _ service.WalletUtXoService = (*UnimplementedService)(nil)

type UnimplementedService struct{}

func (s *UnimplementedService) GetBalanceByAddress(ctx context.Context, param domain.BalanceByAddressParam) (domain.Balance, error) {
	return domain.Balance{}, nil
}

func (s *UnimplementedService) GetUnspentOutputs(ctx context.Context, req domain.UnspentOutputsParam) ([]domain.UnspentOutput, error) {
	return nil, nil
}

func (s *UnimplementedService) GetBlockHeaderByNumber(ctx context.Context, param domain.BlockHeaderNumberParam) (domain.BlockHeader, error) {
	return domain.BlockHeader{}, nil
}

func (s *UnimplementedService) GetSupportChains(ctx context.Context, param domain.SupportChainsParam) (bool, error) {
	return true, nil
}

func (s *UnimplementedService) ConvertAddress(ctx context.Context, param domain.ConvertAddressParam) (string, error) {
	return "", nil
}

func (s *UnimplementedService) ValidAddress(ctx context.Context, param domain.ValidAddressParam) (bool, error) {
	return true, nil
}

func (s *UnimplementedService) GetBlockByNumber(ctx context.Context, param domain.BlockNumberParam) (domain.Block, error) {
	return domain.Block{}, nil
}

func (s *UnimplementedService) GetBlockByHash(ctx context.Context, param domain.BlockHashParam) (domain.Block, error) {
	return domain.Block{}, nil
}

func (s *UnimplementedService) GetBlockHeaderByHash(ctx context.Context, param domain.BlockHeaderHashParam) (domain.BlockHeader, error) {
	return domain.BlockHeader{}, nil
}

func (s *UnimplementedService) GetFee(ctx context.Context, param domain.FeeParam) (domain.Fee, error) {
	return domain.Fee{}, nil
}

func (s *UnimplementedService) SendTx(ctx context.Context, param domain.SendTxParam) (string, error) {
	return "", nil
}

func (s *UnimplementedService) ListTxByAddress(ctx context.Context, param domain.TxAddressParam) ([]domain.TxMessage, error) {
	return nil, nil
}

func (s *UnimplementedService) GetTxByHash(ctx context.Context, param domain.GetTxByHashParam) (domain.TxMessage, error) {
	return domain.TxMessage{}, nil
}

func (s *UnimplementedService) CreateUnSignTransaction(ctx context.Context, param domain.UnSignTransactionParam) (domain.UnSignTransactionResult, error) {
	return domain.UnSignTransactionResult{}, nil
}

func (s *UnimplementedService) BuildSignedTransaction(ctx context.Context, param domain.SignedTransactionParam) (domain.SignedTransaction, error) {
	return domain.SignedTransaction{}, nil
}

func (s *UnimplementedService) DecodeTransaction(ctx context.Context, param domain.DecodeTransactionParam) (string, error) {
	return "", nil
}

func (s *UnimplementedService) VerifySignedTransaction(ctx context.Context, param domain.VerifyTransactionParam) (bool, error) {
	return true, nil
}

func (s *UnimplementedService) GetExtraData(ctx context.Context, param domain.ExtraDataParam) (string, error) {
	return "", nil
}
