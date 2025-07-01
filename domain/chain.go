package domain

import (
	"math/big"
)

type DecodeTx struct {
	Hash       string
	SignHashes [][]byte
	Vins       []Vin
	Vouts      []Vout
	CostFee    *big.Int
}

type SignedTransactionParam struct {
	ConsumerToken string   `json:"consumer_token,omitempty"`
	Chain         string   `json:"chain,omitempty"`
	Network       string   `json:"network,omitempty"`
	TxData        []byte   `json:"tx_data,omitempty"`
	Signatures    [][]byte `json:"signatures,omitempty"`
	PublicKeys    [][]byte `json:"public_keys,omitempty"`
}

type SignedTransaction struct {
	SignedTxData []byte `json:"signed_tx_data,omitempty"`
	Hash         []byte `json:"hash,omitempty"`
}

type UnSignTransactionParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Network       string `json:"network,omitempty"`
	Fee           string `json:"fee,omitempty"`   // 可选，最终会体现在 output 金额里"`
	Vins          []Vin  `json:"vins,omitempty"`  // 输入（来源地址、引用 txid 和 vout）
	Vouts         []Vout `json:"vouts,omitempty"` // 输出（目标地址、金额）
}

type UnSignTransactionResult struct {
	TxData     []byte   `json:"tx_data,omitempty"`
	SignHashes [][]byte `json:"sign_hashes,omitempty"`
}

type TxByHashParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Coin          string `json:"coin,omitempty"`
	Network       string `json:"network,omitempty"`
	Hash          string `json:"hash,omitempty"`
}

type TxAddressParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Coin          string `json:"coin,omitempty"`
	Network       string `json:"network,omitempty"`
	Address       string `json:"address,omitempty"`
	Brc20Address  string `json:"brc20_address,omitempty"`
	Page          uint32 `json:"page,omitempty"`
	Pagesize      uint32 `json:"pagesize,omitempty"`
	Cursor        string `json:"cursor,omitempty"`
}

type TxMessage struct {
	Hash         string   `json:"hash,omitempty"`          // 交易哈希
	Index        uint32   `json:"index,omitempty"`         // 可选，交易在区块中的顺序
	Froms        []string `json:"froms,omitempty"`         // 所有输入地址
	Tos          []string `json:"tos,omitempty"`           // 所有输出地址
	Values       []Value  `json:"values,omitempty"`        // 输出金额（字符串格式，单位可能为 satoshi）
	Fee          string   `json:"fee,omitempty"`           // 手续费（字符串格式）
	Status       TxStatus `json:"status,omitempty"`        // 状态（通常为成功）
	Type         int32    `json:"type,omitempty"`          // 方向：0 = 转出，1 = 转入
	Height       string   `json:"height,omitempty"`        // 所在区块高度
	Brc20Address string   `json:"brc20_address,omitempty"` // 可选，BRC-20 标准相关字段
	Datetime     string   `json:"datetime,omitempty"`      // 区块时间（交易时间）
}

type Value struct {
	Address string `json:"address,omitempty"`
	Value   string `json:"value,omitempty"`
	//Type string  `json:"type,omitempty"`
}

type SendTxParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Coin          string `json:"coin,omitempty"`
	Network       string `json:"network,omitempty"`
	RawTx         string `json:"raw_tx,omitempty"`
}

type Block struct {
	Msg    string              `json:"msg,omitempty"`
	Height uint64              `json:"height,omitempty"`
	Hash   string              `json:"hash,omitempty"`
	TxList []*BlockTransaction `json:"tx_list,omitempty"`
}

type Vin struct {
	Hash    string `json:"hash,omitempty"`
	Index   uint32 `json:"index,omitempty"`
	Amount  int64  `json:"amount,omitempty"`
	Address string `json:"address,omitempty"`
}

type Vout struct {
	Address string `json:"address,omitempty"`
	Amount  int64  `json:"amount,omitempty"`
	Index   uint32 `json:"index,omitempty"`
}

type BlockTransaction struct {
	Hash          string  `son:"hash,omitempty"`
	Fee           string  `json:"fee,omitempty"`
	Version       uint64  `json:"version"`
	Size          uint64  `json:"size"`
	VSize         uint64  `json:"vsize"`
	Weight        uint64  `json:"weight"`
	LockTime      uint64  `json:"locktime"`
	Hex           string  `json:"hex"`
	Blockhash     string  `json:"blockhash"`
	Confirmations uint64  `json:"confirmations"`
	BlockTime     uint64  `json:"blocktime"`
	Time          uint64  `json:"time"`
	BlockHeight   uint64  `json:"block_height"`
	Status        string  `json:"status,omitempty"`
	Vin           []*Vin  `json:"vin,omitempty"`
	Vout          []*Vout `json:"vout,omitempty"`
}

type UnspentOutputsParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Network       string `json:"network,omitempty"`
	Address       string `json:"address,omitempty"`
}

type UnspentOutput struct {
	TxId            string `json:"tx_id,omitempty"`              // 原始交易哈希（正常顺序）
	TxHashBigEndian string `json:"tx_hash_big_endian,omitempty"` // big-endian 格式的交易哈希（某些系统使用）
	TxOutputN       uint64 `json:"tx_output_n,omitempty"`        // 这个 UTXO 在交易中是第几个输出
	Script          string `json:"script,omitempty"`             // 脚本字段（scriptPubKey），后续签名需要
	Height          string `json:"height,omitempty"`             // 区块高度（有可能没赋值）
	BlockTime       string `json:"block_time,omitempty"`         // 区块时间（可选）
	Address         string `json:"address,omitempty"`            // 属于哪个地址
	UnspentAmount   string `json:"unspent_amount,omitempty"`     // 金额（单位：satoshi，string）
	ValueHex        string `json:"value_hex,omitempty"`          // 金额（十六进制字符串，可用于某些原始格式）
	Confirmations   uint64 `json:"confirmations,omitempty"`      // 确认数
	Index           uint64 `json:"index,omitempty"`              // TxIndex：交易在区块中的顺序索引（非 vout）
}

type BlockHashParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Hash          string `json:"hash,omitempty"`
	ViewTx        bool   `json:"view_tx,omitempty"`
}

type BlockNumberParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Height        int64  `json:"height,omitempty"`
	ViewTx        bool   `json:"view_tx,omitempty"`
}

type BlockHeaderHashParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Network       string `json:"network,omitempty"`
	Hash          string `json:"hash,omitempty"`
}

type BlockHeader struct {
	ParentHash string `json:"parent_hash,omitempty"`
	BlockHash  string `json:"block_hash,omitempty"`
	MerkleRoot string `json:"merkle_root,omitempty"`
	Number     string `json:"number,omitempty"`
}

type BlockHeaderNumberParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Network       string `json:"network,omitempty"`
	Height        int64  `json:"height,omitempty"`
}

type ValidAddressParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Network       string `json:"network,omitempty"`
	Address       string `json:"address,omitempty"`
}

type ConvertAddressParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Network       string `json:"network,omitempty"`
	Format        string `json:"format,omitempty"`
	PublicKey     string `json:"public_key,omitempty"`
}

type BalanceByAddressParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Network       string `json:"network,omitempty"`
	Address       string `json:"address,omitempty"`
	Brc20Address  string `json:"brc20_address,omitempty"`
}

type Balance struct {
	Network string `json:"network,omitempty"`
	Balance string `json:"balance,omitempty"`
}

type FeeParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Coin          string `json:"coin,omitempty"`
	Network       string `json:"network,omitempty"`
	RawTx         string `json:"rawTx,omitempty"`
}

/*
Fee
字段	含义

	BestFee	综合推荐费率（有可能是 "NormalFee" + 系统调控）
	BestFeeSat	推荐费用的 satoshi 表示（可能用于预估费用）
	SlowFee	慢速确认费率（较便宜，但可能等几十分钟）
	NormalFee	平均确认时间（一般 1-3 blocks）
	FastFee	快速确认（通常1个区块内）
*/
type Fee struct {
	BestFee    string `json:"best_fee,omitempty"`
	BestFeeSat string `json:"best_fee_sat,omitempty"`
	SlowFee    string `json:"slow_fee,omitempty"`
	NormalFee  string `json:"normal_fee,omitempty"`
	FastFee    string `json:"fast_fee,omitempty"`
}

type DecodeTransactionParam struct {
	Chain   string `json:"chain,omitempty"`
	Network string `json:"network,omitempty"`
	RawData []byte `json:"raw_data,omitempty"`
	Vins    []Vin  `json:"vins,omitempty"`
}

type DecodedTransaction struct {
	TxHash      string   `json:"tx_hash,omitempty"`
	Status      TxStatus `json:"status,omitempty"`
	Vins        []Vin    `json:"vins,omitempty"`
	Vouts       []Vout   `json:"vouts,omitempty"`
	SignHashes  [][]byte `json:"sign_hashes,omitempty"`
	CostFee     string   `json:"cost_fee,omitempty"`
	BlockHeight uint64   `json:"block_height,omitempty"`
	BlockTime   uint64   `json:"block_time,omitempty"`
}

type VerifyTransactionParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Network       string `json:"network,omitempty"`
	RawData       []byte `json:"raw_data,omitempty"`
	Vins          []Vin  `json:"vins,omitempty"`
}

//  TODO-----------------------------

type CommonParam struct {
	ConsumerToken string `protobuf:"bytes,1,opt,name=consumer_token,json=consumerToken,proto3" json:"consumer_token,omitempty"`
	Chain         string `protobuf:"bytes,2,opt,name=chain,proto3" json:"chain,omitempty"`
	Coin          string `protobuf:"bytes,3,opt,name=coin,proto3" json:"coin,omitempty"`
	Network       string `protobuf:"bytes,4,opt,name=network,proto3" json:"network,omitempty"`
}

type ExtraDataParam struct {
	ConsumerToken string `protobuf:"bytes,1,opt,name=consumer_token,json=consumerToken,proto3" json:"consumer_token,omitempty"`
	Chain         string `protobuf:"bytes,2,opt,name=chain,proto3" json:"chain,omitempty"`
	Network       string `protobuf:"bytes,3,opt,name=network,proto3" json:"network,omitempty"`
	Address       string `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Coin          string `protobuf:"bytes,5,opt,name=coin,proto3" json:"coin,omitempty"`
}

type BlockHeaderByRangeParam struct {
	ConsumerToken string `protobuf:"bytes,1,opt,name=consumer_token,json=consumerToken,proto3" json:"consumer_token,omitempty"`
	Chain         string `protobuf:"bytes,2,opt,name=chain,proto3" json:"chain,omitempty"`
	Network       string `protobuf:"bytes,3,opt,name=network,proto3" json:"network,omitempty"`
	Start         string `protobuf:"bytes,4,opt,name=start,proto3" json:"start,omitempty"`
	End           string `protobuf:"bytes,5,opt,name=end,proto3" json:"end,omitempty"`
}

type SupportChainsParam struct {
	ConsumerToken string `json:"consumer_token,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Network       string `json:"network,omitempty"`
}

type TxStatus int32

const (
	TxStatus_NotFound              TxStatus = iota
	TxStatus_Pending               TxStatus = 1
	TxStatus_Failed                TxStatus = 2
	TxStatus_Success               TxStatus = 3
	TxStatus_ContractExecuteFailed TxStatus = 4
	TxStatus_Other                 TxStatus = 5
)

const (
	ChainIdRedisKey     = "chainId:%s"
	ChainId             = "chainId"
	ChainName           = "chainName"
	Polygon             = "Polygon"
	GormInfoFmt         = "%s\n[%.3fms] [rows:%v] %s"
	ZeroAddress         = "0x0000000000000000000000000000000000000000"
	WEthAddress         = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	EthAddress          = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
	SepoliaWETH         = ""
	LogTimeFormat       = "2006-01-02"
	LayerTypeOne        = 1
	LayerTypeTwo        = 2
	StakingTypeStake    = 1
	StakingTypeReward   = 2
	StakingTypeWithdraw = 3

	BridgeOperaInitType     = 1
	BridgeOperaFinalizeType = 2

	ScrollChainId          uint64 = 534352
	PolygonChainId         uint64 = 1101
	PolygonSepoliaChainId  uint64 = 1442
	EthereumChainId        uint64 = 1
	EthereumSepoliaChainId uint64 = 11155111
	BaseChainId            uint64 = 8453
	BaseSepoliaChainId     uint64 = 84532
	MantaChainId           uint64 = 169
	MantaSepoliaChainId    uint64 = 3441006
	MantleSepoliaChainId   uint64 = 5003
	MantleChainId          uint64 = 5000
	ZkFairSepoliaChainId   uint64 = 43851
	ZkFairChainId          uint64 = 42766
	OkxSepoliaChainId      uint64 = 195
	OkxChainId             uint64 = 66
	OpChinId               uint64 = 10
	OpTestChinId           uint64 = 11155420
	LineaChainId           uint64 = 59144
	BlocksLimit                   = 10000
)
