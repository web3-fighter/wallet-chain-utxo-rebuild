package bitcoin

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ethereum/go-ethereum/log"
	"github.com/shopspring/decimal"
	"github.com/web3-fighter/wallet-chain-utxo/domain"
	"github.com/web3-fighter/wallet-chain-utxo/service"
	"github.com/web3-fighter/wallet-chain-utxo/service/base"
	"github.com/web3-fighter/wallet-chain-utxo/service/unimplemente"
	"math"
	"math/big"
	"strconv"
	"strings"
)

const (
	p2pkhFormat  = "p2pkh"
	p2wpkhFormat = "p2wpkh"
	p2shFormat   = "p2sh"
	p2trFormat   = "p2tr"
)

const (
	confirms     = 1
	btcDecimals  = 8
	btcFeeBlocks = 3
	ChainName    = "Bitcoin"
	Symbol       = "BTC"
)

var _ service.WalletUtXoService = (*BitcoinNodeService)(nil)

type BitcoinNodeService struct {
	btcClient       *base.BtcClient
	btcDataClient   *base.BaseDataClient
	thirdPartClient *BlockChainClient
	unimplemente.UnimplementedService
}

// GetUnspentOutputs 查询某个地址的 UTXO 列表（未花费输出），
func (s *BitcoinNodeService) GetUnspentOutputs(ctx context.Context, req domain.UnspentOutputsParam) ([]domain.UnspentOutput, error) {
	utxoList, err := s.thirdPartClient.GetAccountUtxo(req.Address)
	if err != nil {
		return nil, err
	}
	var unspentOutputList []domain.UnspentOutput
	for _, value := range utxoList {
		unspentOutput := domain.UnspentOutput{
			TxHashBigEndian: value.TxHashBigEndian,
			TxId:            value.TxHash,
			TxOutputN:       value.TxOutputN,
			Script:          value.Script,
			UnspentAmount:   strconv.FormatUint(value.Value, 10),
			Index:           value.TxIndex,
			Address:         req.Address,
			ValueHex:        value.ValueHex,
			Confirmations:   value.Confirmations,
			BlockTime:       value.BlockTime,
		}
		unspentOutputList = append(unspentOutputList, unspentOutput)
	}
	return unspentOutputList, nil
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
	address, err := btcutil.DecodeAddress(param.Address, &chaincfg.MainNetParams)
	if err != nil {
		return false, nil
	}
	if !address.IsForNet(&chaincfg.MainNetParams) {
		return false, nil
	}
	return true, nil
}

// GetBlockByNumber 这段代码是一个用于解析 Bitcoin 区块数据 和 交易详情 的完整实现，包含了区块信息获取、
// 交易列表拉取、UTXO 构建、输入输出金额计算、手续费估算等功能。下面我逐步为你详细解析逻辑和每一步的含义。
func (s *BitcoinNodeService) GetBlockByNumber(ctx context.Context, param domain.BlockNumberParam) (domain.Block, error) {
	blockHash, err := s.btcClient.GetBlockHash(param.Height)
	if err != nil {
		log.Error("get block hash by number fail", "err", err)
		return domain.Block{}, err
	}
	var params []json.RawMessage
	numBlocksJSON, _ := json.Marshal(blockHash)
	params = []json.RawMessage{numBlocksJSON}
	// 通过 getblock 接口拿到 block 的详细数据（包含 txids，但不含交易内容）。
	block, err := s.btcClient.RawRequest("getblock", params)
	if err != nil {
		log.Error("get block by number fail", "err", err)
		return domain.Block{}, err
	}
	var resultBlock BlockData
	err = json.Unmarshal(block, &resultBlock)
	if err != nil {
		log.Error("Unmarshal json fail", "err", err)
	}
	var txList []*domain.BlockTransaction
	//  遍历交易 TxIDs，获取交易内容
	for _, txId := range resultBlock.Tx {
		txIdJson, _ := json.Marshal(txId)
		boolJSON, _ := json.Marshal(true)
		dataJSON := []json.RawMessage{txIdJson, boolJSON}
		// 通过 getrawtransaction(txid, true)，逐个获取每个交易的完整结构（包括 inputs, outputs, scriptSig 等）
		tx, err := s.btcClient.RawRequest("getrawtransaction", dataJSON)
		if err != nil {
			fmt.Println("get raw transaction fail", "err", err)
		}
		var rawTx RawTransactionData
		err = json.Unmarshal(tx, &rawTx)
		if err != nil {
			log.Error("json unmarshal fail", "err", err)
			return domain.Block{}, err
		}
		//  解析交易详情（assembleUtxoTransactionReply）
		/*
			这个函数的作用是：从原始交易数据中提取出：
			输入（vin）：根据输入中的前置交易 txid+vout index 去拉取上一个交易的 output 金额+地址
			输出（vout）：直接解析 value 和地址
			手续费：输入金额 - 输出金额
			返回一个结构化的 UtxoTransaction
			这一步 很关键，构建了完整的 UTXO 交易结构。
		*/
		reply, err := s.assembleUtXoTransactionReply(rawTx, int64(resultBlock.Height), int64(resultBlock.Time),
			func(txId string, index uint32) (int64, string, error) {
				preHash, err2 := chainhash.NewHashFromStr(txId)
				if err2 != nil {
					return 0, "", err2
				}
				preHashJson, _ := json.Marshal(preHash)
				preHashBoolJSON, _ := json.Marshal(true)
				preDataJSON := []json.RawMessage{preHashJson, preHashBoolJSON}
				preTx, err2 := s.btcClient.RawRequest("getrawtransaction", preDataJSON)
				if err2 != nil {
					return 0, "", err2
				}
				var preRawTx RawTransactionData
				err2 = json.Unmarshal(preTx, &preRawTx)
				if err2 != nil {
					log.Error("json unmarshal fail", "err", err2)
					return 0, "", err2
				}
				amount := btcToSatoshi(preRawTx.Vout[index].Value).Int64()

				return amount, preRawTx.Vout[index].ScriptPubKey.Addresses[0], nil
			})
		txList = append(txList, &domain.BlockTransaction{
			Hash:          reply.TxHash,
			Fee:           reply.CostFee,
			Size:          rawTx.Size,
			VSize:         rawTx.VSize,
			Weight:        rawTx.Weight,
			LockTime:      rawTx.LockTime,
			Hex:           rawTx.Hex,
			Version:       rawTx.Version,
			Time:          rawTx.Time,
			BlockHeight:   reply.BlockHeight,
			BlockTime:     reply.BlockTime,
			Blockhash:     rawTx.Blockhash,
			Confirmations: rawTx.Confirmations,
			Status:        TxStatus_name[reply.Status.ToInt32()],
			Vin: slice.Map(reply.Vins, func(idx int, item *VinItem) *domain.Vin {
				return &domain.Vin{
					Hash:    item.Hash,
					Index:   item.Index,
					Amount:  item.Amount,
					Address: item.Address,
				}
			}),
			Vout: slice.Map(reply.Vouts, func(idx int, item *VoutItem) *domain.Vout {
				return &domain.Vout{
					Address: item.Address,
					Amount:  item.Amount,
					Index:   item.Index,
				}
			}),
		})
	}
	return domain.Block{
		Height: uint64(param.Height),
		Hash:   blockHash.String(),
		TxList: txList,
	}, nil
}

// assembleUtXoTransactionReply 把一笔原始 Bitcoin 交易 RawTransactionData 解析为带有输入、
// 输出金额、地址、手续费等完整 UTXO 信息的 UtxoTransaction 结构体。
/*
	tx：当前解析的交易原始数据（含 vin/vout）
	blockHeight / blockTime：这笔交易所在区块的高度与时间戳
	getPrevTxInfo：一个回调函数，用来拿 vin 输入的前置交易的金额和地址
*/
func (s *BitcoinNodeService) assembleUtXoTransactionReply(tx RawTransactionData, blockHeight, blockTime int64, getPrevTxInfo func(txid string, index uint32) (int64, string, error)) (*UtxoTransaction, error) {
	var totalAmountIn, totalAmountOut int64
	ins := make([]*VinItem, 0, len(tx.Vin))
	outs := make([]*VoutItem, 0, len(tx.Vout))
	// 遍历交易的输入（Vin），获取每个输入的金额和地址
	/*
		每个 vin 是对前一笔交易的引用（txid + vout index）
		通过 getPrevTxInfo 取得：
			前一笔的金额（amount）
			对应的地址（address）
	*/
	for _, in := range tx.Vin {
		amount, address, err := getPrevTxInfo(in.TxId, in.Vout)
		if err != nil {
			return nil, err
		}
		totalAmountIn += amount

		t := VinItem{
			Hash:    in.TxId,
			Index:   in.Vout,
			Amount:  amount,
			Address: address,
		}
		ins = append(ins, &t)
	}

	// 处理交易输出（vout）
	/*
		把 BTC 金额转为 satoshi（btcToSatoshi 函数）
		获取 vout.scriptPubKey.addresses[0] 为输出地址
		输出金额累计
	*/
	for index, out := range tx.Vout {
		amount := btcToSatoshi(out.Value).Int64()
		addr := ""
		if len(out.ScriptPubKey.Addresses) > 0 {
			addr = out.ScriptPubKey.Addresses[0]
		}

		totalAmountOut += amount
		t := VoutItem{
			Address: addr,
			Amount:  amount,
			Index:   uint32(index),
		}
		outs = append(outs, &t)
	}
	// Bitcoin 的手续费 = 输入金额 - 输出金额
	gasUsed := totalAmountIn - totalAmountOut
	reply := &UtxoTransaction{
		TxHash:      tx.TxId,
		Status:      TxStatus_Success,
		Vins:        ins,
		Vouts:       outs,
		CostFee:     strconv.FormatInt(gasUsed, 10),
		BlockHeight: uint64(blockHeight),
		BlockTime:   uint64(blockTime),
	}
	return reply, nil
}

/*
btcToSatoshi 作用：将 BTC（浮点）精确地转换为 satoshi（整数），避免浮点误差。
实现逻辑：

	把 float64 格式化成 string，防止精度损失
	使用 decimal 库进行乘法精确计算：BTC * 10^8
	转为 *big.Int 返回
*/
func btcToSatoshi(btcCount float64) *big.Int {
	amount := strconv.FormatFloat(btcCount, 'f', -1, 64)
	amountDm, _ := decimal.NewFromString(amount)
	tenDm := decimal.NewFromFloat(math.Pow(10, float64(btcDecimals)))
	satoshiDm, _ := big.NewInt(0).SetString(amountDm.Mul(tenDm).String(), 10)
	return satoshiDm
}

func (s *BitcoinNodeService) GetBlockByHash(ctx context.Context, param domain.BlockHashParam) (domain.Block, error) {
	var params []json.RawMessage
	numBlocksJSON, err := json.Marshal(param.Hash)
	if err != nil {
		log.Error("marshal block hash fail", "err", err)
		return domain.Block{}, err
	}
	params = []json.RawMessage{numBlocksJSON}
	block, err := s.btcClient.RawRequest("getblock", params)
	if err != nil {
		log.Error("get block by hash fail", "err", err)
		return domain.Block{}, err
	}

	var resultBlock BlockData
	err = json.Unmarshal(block, &resultBlock)
	if err != nil {
		log.Error("Unmarshal json fail", "err", err)
	}
	var txList []*domain.BlockTransaction
	//  遍历交易 TxIDs，获取交易内容
	for _, txId := range resultBlock.Tx {
		txIdJson, _ := json.Marshal(txId)
		boolJSON, _ := json.Marshal(true)
		dataJSON := []json.RawMessage{txIdJson, boolJSON}
		// 通过 getrawtransaction(txid, true)，逐个获取每个交易的完整结构（包括 inputs, outputs, scriptSig 等）
		tx, err := s.btcClient.RawRequest("getrawtransaction", dataJSON)
		if err != nil {
			fmt.Println("get raw transaction fail", "err", err)
		}
		var rawTx RawTransactionData
		err = json.Unmarshal(tx, &rawTx)
		if err != nil {
			log.Error("json unmarshal fail", "err", err)
			return domain.Block{}, err
		}
		reply, err := s.assembleUtXoTransactionReply(rawTx, int64(resultBlock.Height), int64(resultBlock.Time),
			func(txId string, index uint32) (int64, string, error) {
				preHash, err2 := chainhash.NewHashFromStr(txId)
				if err2 != nil {
					return 0, "", err2
				}
				preHashJson, _ := json.Marshal(preHash)
				preHashBoolJSON, _ := json.Marshal(true)
				preDataJSON := []json.RawMessage{preHashJson, preHashBoolJSON}
				preTx, err2 := s.btcClient.RawRequest("getrawtransaction", preDataJSON)
				if err2 != nil {
					return 0, "", err2
				}
				var preRawTx RawTransactionData
				err2 = json.Unmarshal(preTx, &preRawTx)
				if err2 != nil {
					log.Error("json unmarshal fail", "err", err2)
					return 0, "", err2
				}
				amount := btcToSatoshi(preRawTx.Vout[index].Value).Int64()

				return amount, preRawTx.Vout[index].ScriptPubKey.Addresses[0], nil
			})
		txList = append(txList, &domain.BlockTransaction{
			Hash:          reply.TxHash,
			Fee:           reply.CostFee,
			Size:          rawTx.Size,
			VSize:         rawTx.VSize,
			Weight:        rawTx.Weight,
			LockTime:      rawTx.LockTime,
			Hex:           rawTx.Hex,
			Version:       rawTx.Version,
			Time:          rawTx.Time,
			BlockHeight:   reply.BlockHeight,
			BlockTime:     reply.BlockTime,
			Blockhash:     rawTx.Blockhash,
			Confirmations: rawTx.Confirmations,
			Status:        TxStatus_name[reply.Status.ToInt32()],
			Vin: slice.Map(reply.Vins, func(idx int, item *VinItem) *domain.Vin {
				return &domain.Vin{
					Hash:    item.Hash,
					Index:   item.Index,
					Amount:  item.Amount,
					Address: item.Address,
				}
			}),
			Vout: slice.Map(reply.Vouts, func(idx int, item *VoutItem) *domain.Vout {
				return &domain.Vout{
					Address: item.Address,
					Amount:  item.Amount,
					Index:   item.Index,
				}
			}),
		})
	}
	return domain.Block{
		Height: resultBlock.Height,
		Hash:   param.Hash,
		TxList: txList,
	}, nil
}

func (s *BitcoinNodeService) GetBlockHeaderByHash(ctx context.Context, param domain.BlockHeaderHashParam) (domain.BlockHeader, error) {
	hash, err := chainhash.NewHashFromStr(param.Hash)
	if err != nil {
		log.Error("format string to hash fail", "err", err)
	}
	blockHeader, err := s.btcClient.GetBlockHeader(hash)
	if err != nil {
		return domain.BlockHeader{}, err
	}
	return domain.BlockHeader{
		ParentHash: blockHeader.PrevBlock.String(),
		Number:     string(blockHeader.Version),
		BlockHash:  param.Hash,
		MerkleRoot: blockHeader.MerkleRoot.String(),
	}, nil
}

func (s *BitcoinNodeService) GetBlockHeaderByNumber(ctx context.Context, param domain.BlockHeaderNumberParam) (domain.BlockHeader, error) {
	blockNumber := param.Height
	if blockNumber == 0 {
		latestBlock, err := s.btcClient.GetBlockCount()
		if err != nil {
			return domain.BlockHeader{}, err
		}
		blockNumber = latestBlock
	}
	blockHash, err := s.btcClient.GetBlockHash(blockNumber)
	if err != nil {
		log.Error("get block hash by number fail", "err", err)
		return domain.BlockHeader{}, err
	}
	blockHeader, err := s.btcClient.GetBlockHeader(blockHash)
	if err != nil {
		return domain.BlockHeader{}, err
	}
	return domain.BlockHeader{
		ParentHash: blockHeader.PrevBlock.String(),
		Number:     string(blockHeader.Version),
		BlockHash:  blockHash.String(),
		MerkleRoot: blockHeader.MerkleRoot.String(),
	}, nil
}

func (s *BitcoinNodeService) GetBalanceByAddress(ctx context.Context, param domain.BalanceByAddressParam) (domain.Balance, error) {
	balance, err := s.thirdPartClient.GetBalanceByAddress(param.Address)
	if err != nil {
		return domain.Balance{}, err
	}
	return domain.Balance{Balance: balance}, nil
}

// GetFee 获取 BTC 网络手续费（Fee）推荐值
func (s *BitcoinNodeService) GetFee(ctx context.Context, param domain.FeeParam) (domain.Fee, error) {
	gasFeeResp, err := s.btcDataClient.GetFee(Symbol)
	if err != nil {
		return domain.Fee{}, err
	}
	return domain.Fee{
		BestFee:    gasFeeResp.BestTransactionFee,
		BestFeeSat: gasFeeResp.BestTransactionFeeSat,
		SlowFee:    gasFeeResp.SlowGasPrice,
		NormalFee:  gasFeeResp.StandardGasPrice,
		FastFee:    gasFeeResp.RapidGasPrice,
	}, nil
}

func (s *BitcoinNodeService) SendTx(ctx context.Context, param domain.SendTxParam) (string, error) {
	r := bytes.NewReader([]byte(param.RawTx))
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(r)
	if err != nil {
		return "", err
	}
	txHash, err := s.btcClient.SendRawTransaction(&msgTx, true)
	if err != nil {
		return "", err
	}
	if strings.Compare(msgTx.TxHash().String(), txHash.String()) != 0 {
		log.Error("broadcast transaction, tx hash mismatch", "local hash", msgTx.TxHash().String(), "hash from net", txHash.String(), "signedTx", param.RawTx)
	}
	return txHash.String(), nil
}

// ListTxByAddress 查询某个地址的历史交易记录（分页），并将每一笔交易转换成统一格式 TxMessage 返回
func (s *BitcoinNodeService) ListTxByAddress(ctx context.Context, param domain.TxAddressParam) ([]domain.TxMessage, error) {
	// thirdPartClient：调用外部服务（如 Blockstream、Mempool.space、BitGo 等）查询地址相关交易。
	//返回结果里通常有：交易 hash、input、output、block 高度、时间、手续费等。
	transaction, err := s.thirdPartClient.GetTransactionsByAddress(param.Address,
		strconv.Itoa(int(param.Page)), strconv.Itoa(int(param.Pagesize)))
	if err != nil {
		return nil, err
	}
	var txMessages []domain.TxMessage
	for _, txItems := range transaction.Txs {
		var fromAddrs []string
		var toAddrs []string
		var values []domain.Value
		for _, inputs := range txItems.Inputs {
			fromAddrs = append(fromAddrs, inputs.PrevOut.Addr)
			values = append(values, domain.Value{
				Address: inputs.PrevOut.Addr,
				Value:   inputs.PrevOut.Value.String(),
			})
		}
		txFee := txItems.Fee
		for _, out := range txItems.Out {
			toAddrs = append(toAddrs, out.Addr)
			values = append(values, domain.Value{
				Address: out.Addr,
				Value:   out.Value.String(),
			})
		}
		// 方向判断（是转出还是转入）
		datetime := txItems.Time.String()
		direction := int32(1)
		for _, fromAddr := range fromAddrs {
			if strings.EqualFold(fromAddr, param.Address) {
				direction = 0 // 出现在 inputs 中就是转出
				break
			}
		}
		txMessages = append(txMessages, domain.TxMessage{
			Hash:     txItems.Hash,
			Froms:    fromAddrs,
			Tos:      toAddrs,
			Values:   values,
			Fee:      txFee.String(),
			Status:   domain.TxStatus_Success,
			Type:     direction,
			Height:   txItems.BlockHeight.String(),
			Datetime: datetime,
		})
	}
	return txMessages, nil
}

func (s *BitcoinNodeService) GetTxByHash(ctx context.Context, param domain.GetTxByHashParam) (domain.TxMessage, error) {
	transaction, err := s.thirdPartClient.GetTransactionsByHash(param.Hash)
	if err != nil {
		return domain.TxMessage{}, err
	}
	var fromAddrs []string
	var toAddrs []string
	var values []domain.Value
	for _, inputs := range transaction.Inputs {
		fromAddrs = append(fromAddrs, inputs.PrevOut.Addr)
		values = append(values, domain.Value{
			Address: inputs.PrevOut.Addr,
			Value:   inputs.PrevOut.Value.String(),
		})
	}
	txFee := transaction.Fee
	for _, out := range transaction.Out {
		toAddrs = append(toAddrs, out.Addr)
		values = append(values, domain.Value{
			Address: out.Addr,
			Value:   out.Value.String(),
		})
	}
	datetime := transaction.Time.String()
	// 方向判断（是转出还是转入）
	return domain.TxMessage{
		Hash:     transaction.Hash,
		Froms:    fromAddrs,
		Tos:      toAddrs,
		Values:   values,
		Fee:      txFee.String(),
		Status:   domain.TxStatus_Success,
		Height:   transaction.BlockHeight.String(),
		Datetime: datetime,
	}, nil
}

//func (s *BitcoinNodeService) assembleUtXoTransactionReplyForTxHash(tx *btcjson.TxRawResult, blockHeight int64,
//	getPrevTxInfo func(txid string, index uint32) (int64, string, error)) (*wallet2.TxHashResponse, error) {
//	var totalAmountIn, totalAmountOut int64
//	var from_addrs []*wallet2.Address
//	var to_addrs []*wallet2.Address
//	var value_list []*wallet2.Value
//	var direction int32
//	for _, in := range tx.Vin {
//		amount, address, err := getPrevTxInfo(in.Txid, in.Vout)
//		if err != nil {
//			return nil, err
//		}
//		totalAmountIn += amount
//		from_addrs = append(from_addrs, &wallet2.Address{Address: address})
//		value_list = append(value_list, &wallet2.Value{Value: strconv.FormatInt(totalAmountIn, 10)})
//	}
//	for _, out := range tx.Vout {
//		amount := btcToSatoshi(out.Value).Int64()
//		totalAmountOut += amount
//		addr := ""
//		if len(out.ScriptPubKey.Addresses) > 0 {
//			addr = out.ScriptPubKey.Addresses[0]
//		}
//		to_addrs = append(to_addrs, &wallet2.Address{Address: addr})
//		value_list = append(value_list, &wallet2.Value{Value: strconv.FormatInt(totalAmountOut, 10)})
//	}
//	gasUsed := totalAmountIn - totalAmountOut
//	return wallet2.TxMessage{
//		Hash:   tx.Hash,
//		Status: wallet2.TxStatus_Success,
//		Froms:  from_addrs,
//		Tos:    to_addrs,
//		Fee:    strconv.FormatInt(gasUsed, 10),
//		Values: value_list,
//		Height: strconv.FormatInt(blockHeight, 10),
//		Type:   direction,
//	}, nil
//}

// CreateUnSignTransaction 构建一个 比特币未签名交易（Unsigned Tx）生成接口，用于支持离线签名、冷签名等场景。
func (s *BitcoinNodeService) CreateUnSignTransaction(_ context.Context, param domain.UnSignTransactionParam) (domain.UnSignTransactionResult, error) {
	txHash, buf, err := s.CalcSignHashes(param.Vin, param.Vout)
	if err != nil {
		log.Error("calc sign hashes fail", "err", err)
		return domain.UnSignTransactionResult{}, err
	}
	return domain.UnSignTransactionResult{
		TxData:     buf,
		SignHashes: txHash,
	}, nil
}

// CalcSignHashes 根据传入的 UTXO（Vin）和目标地址输出（Vout），
// 构造出一笔 RawTx，并为每个 Vin 生成签名哈希（signHash），以供后续进行签名。
// 调用者可以拿到：
// SignHashes：每个 input 的签名哈希（需要私钥签名）
// TxData：未序列化的交易结构体数据（这里其实是空的，有个 BUG，详见下文）
func (s *BitcoinNodeService) CalcSignHashes(Vins []domain.Vin, Vouts []domain.Vout) ([][]byte, []byte, error) {
	if len(Vins) == 0 || len(Vouts) == 0 {
		return nil, nil, errors.New("invalid len in or out")
	}
	// 构建原始交易结构 rawTx， 这就是一个未签名的 MsgTx 对象。
	rawTx := wire.NewMsgTx(wire.TxVersion)
	/*
		构建交易输入
		每个 Vin 都构建为一个 TxIn，其引用的是：
			前一个交易的 txid 和 pre vout index
			这里暂时不带 ScriptSig（因为还没签名）
	*/
	for _, in := range Vins {
		utxoHash, err := chainhash.NewHashFromStr(in.Hash)
		if err != nil {
			return nil, nil, err
		}
		txIn := wire.NewTxIn(wire.NewOutPoint(utxoHash, in.Index), nil, nil)
		rawTx.AddTxIn(txIn)
	}
	/*
		构建交易输出
		每个 Vout 构建为一个 TxOut：
			地址转为公钥脚本 PayToAddrScript
			添加到交易中
	*/
	for _, out := range Vouts {
		toAddress, err := btcutil.DecodeAddress(out.Address, &chaincfg.MainNetParams)
		if err != nil {
			return nil, nil, err
		}
		toPkScript, err := txscript.PayToAddrScript(toAddress)
		if err != nil {
			return nil, nil, err
		}
		rawTx.AddTxOut(wire.NewTxOut(out.Amount, toPkScript))
	}
	signHashes := make([][]byte, len(Vins))
	/*
		为每个输入生成签名哈希
			每个 input 的签名哈希是针对该 input 所引用的 UTXO 地址脚本计算的
			采用 SigHashAll：表示签名这整个交易（最常见的方式）
	*/
	for i, in := range Vins {
		from := in.Address
		fromAddr, err := btcutil.DecodeAddress(from, &chaincfg.MainNetParams)
		if err != nil {
			log.Info("decode address error", "from", from, "err", err)
			return nil, nil, err
		}
		fromPkScript, err := txscript.PayToAddrScript(fromAddr)
		if err != nil {
			log.Info("pay to addr script err", "err", err)
			return nil, nil, err
		}
		signHash, err := txscript.CalcSignatureHash(fromPkScript, txscript.SigHashAll, rawTx, i)
		if err != nil {
			log.Info("Calc signature hash error", "err", err)
			return nil, nil, err
		}
		signHashes[i] = signHash
	}
	// 交易序列化数据
	buf := bytes.NewBuffer(make([]byte, 0, rawTx.SerializeSize()))
	err := rawTx.Serialize(buf)
	if err != nil {
		return nil, nil, err
	}
	return signHashes, buf.Bytes(), nil
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
