package bitcoin

import (
	"math/big"
)

type TxStatus int32

func (x TxStatus) ToInt32() int32 {
	return int32(x)
}

const (
	TxStatus_NotFound              TxStatus = 0
	TxStatus_Pending               TxStatus = 1
	TxStatus_Failed                TxStatus = 2
	TxStatus_Success               TxStatus = 3
	TxStatus_ContractExecuteFailed TxStatus = 4
	TxStatus_Other                 TxStatus = 5
)

// Enum value maps for TxStatus.
var (
	TxStatus_name = map[int32]string{
		0: "NotFound",
		1: "Pending",
		2: "Failed",
		3: "Success",
		4: "ContractExecuteFailed",
		5: "Other",
	}
	TxStatus_value = map[string]int32{
		"NotFound":              0,
		"Pending":               1,
		"Failed":                2,
		"Success":               3,
		"ContractExecuteFailed": 4,
		"Other":                 5,
	}
)

type SpendingOutpointsItem struct {
	N       uint64  `json:"n"`
	TxIndex big.Int `json:"tx_index"`
}

type PrevOut struct {
	Addr              string                  `json:"addr"`
	N                 uint64                  `json:"n"`
	Script            string                  `json:"script"`
	SpendingOutpoints []SpendingOutpointsItem `json:"spending_outpoints"`
	Spent             bool                    `json:"spent"`
	TxIndex           big.Int                 `json:"tx_index"`
	Type              uint64                  `json:"type"`
	Value             big.Int                 `json:"value"`
}

type InputItem struct {
	Sequence big.Int `json:"sequence"`
	Witness  string  `json:"witness"`
	Script   string  `json:"script"`
	Index    uint64  `json:"index"`
	PrevOut  PrevOut `json:"prev_out"`
}

type OutItem struct {
	Type              uint64                  `json:"type"`
	Spent             bool                    `json:"spent"`
	Value             big.Int                 `json:"value"`
	SpendingOutpoints []SpendingOutpointsItem `json:"spending_outpoints"`
	N                 uint64                  `json:"n"`
	TxIndex           big.Int                 `json:"tx_index"`
	Script            string                  `json:"script"`
	Addr              string                  `json:"addr"`
}

type TxsItem struct {
	Hash        string      `json:"hash"`
	Ver         uint64      `json:"ver"`
	VinSz       uint64      `json:"vin_sz"`
	VoutSz      uint64      `json:"vout_sz"`
	Size        uint64      `json:"size"`
	Weight      uint64      `json:"weight"`
	Fee         big.Int     `json:"fee"`
	RelayedBy   string      `json:"relayed_by"`
	LockTime    big.Int     `json:"lock_time"`
	TxIndex     uint64      `json:"tx_index"`
	DoubleSpend bool        `json:"double_spend"`
	Time        big.Int     `json:"time"`
	BlockIndex  big.Int     `json:"block_index"`
	BlockHeight big.Int     `json:"block_height"`
	Inputs      []InputItem `json:"inputs"`
	Out         []OutItem   `json:"out"`
	Result      big.Int     `json:"result"`
	Balance     big.Int     `json:"balance"`
}

type Transaction struct {
	Hash160       string    `json:"hash160"`
	Address       string    `json:"address"`
	NTx           uint64    `json:"n_tx"`
	NUnredeemed   big.Int   `json:"n_unredeemed"`
	TotalReceived big.Int   `json:"total_received"`
	TotalSent     big.Int   `json:"total_sent"`
	FinalBalance  big.Int   `json:"final_balance"`
	Txs           []TxsItem `json:"txs"`
}

type BlockData struct {
	Hash              string   `json:"hash"`
	Confirmations     uint64   `json:"confirmations"`
	Size              uint64   `json:"size"`
	StrippedSize      uint64   `json:"strippedsize"`
	Weight            uint64   `json:"weight"`
	Height            uint64   `json:"height"`
	Version           uint64   `json:"version"`
	VersionHex        string   `json:"version_hex"`
	Merkleroot        string   `json:"merkleroot"`
	Tx                []string `json:"tx"`
	Time              uint64   `json:"time"`
	MedianTime        uint64   `json:"mediantime"`
	Nonce             uint64   `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        uint64   `json:"difficulty"`
	ChainWork         string   `json:"chainwork"`
	NTx               uint64   `json:"nTx"`
	PreviousBlockHash string   `json:"previousblockhash"`
	NextBlockHash     string   `json:"nextblockhash"`
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type Vin struct {
	Coinbase    string    `json:"coinbase"`
	TxId        string    `json:"txid"`
	Vout        uint32    `json:"vout"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	Sequence    uint64    `json:"sequence"`
	TxInWitness []string  `json:"txinwitness"`
}

type ScriptPubKey struct {
	//Asm     string `json:"asm"`
	//Hex     string `json:"hex"`
	//Desc    string `json:"desc"`
	//Address string `json:"addresses"`
	//Type    string `json:"type"`

	Asm       string   `json:"asm"`
	Hex       string   `json:"hex,omitempty"`
	ReqSigs   int32    `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
}

type Vout struct {
	Value        float64      `json:"value"`
	N            uint32       `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptpubkey"`
}

type RawTransactionData struct {
	TxId          string `json:"txid"`
	Hash          string `json:"hash"`
	Version       uint64 `json:"version"`
	Size          uint64 `json:"size"`
	VSize         uint64 `json:"vsize"`
	Weight        uint64 `json:"weight"`
	LockTime      uint64 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	Hex           string `json:"hex"`
	Blockhash     string `json:"blockhash"`
	Confirmations uint64 `json:"confirmations"`
	BlockTime     uint64 `json:"blocktime"`
	Time          uint64 `json:"time"`
}

type VinItem struct {
	Hash    string `json:"hash,omitempty"`
	Index   uint32 `json:"index,omitempty"`
	Amount  int64  `json:"amount,omitempty"`
	Address string `json:"address,omitempty"`
}

type VoutItem struct {
	Address string `json:"address,omitempty"`
	Amount  int64  `json:"amount,omitempty"`
	Index   uint32 `json:"index,omitempty"`
}

type UtxoTransaction struct {
	TxHash      string      `protobuf:"bytes,3,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	Status      TxStatus    `protobuf:"varint,4,opt,name=status,proto3,enum=savour_rpc.wallet.TxStatus" json:"status,omitempty"`
	Vins        []*VinItem  `protobuf:"bytes,5,rep,name=vins,proto3" json:"vins,omitempty"`
	Vouts       []*VoutItem `protobuf:"bytes,6,rep,name=vouts,proto3" json:"vouts,omitempty"`
	SignHashes  [][]byte    `protobuf:"bytes,7,rep,name=sign_hashes,json=signHashes,proto3" json:"sign_hashes,omitempty"`
	CostFee     string      `protobuf:"bytes,8,opt,name=cost_fee,json=costFee,proto3" json:"cost_fee,omitempty"`
	BlockHeight uint64      `protobuf:"varint,9,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	BlockTime   uint64      `protobuf:"varint,10,opt,name=block_time,json=blockTime,proto3" json:"block_time,omitempty"`
}

type AccountBalance struct {
	FinalBalance  big.Int `json:"final_balance"`
	NTx           big.Int `json:"n_tx"`
	TotalReceived big.Int `json:"total_received"`
}

type UnspentOutput struct {
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxHash          string `json:"tx_hash"`
	TxOutputN       uint64 `json:"tx_output_n"`
	Script          string `json:"script"`
	Value           uint64 `json:"value"`
	ValueHex        string `json:"value_hex"`
	Confirmations   uint64 `json:"confirmations"`
	TxIndex         uint64 `json:"tx_index"`
	BlockTime       string `json:"block_time,omitempty"`
}

type UnspentOutputList struct {
	Notice         string          `json:"notice"`
	UnspentOutputs []UnspentOutput `json:"unspent_outputs"`
}

type GasFee struct {
	ChainFullName       string `json:"chainFullName"`
	ChainShortName      string `json:"chainShortName"`
	Symbol              string `json:"symbol"`
	BestTransactionFee  string `json:"bestTransactionFee"`
	RecommendedGasPrice string `json:"recommendedGasPrice"`
	RapidGasPrice       string `json:"rapidGasPrice"`
	StandardGasPrice    string `json:"standardGasPrice"`
	SlowGasPrice        string `json:"slowGasPrice"`
}

type GasFeeData struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []GasFee `json:"data"`
}
