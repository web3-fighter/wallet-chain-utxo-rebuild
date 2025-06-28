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

func (s *BitcoinNodeService) ListTxByAddress(ctx context.Context, param domain.TxAddressParam) ([]domain.TxMessage, error) {
	transaction, err := s.thirdPartClient.GetTransactionsByAddress(param.Address,
		strconv.Itoa(int(param.Page)), strconv.Itoa(int(param.Pagesize)))
	if err != nil {
		return nil, err
	}
	var txMessages []domain.TxMessage
	for _, txItems := range transaction.Txs {
		var fromAddrs []string
		var toAddrs []string
		var values []string
		var direction int32
		for _, inputs := range txItems.Inputs {
			fromAddrs = append(fromAddrs, inputs.PrevOut.Addr)
		}
		txFee := txItems.Fee
		for _, out := range txItems.Out {
			toAddrs = append(toAddrs, out.Addr)
			values = append(values, out.Value.String())
		}
		datetime := txItems.Time.String()
		if strings.EqualFold(param.Address, fromAddrs[0]) {
			direction = 0
		} else {
			direction = 1
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
