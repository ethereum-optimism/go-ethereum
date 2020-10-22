package rollup

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	ctc "github.com/ethereum/go-ethereum/contracts/canonicaltransactionchain"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

// Mock deployed address of canonical transaction chain
var ctcAddress = common.HexToAddress("0xE894780e35530557B152281e8828339303aE33e5")

func TestSyncServiceDatabase(t *testing.T) {
	service, err := newTestSyncService()
	if err != nil {
		t.Fatal(err)
	}

	mockEthClient(service, map[string]interface{}{})
	mockLogClient(service, [][]types.Log{})

	go service.Loop()

	headers := []types.Header{
		{Number: big.NewInt(1)},
		{Number: big.NewInt(2)},
	}

	for _, header := range headers {
		service.heads <- &header

		height := <-service.doneProcessing
		if height != header.Number.Uint64() {
			t.Fatal("Wrong height received")
		}

		// The lastestEth1Data should be kept up to data
		if service.Eth1Data.BlockHeight != header.Number.Uint64() {
			t.Fatalf("Mismatched eth1 data blockheight: got %d, expect %d", service.Eth1Data.BlockHeight, header.Number.Uint64())
		}
		if !bytes.Equal(service.Eth1Data.BlockHash.Bytes(), header.Hash().Bytes()) {
			t.Fatalf("Mismatched eth1 blockhash")
		}

		// The database should be kept up to date
		eth1data := service.GetLastProcessedEth1Data()
		if eth1data.BlockHeight != height {
			t.Fatal("Wrong height in database")
		}
		if !bytes.Equal(eth1data.BlockHash.Bytes(), header.Hash().Bytes()) {
			t.Fatal("Wrong hash in database")
		}
	}
}

func mustABINewType(s string) abi.Type {
	typ, err := abi.NewType(s, s, []abi.ArgumentMarshaling{})
	if err != nil {
		fmt.Println(err)
	}
	return typ
}

func abiEncodeCTCEnqueued(origin, target *common.Address, gasLimit, queueIndex, timestamp *big.Int, data []byte) []byte {
	args := abi.Arguments{
		{Name: "l1TxOrigin", Type: mustABINewType("address")},
		{Name: "target", Type: mustABINewType("address")},
		{Name: "gasLimit", Type: mustABINewType("uint256")},
		{Name: "data", Type: mustABINewType("bytes")},
		{Name: "queueIndex", Type: mustABINewType("uint256")},
		{Name: "timestamp", Type: mustABINewType("uint256")},
	}
	raw, err := args.PackValues([]interface{}{
		origin,
		target,
		gasLimit,
		data,
		queueIndex,
		timestamp,
	})
	if err != nil {
		fmt.Printf("Cannot abi encode: %s", err)
		return []byte{}
	}
	return raw
}

// Can be used for both queue batch appended and sequencer batch appended
func abiEncodeBatchAppended(startingQueueIndex, numQueueElements, totalElements *big.Int) []byte {
	args := abi.Arguments{
		{Name: "startingQueueIndex", Type: mustABINewType("uint256")},
		{Name: "numQueueElements", Type: mustABINewType("uint256")},
		{Name: "totalElements", Type: mustABINewType("uint256")},
	}
	raw, err := args.PackValues([]interface{}{
		startingQueueIndex,
		numQueueElements,
		totalElements,
	})
	if err != nil {
		fmt.Printf("Cannot abi encode: %s", err)
		return []byte{}
	}
	return raw
}

// Test that the `RollupTransaction` ends up in the transaction cache
// after the transaction enqueued event is emitted.
func TestSyncServiceTransactionEnqueued(t *testing.T) {
	service, err := newTestSyncService()
	if err != nil {
		t.Fatal(err)
	}

	// The queue index is used as the key in the transaction cache
	queueIndex := big.NewInt(0)
	// The timestamp is in the rollup transaction
	timestamp := big.NewInt(24)
	// The target is the `to` field on the transaction
	target := common.HexToAddress("0x04668ec2f57cc15c381b461b9fedab5d451c8f7f")
	// The layer one transaction origin is in the txmeta on the transaction
	l1TxOrigin := common.HexToAddress("0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8")
	// The gasLimit is the `gasLimit` on the transaction
	gasLimit := big.NewInt(66)
	// The data is the `data` on the transaction
	data := []byte{0x02, 0x92}

	mockEthClient(service, map[string]interface{}{})
	mockLogClient(service, [][]types.Log{
		{
			{
				Address:     ctcAddress,
				BlockNumber: 1,
				Topics: []common.Hash{
					common.BytesToHash(transactionEnqueuedEventSignature),
				},
				Data: abiEncodeCTCEnqueued(&l1TxOrigin, &target, gasLimit, queueIndex, timestamp, data),
			},
		},
	})

	// Start up the main loop
	go service.Loop()

	service.heads <- &types.Header{Number: big.NewInt(1)}
	_ = <-service.doneProcessing

	rtx, ok := service.txCache.Load(queueIndex.Uint64())
	if !ok {
		t.Fatal("Transaction not found in cache")
	}

	// The timestamps should be equal
	if big.NewInt(rtx.timestamp.Unix()).Cmp(timestamp) != 0 {
		t.Fatal("Incorrect time recovered")
	}

	// The target from the calldata should be the `to` in the transaction
	if !bytes.Equal(rtx.tx.To().Bytes(), target.Bytes()) {
		t.Fatal("Incorrect target")
	}

	if !bytes.Equal(rtx.tx.L1MessageSender().Bytes(), l1TxOrigin.Bytes()) {
		t.Fatal("L1TxOrigin not set correctly")
	}

	if rtx.tx.Gas() != gasLimit.Uint64() {
		t.Fatal("Incorrect gas limit")
	}

	if !bytes.Equal(rtx.tx.Data(), data) {
		t.Fatal("Incorrect data")
	}
}

// Tests that a queue batch append results in the transaction
// from the cache is played against the state.
func TestSyncServiceQueueBatchAppend(t *testing.T) {
	service, err := newTestSyncService()
	if err != nil {
		t.Fatal(err)
	}

	// The queue index is 0 as well as the starting queue index below.
	// These must match for this to work.
	queueIndex, timestamp, gasLimit := big.NewInt(0), big.NewInt(97538), big.NewInt(210000)
	target := common.HexToAddress("0x04668ec2f57cc15c381b461b9fedab5d451c8f7f")
	l1TxOrigin := common.HexToAddress("0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8")
	data := []byte{0x02, 0x92}

	startingQueueIndex := big.NewInt(0)
	numQueueElements := big.NewInt(1)
	totalElements := big.NewInt(0)

	mockEthClient(service, map[string]interface{}{})
	mockLogClient(service, [][]types.Log{
		{
			// This transaction will end up in the tx cache
			{
				Address:     ctcAddress,
				BlockNumber: 1,
				Topics: []common.Hash{
					common.BytesToHash(transactionEnqueuedEventSignature),
				},
				Data: abiEncodeCTCEnqueued(&l1TxOrigin, &target, gasLimit, queueIndex, timestamp, data),
			},
			// This should pull the tx out of the tx cache and then play it evaluate it
			{
				Address:     ctcAddress,
				BlockNumber: 1,
				Topics: []common.Hash{
					common.BytesToHash(queueBatchAppendedEventSignature),
				},
				Data: abiEncodeBatchAppended(startingQueueIndex, numQueueElements, totalElements),
			},
		},
	})

	go service.Loop()

	service.heads <- &types.Header{Number: big.NewInt(1)}
	_ = <-service.doneProcessing
	rtx, _ := service.txCache.Load(queueIndex.Uint64())

	ok, err := txProcessed(t, rtx, service)
	if !ok {
		t.Fatal(err)
	}
}

func txProcessed(t *testing.T, rtx *RollupTransaction, service *SyncService) (bool, error) {
	// Due to the current architecture of the system, the transaction should end
	// up in the mempool. Downstream services are responsible for applying
	// transactions to the state from the mempool.
	pending, _ := service.txpool.Pending()
	count := 0
	for from, txs := range pending {
		// The from should be the god key
		if bytes.Equal(from.Bytes(), service.address.Bytes()) {
			if len(txs) != 1 {
				t.Fatal("More transactions in mempool than expected")
			}
			tx := txs[0]
			//fmt.Println(tx.Hash().Hex())

			if rtx.tx.Nonce() != tx.Nonce() {
				t.Fatal("Nonce mismatch")
			}
			if !bytes.Equal(rtx.tx.To().Bytes(), tx.To().Bytes()) {
				t.Fatal("To mismatch")
			}
			if rtx.tx.Gas() != tx.Gas() {
				t.Fatal("Gas mismatch")
			}
			if !bytes.Equal(rtx.tx.GasPrice().Bytes(), tx.GasPrice().Bytes()) {
				t.Fatal("GasPrice mismatch")
			}
			if !bytes.Equal(rtx.tx.Value().Bytes(), tx.Value().Bytes()) {
				t.Fatal("Value mismatch")
			}
			if !bytes.Equal(rtx.tx.Data(), tx.Data()) {
				t.Fatal("Data mismatch")
			}
			// remove the signature from the tx by creating a new tx with all
			// of the information and then compare hashes.
			fresh := types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), tx.Data(), nil, nil, types.QueueOriginL1ToL2, types.SighashEIP155)

			if !bytes.Equal(fresh.Hash().Bytes(), rtx.tx.Hash().Bytes()) {
				t.Fatal("Hash mismatch")
			}
		}
		// Keep track of all pending tranasctions
		count++
	}

	// There should only be one transaction in the mempool
	if count != 1 {
		t.Fatal("More transactions in mempool than expected")
	}

	return true, nil
}

func TestSyncServiceSequencerBatchAppend(t *testing.T) {
	service, err := newTestSyncService()
	if err != nil {
		t.Fatal(err)
	}

	raw := hexutil.MustDecode("0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301")
	var sig [65]byte
	copy(sig[:], raw)

	ctcTx := CTCTransaction{
		typ: CTCTransactionTypeEIP155,
		tx: &CTCTxEIP155{
			Signature: sig,
			gasLimit:  50000,
			gasPrice:  0,
			nonce:     (1 << 21) + 14,
			target:    common.HexToAddress("0x5769785087b1b64e4cbd9a38d48a1ca35a2fd75cf5cd941d75b2e2fbc6018e8a"),
			data:      raw,
		},
	}

	length, _ := ctcTx.Len()
	txdata := make([]byte, length)
	err = ctcTx.Encode(txdata)
	if err != nil {
		t.Fatal(err)
	}

	cd := appendSequencerBatchCallData{
		ChainElements: []chainElement{
			{
				IsSequenced: true,
				Timestamp:   big.NewInt(1602820447),
				BlockNumber: big.NewInt(0),
				TxData:      txdata,
			},
		},
		Contexts: []ctcBatchContext{
			{
				NumSequencedTransactions:       big.NewInt(1),
				NumSubsequentQueueTransactions: big.NewInt(0),
				Timestamp:                      big.NewInt(1602820447),
				BlockNumber:                    big.NewInt(0),
			},
		},
		ShouldStartAtBatch:    big.NewInt(0),
		TotalElementsToAppend: big.NewInt(1),
	}

	rawCd := new(bytes.Buffer)
	err = cd.Encode(rawCd)
	if err != nil {
		t.Fatal(err)
	}

	// get transaction by hash
	mockEthClient(service, map[string]interface{}{
		"TransactionByHash": []*types.Transaction{
			types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), rawCd.Bytes(), nil, nil, types.QueueOriginSequencer, types.SighashEIP155),
		},
	})
	mockLogClient(service, [][]types.Log{
		{
			{
				Address:     ctcAddress,
				BlockNumber: 1,
				Topics: []common.Hash{
					common.BytesToHash(sequencerBatchAppendedEventSignature),
				},
				Data: abiEncodeBatchAppended(big.NewInt(0), big.NewInt(1), big.NewInt(1)),
			},
		},
	})

	go service.Loop()

	service.heads <- &types.Header{Number: big.NewInt(1)}
	_ = <-service.doneProcessing

	// The god key should be a key containing the transaction in queued
	_, queued := service.txpool.Content()

	// queued should only have 1 key/value pair
	for _, transactions := range queued {
		// This stuff should match ctcTx.tx
		tx := transactions[0]
		data := tx.Data()
		to := tx.To().Hex()
		fmt.Printf("%#v\n", data)
		fmt.Printf("%#v\n", to)
	}
}

func newTestSyncService() (*SyncService, error) {
	chainCfg := params.AllEthashProtocolChanges
	chainID := big.NewInt(420)
	chainCfg.ChainID = chainID

	engine := ethash.NewFaker()
	db := rawdb.NewMemoryDatabase()
	_ = new(core.Genesis).MustCommit(db)
	chain, err := core.NewBlockChain(db, nil, chainCfg, engine, vm.Config{}, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize blockchain: %w", err)
	}
	chaincfg := params.ChainConfig{ChainID: chainID}

	txPool := core.NewTxPool(core.TxPoolConfig{PriceLimit: 0}, &chaincfg, chain)

	// Hardcoded god key for determinism
	d := "0xcb27a3fd66eeb29699d37c860f4b3545dad264aa70d2afdd92a454f30e3ae560"
	key, err = crypto.ToECDSA(hexutil.MustDecode(d))

	cfg := Config{
		CanonicalTransactionChainDeployHeight: big.NewInt(0),
		CanonicalTransactionChainAddress:      ctcAddress,
		TxIngestionSignerKey:                  key,
	}

	service, err := NewSyncService(context.Background(), cfg, txPool, chain, db)
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize syncservice: %w", err)
	}

	return service, nil
}

// Mock setup functions
func mockLogClient(service *SyncService, logs [][]types.Log) {
	service.logClient = newMockBoundCTCContract(logs)
	ctcFilterer, _ := ctc.NewOVMCanonicalTransactionChainFilterer(ctcAddress, service.logClient)
	service.ctcFilterer = ctcFilterer
}

func mockEthClient(service *SyncService, responses map[string]interface{}) {
	service.ethclient = newMockEthereumClient(responses)
}

// Test utilities
type mockEthereumClient struct {
	transactionByHashCallCount int
	transactionByHashResponses []*types.Transaction
}

func (m *mockEthereumClient) ChainID(context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}
func (m *mockEthereumClient) NetworkID(context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}
func (m *mockEthereumClient) SyncProgress(context.Context) (*ethereum.SyncProgress, error) {
	sp := ethereum.SyncProgress{}
	return &sp, nil
}
func (m *mockEthereumClient) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	h := types.Header{}
	return &h, nil
}
func (m *mockEthereumClient) TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error) {
	if m.transactionByHashCallCount < len(m.transactionByHashResponses) {
		res := m.transactionByHashResponses[m.transactionByHashCallCount]
		return res, false, nil
	}
	t := types.Transaction{}
	return &t, false, nil
}

func newMockEthereumClient(responses map[string]interface{}) *mockEthereumClient {
	transactionByHashResponses := []*types.Transaction{}

	txByHash, ok := responses["TransactionByHash"]
	if ok {
		transactionByHashResponses = txByHash.([]*types.Transaction)
	}
	return &mockEthereumClient{
		transactionByHashResponses: transactionByHashResponses,
	}
}

type mockBoundCTCContract struct {
	filterLogsResponses [][]types.Log
	filterLogsCallCount int
}

func (m *mockBoundCTCContract) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	if m.filterLogsCallCount < len(m.filterLogsResponses) {
		res := m.filterLogsResponses[m.filterLogsCallCount]
		m.filterLogsCallCount++
		return res, nil
	}
	return []types.Log{}, nil
}
func (m *mockBoundCTCContract) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return newMockSubscription(), nil
}
func newMockBoundCTCContract(responses [][]types.Log) *mockBoundCTCContract {
	return &mockBoundCTCContract{
		filterLogsResponses: responses,
	}
}

type mockSubscription struct {
	e <-chan error
}

func (m *mockSubscription) Unsubscribe() {}
func (m *mockSubscription) Err() <-chan error {
	return m.e
}
func newMockSubscription() *mockSubscription {
	e := make(chan error)
	return &mockSubscription{
		e: e,
	}
}
