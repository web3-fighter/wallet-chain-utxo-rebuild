package base

import (
	"github.com/btcsuite/btcd/rpcclient"
)

type BtcClient struct {
	*rpcclient.Client
	compressed bool
}

func NewBtcClient(client *rpcclient.Client) *BtcClient {
	return &BtcClient{
		Client:     client,
		compressed: true,
	}
}
