package rollup

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
)

// These variables represent the event signatures
var (
	transactionEnqueuedEventSignature      = crypto.Keccak256([]byte("TransactionEnqueued(address,address,uint256,bytes,uint256,uint256)"))
	queueBatchAppendedEventSignature       = crypto.Keccak256([]byte("QueueBatchAppended(uint256,uint256,uint256)"))
	sequencerBatchAppendedEventSignature   = crypto.Keccak256([]byte("SequencerBatchAppended(uint256,uint256,uint256)"))
	transactionBatchAppendedEventSignature = crypto.Keccak256([]byte("TransactionBatchAppended(uint256,bytes32,uint256,uint256,bytes)"))
)

// Test that the `RollupTransaction` ends up in the transaction cache
// after the transaction enqueued event is emitted. Set `false` as
// the argument to start as a sequencer
func TestSyncServiceTransactionEnqueued(t *testing.T) {
	service, txCh, _, err := newTestSyncService(false)
	if err != nil {
		t.Fatal(err)
	}

	// The timestamp is in the rollup transaction
	timestamp := uint64(24)
	// The target is the `to` field on the transaction
	target := common.HexToAddress("0x04668ec2f57cc15c381b461b9fedab5d451c8f7f")
	// The layer one transaction origin is in the txmeta on the transaction
	l1TxOrigin := common.HexToAddress("0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8")
	// The gasLimit is the `gasLimit` on the transaction
	gasLimit := uint64(66)
	// The data is the `data` on the transaction
	data := []byte{0x02, 0x92}
	// The L1 blocknumber for the transaction's evm context
	l1BlockNumber := big.NewInt(100)
	// The queue index of the L1 to L2 transaction
	queueIndex := uint64(100)
	//The index in the ctc
	index := uint64(120)

	tx := types.NewTransaction(0, target, big.NewInt(0), gasLimit, big.NewInt(0), data, &l1TxOrigin, l1BlockNumber, types.QueueOriginL1ToL2, types.SighashEIP155)
	meta := types.TransactionMeta{
		L1BlockNumber:     l1BlockNumber,
		L1Timestamp:       timestamp,
		L1MessageSender:   &l1TxOrigin,
		SignatureHashType: types.SighashEIP155,
		QueueOrigin:       big.NewInt(int64(types.QueueOriginL1ToL2)),
		Index:             &index,
		QueueIndex:        &queueIndex,
	}
	tx.SetTransactionMeta(&meta)

	setupMockClient(service, map[string]interface{}{
		"GetEnqueue": []*types.Transaction{
			tx,
		},
	})

	// Start up the main loop
	go service.Loop()
	// Wait for the tx to be confirmed into the chain and then
	// make sure it is the transactions that was set up with in the mockclient
	event := <-txCh
	if len(event.Txs) != 1 {
		t.Fatal("Unexpected number of transactions")
	}
	confirmed := event.Txs[0]

	if !reflect.DeepEqual(tx, confirmed) {
		t.Fatal("different txs")
	}
}

// Pass true to set as a verifier
func TestSyncServiceSync(t *testing.T) {
	service, txCh, sub, err := newTestSyncService(true)
	defer sub.Unsubscribe()
	if err != nil {
		t.Fatal(err)
	}

	timestamp := uint64(24)
	target := common.HexToAddress("0x04668ec2f57cc15c381b461b9fedab5d451c8f7f")
	l1TxOrigin := common.HexToAddress("0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8")
	gasLimit := uint64(66)
	data := []byte{0x02, 0x92}
	l1BlockNumber := big.NewInt(100)
	queueIndex := uint64(0)
	index := uint64(0)
	tx := types.NewTransaction(0, target, big.NewInt(0), gasLimit, big.NewInt(0), data, &l1TxOrigin, l1BlockNumber, types.QueueOriginL1ToL2, types.SighashEIP155)
	meta := types.TransactionMeta{
		L1BlockNumber:     l1BlockNumber,
		L1Timestamp:       timestamp,
		L1MessageSender:   &l1TxOrigin,
		SignatureHashType: types.SighashEIP155,
		QueueOrigin:       big.NewInt(int64(types.QueueOriginL1ToL2)),
		Index:             &index,
		QueueIndex:        &queueIndex,
	}
	tx.SetTransactionMeta(&meta)

	setupMockClient(service, map[string]interface{}{
		"GetTransaction": []*types.Transaction{
			tx,
		},
	})

	go service.Loop()

	event := <-txCh
	if len(event.Txs) != 1 {
		t.Fatal("Unexpected number of transactions")
	}
	confirmed := event.Txs[0]

	if !reflect.DeepEqual(tx, confirmed) {
		t.Fatal("different txs")
	}
}

func TestInitializeL1ContextPostGenesis(t *testing.T) {
	service, _, _, err := newTestSyncService(true)
	if err != nil {
		t.Fatal(err)
	}

	timestamp := uint64(24)
	target := common.HexToAddress("0x04668ec2f57cc15c381b461b9fedab5d451c8f7f")
	l1TxOrigin := common.HexToAddress("0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8")
	gasLimit := uint64(66)
	data := []byte{0x02, 0x92}
	l1BlockNumber := big.NewInt(100)
	queueIndex := uint64(100)
	index := uint64(120)
	tx := types.NewTransaction(0, target, big.NewInt(0), gasLimit, big.NewInt(0), data, &l1TxOrigin, l1BlockNumber, types.QueueOriginL1ToL2, types.SighashEIP155)
	meta := types.TransactionMeta{
		L1BlockNumber:     l1BlockNumber,
		L1Timestamp:       timestamp,
		L1MessageSender:   &l1TxOrigin,
		SignatureHashType: types.SighashEIP155,
		QueueOrigin:       big.NewInt(int64(types.QueueOriginL1ToL2)),
		Index:             &index,
		QueueIndex:        &queueIndex,
	}
	tx.SetTransactionMeta(&meta)

	setupMockClient(service, map[string]interface{}{
		"GetEnqueue": []*types.Transaction{
			tx,
		},
	})

	header := types.Header{
		Number: big.NewInt(0),
		Time:   0,
	}

	number := uint64(1)
	tx.SetL1Timestamp(timestamp)
	tx.SetL1BlockNumber(number)
	block := types.NewBlock(&header, []*types.Transaction{tx}, []*types.Header{}, []*types.Receipt{})
	service.bc.SetCurrentBlock(block)

	err = service.initializeLatestL1(big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	latestL1Timestamp := service.GetLatestL1Timestamp()
	latestL1BlockNumber := service.GetLatestL1BlockNumber()
	if number != latestL1BlockNumber {
		t.Fatalf("number does not match, got %d, expected %d", latestL1BlockNumber, number)
	}
	if latestL1Timestamp != timestamp {
		t.Fatal("timestamp does not match")
	}
}

func newTestSyncService(isVerifier bool) (*SyncService, chan core.NewTxsEvent, event.Subscription, error) {
	chainCfg := params.AllEthashProtocolChanges
	chainID := big.NewInt(420)
	chainCfg.ChainID = chainID

	engine := ethash.NewFaker()
	db := rawdb.NewMemoryDatabase()
	_ = new(core.Genesis).MustCommit(db)
	chain, err := core.NewBlockChain(db, nil, chainCfg, engine, vm.Config{}, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Cannot initialize blockchain: %w", err)
	}
	chaincfg := params.ChainConfig{ChainID: chainID}

	txPool := core.NewTxPool(core.TxPoolConfig{PriceLimit: 0}, &chaincfg, chain)
	cfg := Config{
		CanonicalTransactionChainDeployHeight: big.NewInt(0),
		IsVerifier:                            isVerifier,
		// Set as an empty string as this is a dummy value anyways.
		// The client needs to be mocked with a mockClient
		RollupClientHttp: "",
	}

	service, err := NewSyncService(context.Background(), cfg, txPool, chain, db)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Cannot initialize syncservice: %w", err)
	}

	txCh := make(chan core.NewTxsEvent, 1)
	sub := service.SubscribeNewTxsEvent(txCh)

	return service, txCh, sub, nil
}

type mockClient struct {
	getEnqueueCallCount     int
	getEnqueue              []*types.Transaction
	getTransactionCallCount int
	getTransaction          []*types.Transaction
	getEthContextCallCount  int
	getEthContext           []*EthContext
}

func setupMockClient(service *SyncService, responses map[string]interface{}) {
	client := newMockClient(responses)
	service.client = client
}

func newMockClient(responses map[string]interface{}) *mockClient {
	getEnqueueResponses := []*types.Transaction{}
	getTransactionResponses := []*types.Transaction{}
	getEthContextResponses := []*EthContext{}

	enqueue, ok := responses["GetEnqueue"]
	if ok {
		getEnqueueResponses = enqueue.([]*types.Transaction)
	}
	getTx, ok := responses["GetTransaction"]
	if ok {
		getTransactionResponses = getTx.([]*types.Transaction)
	}
	getCtx, ok := responses["GetEthContext"]
	if ok {
		getEthContextResponses = getCtx.([]*EthContext)
	}
	return &mockClient{
		getEnqueue:     getEnqueueResponses,
		getTransaction: getTransactionResponses,
		getEthContext:  getEthContextResponses,
	}
}

func (m *mockClient) GetEnqueue(index uint64) (*types.Transaction, error) {
	if m.getEnqueueCallCount < len(m.getEnqueue) {
		tx := m.getEnqueue[m.getEnqueueCallCount]
		m.getEnqueueCallCount++
		return tx, nil
	}
	return nil, errors.New("")
}

func (m *mockClient) GetLatestEnqueue() (*types.Transaction, error) {
	if len(m.getEnqueue) == 0 {
		return &types.Transaction{}, errors.New("")
	}
	return m.getEnqueue[len(m.getEnqueue)-1], nil
}

func (m *mockClient) GetTransaction(index uint64) (*types.Transaction, error) {
	if m.getTransactionCallCount < len(m.getTransaction) {
		tx := m.getTransaction[m.getTransactionCallCount]
		m.getTransactionCallCount++
		return tx, nil
	}
	return nil, errors.New("")
}

func (m *mockClient) GetLatestTransaction() (*types.Transaction, error) {
	return m.getTransaction[len(m.getTransaction)-1], nil
}

func (m *mockClient) GetEthContext(index uint64) (*EthContext, error) {
	if m.getEthContextCallCount < len(m.getEthContext) {
		ctx := m.getEthContext[m.getEthContextCallCount]
		m.getEthContextCallCount++
		return ctx, nil
	}
	return nil, errors.New("")
}

func (m *mockClient) GetLatestEthContext() (*EthContext, error) {
	return nil, nil
}
