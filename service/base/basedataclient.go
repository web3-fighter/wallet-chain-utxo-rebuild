package base

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/web3-fighter/chain-explorer-api/client"
	"github.com/web3-fighter/chain-explorer-api/client/oklink"
	"github.com/web3-fighter/chain-explorer-api/types"
)

type BaseDataClient struct {
	chainShortName string
	explorerName   string
	baseDataCli    client.ChainExplorer
}

func (c *BaseDataClient) GetFee(chainShortName string) (*types.GasEstimateFeeResponse, error) {
	gefr := &types.GasEstimateFeeRequest{
		ChainShortName: chainShortName,
		ExplorerName:   oklink.ChainExplorerName,
	}
	okResp, err := c.baseDataCli.GetEstimateGasFee(gefr)
	if err != nil {
		log.Error("get estimate gas fee fail", "err", err)
		return nil, err
	}
	return okResp, nil
}

func (c *BaseDataClient) GetAccountBalance(chainShortName, address string) (*types.AccountBalanceResponse, error) {
	accountItem := []string{address}
	contractAddress := []string{"0x00"}
	acbr := &types.AccountBalanceRequest{
		ChainShortName:  c.explorerName,
		ExplorerName:    chainShortName,
		Account:         accountItem,
		ContractAddress: contractAddress,
	}
	balanceResponse, err := c.baseDataCli.GetAccountBalance(acbr)
	if err != nil {
		log.Error("get balance response fail", "err", err)
		return nil, err
	}
	return balanceResponse, nil
}

func (c *BaseDataClient) GetAccountUtXoList(chainShortName, address string) ([]*types.AccountUtxoResponse, error) {
	utxoRequest := &types.AccountUtxoRequest{
		ChainShortName: chainShortName,
		ExplorerName:   c.explorerName,
		Address:        address,
	}
	utxoResponse, err := c.baseDataCli.GetAccountUtxo(utxoRequest)
	if err != nil {
		log.Error("get account utxo fail", "err", err)
		return nil, err
	}
	return utxoResponse, nil
}

func (c *BaseDataClient) GetTxListByAddress(chainShortName, address string, page, pageSize uint64) (*types.TransactionResponse[types.AccountTxResponse], error) {
	txRequest := &types.AccountTxRequest{
		ChainShortName: chainShortName,
		ExplorerName:   c.explorerName,
		Action:         types.OkLinkActionUtxo,
		Address:        address,
		PageRequest: types.PageRequest{
			Page:  page,
			Limit: pageSize,
		},
	}
	txListResponse, err := c.baseDataCli.GetTxByAddress(txRequest)
	if err != nil {
		log.Error("get tx by address fail", "err", err)
		return nil, err
	}
	log.Info("tx list response success", "transactionList Length", len(txListResponse.TransactionList))
	return txListResponse, nil
}

func (c *BaseDataClient) GetTxByHash(chainShortName, txId string) (*types.TxResponse, error) {
	txRequest := &types.TxRequest{
		ChainShortName: chainShortName,
		ExplorerName:   c.explorerName,
		Txid:           txId,
	}
	txResponse, err := c.baseDataCli.GetTxByHash(txRequest)
	if err != nil {
		log.Error("get tx by address fail", "err", err)
		return nil, err
	}
	return txResponse, nil
}

func NewBaseDataClient(explorerName string, baseDataCli client.ChainExplorer) *BaseDataClient {
	return &BaseDataClient{
		explorerName: explorerName,
		baseDataCli:  baseDataCli,
	}
}
