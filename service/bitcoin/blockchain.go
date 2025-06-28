package bitcoin

import (
	"errors"
	"fmt"
	gresty "github.com/go-resty/resty/v2"
)

var errBlockChainHTTPError = errors.New("blockchain http error")

type BlockChainClient struct {
	client *gresty.Client
}

func (c *BlockChainClient) GetBalanceByAddress(address string) (string, error) {
	var accountBalance map[string]*AccountBalance
	response, err := c.client.R().
		SetResult(&accountBalance).
		Get("/balance?active=" + address)
	if err != nil {
		return "", fmt.Errorf("cannot get account balance: %w", err)
	}
	if response.StatusCode() != 200 {
		return "", errors.New("get account balance fail")
	}
	return accountBalance[address].FinalBalance.String(), nil
}

func (c *BlockChainClient) GetAccountUtxo(address string) ([]UnspentOutput, error) {
	var utxoUnspentList UnspentOutputList
	response, err := c.client.R().
		SetResult(&utxoUnspentList).
		Get("/unspent?active=" + address)
	if err != nil {
		return nil, fmt.Errorf("cannot utxo fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	return utxoUnspentList.UnspentOutputs, nil
}

func (c *BlockChainClient) GetTransactionsByAddress(address, pageSize, page string) (*Transaction, error) {
	var transactionList Transaction
	response, err := c.client.R().
		SetResult(&transactionList).
		Get("/rawaddr/" + address + "?limit=" + pageSize + "&offset=" + page)
	if err != nil {
		return nil, fmt.Errorf("cannot utxo fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	return &transactionList, nil
}

func (c *BlockChainClient) GetTransactionsByHash(txHash string) (*TxsItem, error) {
	var transaction TxsItem
	response, err := c.client.R().
		SetResult(&transaction).
		Get("/rawtx/" + txHash)
	if err != nil {
		return nil, fmt.Errorf("cannot utxo fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	return &transaction, nil
}

func NewBlockChainClient(url string) (*BlockChainClient, error) {
	// validate if blockchain url is provided or not
	if url == "" {
		return nil, fmt.Errorf("blockchain URL cannot be empty")
	}

	client := gresty.New()
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errBlockChainHTTPError)
		}
		return nil
	})
	return &BlockChainClient{
		client: client,
	}, nil
}
