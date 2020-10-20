package rollup

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

// Replace implementations with interfaces so that they can be mocked here
// need to figure out how to test this

func TestRollupService(t *testing.T) {
	chainCfg := params.AllEthashProtocolChanges
	chainID := big.NewInt(420)
	chainCfg.ChainID = chainID

	engine := ethash.NewFaker()
	db := rawdb.NewMemoryDatabase()
	_ = new(core.Genesis).MustCommit(db)
	chain, err := core.NewBlockChain(db, nil, chainCfg, engine, vm.Config{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	chaincfg := params.ChainConfig{ChainID: chainID}

	txPool := core.NewTxPool(core.TxPoolConfig{}, &chaincfg, chain)

	cfg := Config{
		StateCommitmentChainAddress: common.Address{},
	}

	_, err = NewSyncService(context.Background(), cfg, txPool, chain, db)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO
// get some good test data
func TestCanonicalChainBatchContext(t *testing.T) {
	tests := []struct {
		input  ctcBatchContext
		expect []byte
	}{
		{
			input: ctcBatchContext{
				NumSequencedTransactions:       big.NewInt(0),
				NumSubsequentQueueTransactions: big.NewInt(0),
				Timestamp:                      big.NewInt(0),
				BlockNumber:                    big.NewInt(0),
			},
			expect: hexutil.MustDecode("0x00"),
		},
		{
			input: ctcBatchContext{
				NumSequencedTransactions:       big.NewInt(10),
				NumSubsequentQueueTransactions: big.NewInt(11),
				Timestamp:                      big.NewInt(12),
				BlockNumber:                    big.NewInt(13),
			},
			expect: hexutil.MustDecode("0x00"),
		},
	}

	for _, test := range tests {
		buf := new(bytes.Buffer)
		err := test.input.Encode(buf)
		if err != nil {
			t.Fatal(err)
		}
		reader := bytes.NewReader(buf.Bytes())
		ctcCtx := ctcBatchContext{}
		err = ctcCtx.Decode(reader)
		if err != nil {
			t.Fatal(err)
		}

		if !isCtcBatchContextEqual(&test.input, &ctcCtx) {
			t.Fatal(err)
		}
	}
}

func isCtcBatchContextEqual(one *ctcBatchContext, two *ctcBatchContext) bool {
	if one == nil && two == nil {
		return true
	}
	if !bytes.Equal(one.NumSequencedTransactions.Bytes(), two.NumSequencedTransactions.Bytes()) {
		return false
	}
	if !bytes.Equal(one.NumSubsequentQueueTransactions.Bytes(), two.NumSubsequentQueueTransactions.Bytes()) {
		return false
	}
	if !bytes.Equal(one.Timestamp.Bytes(), two.Timestamp.Bytes()) {
		return false
	}
	if !bytes.Equal(one.BlockNumber.Bytes(), two.BlockNumber.Bytes()) {
		return false
	}
	return true
}

func TestSequencerBatchCalldata(t *testing.T) {
	tests := []struct {
		input  appendSequencerBatchCallData
		expect []byte
	}{
		{
			input: appendSequencerBatchCallData{
				ChainElements: []chainElement{
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(0),
						BlockNumber: big.NewInt(0),
						TxData:      hexutil.MustDecode("0x1234"),
					},
				},
				Contexts: []ctcBatchContext{
					{
						NumSequencedTransactions:       big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(0),
						Timestamp:                      big.NewInt(0),
						BlockNumber:                    big.NewInt(0),
					},
				},
				ShouldStartAtBatch:    big.NewInt(1234),
				TotalElementsToAppend: big.NewInt(1),
			},
			expect: hexutil.MustDecode("0x00000004d2000001000001000001000000000000000000000000000000021234"),
		},
		{
			input: appendSequencerBatchCallData{
				ChainElements: []chainElement{
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602820447),
						BlockNumber: big.NewInt(0),
						TxData:      hexutil.MustDecode("0x12"),
					},
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602820447),
						BlockNumber: big.NewInt(0),
						TxData:      hexutil.MustDecode("0x1234"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
				},
				Contexts: []ctcBatchContext{
					{
						NumSequencedTransactions:       big.NewInt(2),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602820447),
						BlockNumber:                    big.NewInt(0),
					},
				},
				ShouldStartAtBatch:    big.NewInt(0),
				TotalElementsToAppend: big.NewInt(3),
			},
			expect: hexutil.MustDecode("0x0000000000000003000001000002000001005f89195f0000000000000001120000021234"),
		},
		/*
			{ transactions: [ '0x12', '0x1234', '0x123434', '0x12343434' ],
			  contexts:
			   [ { numSequencedTransactions: 1,
			       numSubsequentQueueTransactions: 1,
			       timestamp: 1602821663,
			       blockNumber: 12 },
			     { numSequencedTransactions: 1,
			       numSubsequentQueueTransactions: 1,
			       timestamp: 1602821663,
			       blockNumber: 12 },
			     { numSequencedTransactions: 1,
			       numSubsequentQueueTransactions: 1,
			       timestamp: 1602821663,
			       blockNumber: 12 },
			     { numSequencedTransactions: 1,
			       numSubsequentQueueTransactions: 1,
			       timestamp: 1602821663,
			       blockNumber: 12 } ],

			  shouldStartAtBatch: 0,
			  totalElementsToAppend: 8 }
			0000000000000008000004000001000001005f891e1f000000000c000001000001005f891e1f000000000c000001000001005f891e1f000000000c000001000001005f891e1f000000000c00000112000002123400000312343400000412343434
		*/
		{
			input: appendSequencerBatchCallData{
				ChainElements: []chainElement{
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602821663),
						BlockNumber: big.NewInt(12),
						TxData:      hexutil.MustDecode("0x12"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602821663),
						BlockNumber: big.NewInt(12),
						TxData:      hexutil.MustDecode("0x1234"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602821663),
						BlockNumber: big.NewInt(12),
						TxData:      hexutil.MustDecode("0x123434"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
					{
						IsSequenced: true,
						Timestamp:   big.NewInt(1602821663),
						BlockNumber: big.NewInt(12),
						TxData:      hexutil.MustDecode("0x12343434"),
					},
					{
						IsSequenced: false,
						Timestamp:   nil,
						BlockNumber: nil,
						TxData:      []byte{},
					},
				},
				Contexts: []ctcBatchContext{
					{
						NumSequencedTransactions:       big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602821663),
						BlockNumber:                    big.NewInt(12),
					},
					{
						NumSequencedTransactions:       big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602821663),
						BlockNumber:                    big.NewInt(12),
					},
					{NumSequencedTransactions: big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602821663),
						BlockNumber:                    big.NewInt(12),
					},
					{
						NumSequencedTransactions:       big.NewInt(1),
						NumSubsequentQueueTransactions: big.NewInt(1),
						Timestamp:                      big.NewInt(1602821663),
						BlockNumber:                    big.NewInt(12),
					},
				},
				ShouldStartAtBatch:    big.NewInt(0),
				TotalElementsToAppend: big.NewInt(8),
			},
			expect: hexutil.MustDecode("0x0000000000000008000004000001000001005f891e1f000000000c000001000001005f891e1f000000000c000001000001005f891e1f000000000c000001000001005f891e1f000000000c00000112000002123400000312343400000412343434"),
		},
	}

	for _, test := range tests {
		buf := new(bytes.Buffer)
		err := test.input.Encode(buf)
		if err != nil {
			t.Fatalf("Cannot encode appendSequencerBatchCallData: %s\n", err.Error())
		}

		if !bytes.Equal(test.expect, buf.Bytes()) {
			t.Fatalf("Serialization mismatch: expect\nexpect:\n%x\ngot:\n%x", test.expect, buf.Bytes())
		}

		reader := bytes.NewReader(buf.Bytes())
		cd := appendSequencerBatchCallData{}
		err = cd.Decode(reader)
		if err != nil {
			t.Fatalf("Error decoding: %s", err.Error())
		}

		if !isSequencerBatchCalldataEqual(&test.input, &cd) {
			t.Fatalf("Deserialization result mismatch:\nexpect:\n%#v\ngot:\n%#v", test.input, cd)
		}
	}
}

func isSequencerBatchCalldataEqual(one, two *appendSequencerBatchCallData) bool {
	if one == nil && two == nil {
		return true
	}
	if xornil(one.ShouldStartAtBatch, two.ShouldStartAtBatch) {
		return false
	}
	if one.ShouldStartAtBatch != nil && two.ShouldStartAtBatch != nil {
		if !bytes.Equal(one.ShouldStartAtBatch.Bytes(), two.ShouldStartAtBatch.Bytes()) {
			return false
		}
	}
	if xornil(one.TotalElementsToAppend, two.TotalElementsToAppend) {
		return false
	}
	if one.TotalElementsToAppend != nil && two.TotalElementsToAppend != nil {
		if !bytes.Equal(one.TotalElementsToAppend.Bytes(), two.TotalElementsToAppend.Bytes()) {
			return false
		}
	}
	if len(one.Contexts) != len(two.Contexts) {
		return false
	}
	for i, oneCtx := range one.Contexts {
		twoCtx := two.Contexts[i]
		if !isCtcBatchContextEqual(&oneCtx, &twoCtx) {
			return false
		}
	}
	if len(one.ChainElements) != len(two.ChainElements) {
		return false
	}
	for i, oneEl := range one.ChainElements {
		twoEl := two.ChainElements[i]
		if oneEl.IsSequenced != twoEl.IsSequenced {
			return false
		}
		if xornil(oneEl.Timestamp, twoEl.Timestamp) {
			return false
		}
		if oneEl.Timestamp != nil && twoEl.Timestamp != nil {
			if !bytes.Equal(oneEl.Timestamp.Bytes(), twoEl.Timestamp.Bytes()) {
				return false
			}
		}
		if xornil(oneEl.BlockNumber, twoEl.BlockNumber) {
			return false
		}
		if oneEl.BlockNumber != nil && twoEl.BlockNumber != nil {
			if !bytes.Equal(oneEl.BlockNumber.Bytes(), twoEl.BlockNumber.Bytes()) {
				return false
			}
		}
		if !bytes.Equal(oneEl.TxData, twoEl.TxData) {
			return false
		}
	}
	return true
}

func xornil(one, two interface{}) bool {
	if one == nil && two != nil {
		return true
	}
	if two == nil && one != nil {
		return true
	}

	return false
}
