package bitcoin

import (
	"context"
	"github.com/web3-fighter/wallet-chain-utxo/domain"
	"github.com/web3-fighter/wallet-chain-utxo/service"
	"github.com/web3-fighter/wallet-chain-utxo/service/base"
)

const ChainName = "Bitcoin"

var _ service.WalletUtXoService = (*BitcoinNodeService)(nil)

type BitcoinNodeService struct {
	btcClient       *base.BtcClient
	btcDataClient   *base.BaseDataClient
	thirdPartClient *BlockChainClient
}

func (s *BitcoinNodeService) GetSupportChains(ctx context.Context, param domain.SupportChainsParam) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *BitcoinNodeService) ConvertAddress(ctx context.Context, param domain.ConvertAddressParam) (string, error) {
	//TODO implement me
	panic("implement me")
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
