package rollup

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-resty/resty/v2"
)

/**
 * GET /enqueue/index/{index}
 * GET /transaction/index/{index}
 * GET /batch/transaction/index/{index}
 * GET /stateroot/index/{index}
 * GET /batch/stateroot/index/{index}
 */

type Batch struct {
	Index             uint64         `json:"index"`
	Root              common.Hash    `json:"root,omitempty"`
	Size              uint32         `json:"size,omitempty"`
	PrevTotalElements uint32         `json:"prevTotalElements,omitempty"`
	ExtraData         hexutil.Bytes  `json:"extraData,omitempty"`
	BlockNumber       uint64         `json:"blockNumber"`
	Timestamp         uint64         `json:"timestamp"`
	Submitter         common.Address `json:"submitter"`
}

type stateRoot struct {
	Index uint64   `json:"index"`
	Value [32]byte `json:"value"`
}

type stateRoots []stateRoot

type message struct {
	Target   common.Address `json:"target"`
	Data     hexutil.Bytes  `json:"data"`
	GasLimit uint64         `json:"gasLimit"`
}

type EthContext struct {
	BlockNumber uint64 `json:"blockNumber"`
	Timestamp   uint64 `json:"timestamp"`
}

type transaction struct {
	Index       uint64         `json:"index"`
	BatchIndex  uint64         `json:"batchIndex"`
	BlockNumber uint64         `json:"blockNumber"`
	Timestamp   uint64         `json:"timestamp"`
	GasLimit    uint64         `json:"gasLimit"`
	Target      common.Address `json:"target"`
	Origin      common.Address `json:"origin"`
	Data        hexutil.Bytes  `json:"data"`
	QueueOrigin string         `json:"queueOrigin"`
	Type        string         `json:"type"`
	QueueIndex  uint64         `json:"queueIndex"`
	Decoded     *decoded       `json:"decoded"`
}

type Enqueue struct {
	Index       uint64         `json:"index"`
	Message     message        `json:"message"`
	Data        hexutil.Bytes  `json"data"`
	GasLimit    uint64         `json"gasLimit"`
	Origin      common.Address `json"origin"`
	BlockNumber uint64         `json"blockNumber"`
	Timestamp   uint64         `json"timestamp"`
}

type signature struct {
	R hexutil.Bytes `json:"r"`
	S hexutil.Bytes `json:"s"`
	V hexutil.Bytes `json:"v"`
}

type decoded struct {
	Signautre signature      `json:"sig"`
	GasLimit  uint64         `json:"gasLimit"`
	GasPrice  uint64         `json:"gasPrice"`
	Nonce     uint64         `json:"nonce"`
	Target    common.Address `json:"target"`
	Data      hexutil.Bytes  `json:"data"`
}

type Client struct {
	client *resty.Client
}

type TransactionResponse struct {
	Transaction transaction `json:"transaction"`
	Batch       Batch       `json:"batch"`
}

type BatchResponse struct {
	Batch        Batch         `json:"batch"`
	Transactions []transaction `json:"transactions"`
}

func NewClient(url string, confirmations uint) *Client {
	client := resty.New()
	client.SetHostURL(url)
	s := strconv.Itoa(int(confirmations))
	client.SetQueryParam("confirmations", s)

	return &Client{
		client: client,
	}
}

func (c *Client) GetEnqueue(index uint64) (*Enqueue, error) {
	str := strconv.FormatUint(index, 10)
	response, err := c.client.R().
		SetPathParams(map[string]string{
			"index": str,
		}).
		SetResult(&Enqueue{}).
		Get("/enqueue/index/{index}")

	if err != nil {
		return nil, err
	}
	enqueue, ok := response.Result().(*Enqueue)
	if !ok {
		return nil, errors.New("")
	}
	return enqueue, nil
}

func (c *Client) GetTransaction(index uint64) (*types.Transaction, error) {
	str := strconv.FormatUint(index, 10)
	response, err := c.client.R().
		SetPathParams(map[string]string{
			"index": str,
		}).
		SetResult(&TransactionResponse{}).
		Get("/transaction/index/{index}")

	if err != nil {
		return nil, err
	}
	res, ok := response.Result().(*TransactionResponse)
	if !ok {
		return nil, errors.New("")
	}

	if res.Transaction.Decoded != nil {
		nonce := res.Transaction.Decoded.Nonce
		to := res.Transaction.Target
		value := new(big.Int)
		// Note: there are two gas limits, one top level and
		// another on the raw transaction itself. Maybe maxGasLimit
		// for the top level?
		gasLimit := res.Transaction.Decoded.GasLimit
		gasPrice := new(big.Int).SetUint64(res.Transaction.Decoded.GasPrice)
		data := res.Transaction.Decoded.Data
		l1MessageSender := res.Transaction.Origin
		l1BlockNumber := new(big.Int).SetUint64(res.Transaction.BlockNumber)
		// The queue origin must be either sequencer of l1, otherwise
		// it is considered an unknown queue origin and will not be processed
		var queueOrigin types.QueueOrigin
		if res.Transaction.QueueOrigin == "sequencer" {
			queueOrigin = types.QueueOriginSequencer
		} else if res.Transaction.QueueOrigin == "l1" {
			queueOrigin = types.QueueOriginL1ToL2
		} else {
			return nil, fmt.Errorf("Unknown queue origin: %s", res.Transaction.QueueOrigin)
		}
		// The transaction type must be EIP155 or EthSign. Throughout this
		// codebase, it is referred to as "sighash type" but it could actually
		// be generalized to transaction type. Right now the only different
		// types use a different signature hashing scheme.
		var sighashType types.SignatureHashType
		if res.Transaction.Type == "EIP155" {
			sighashType = types.SighashEIP155
		} else if res.Transaction.Type == "ETH_Sign" {
			sighashType = types.SighashEthSign
		} else {
			return nil, fmt.Errorf("Unknown transaction type: %s", res.Transaction.Type)
		}

		tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data, &l1MessageSender, l1BlockNumber, queueOrigin, sighashType)
		return tx, nil
	}

	return nil, fmt.Errorf("No decoded transaction: %s", res.Transaction.Type)
}

func (c *Client) GetEthContext(index uint64) (*EthContext, error) {
	str := strconv.FormatUint(index, 10)
	response, err := c.client.R().
		SetPathParams(map[string]string{
			"index": str,
		}).
		SetResult(&EthContext{}).
		Get("/eth/context/latest")

	if err != nil {
		return nil, fmt.Errorf("Cannot fetch eth context: %w", err)
	}

	context, ok := response.Result().(*EthContext)
	if !ok {
		return nil, errors.New("Cannot parse EthContext")
	}

	return context, nil
}
