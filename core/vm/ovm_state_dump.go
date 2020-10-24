package vm

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type OVMDumpAccount struct {
	Address common.Address			`json:"address"`
	Code string						`json:"code"`
	CodeHash string					`json:"codeHash"`
	Storage	map[common.Hash]string	`json:"storage"`
	ABI abi.ABI						`json:"abi"`
}

type OVMDump struct {
	Accounts map[string]OVMDumpAccount	`json:"accounts"`
}

var OVMStateDump OVMDump
var OVMStateManager OVMDumpAccount
var OVMExecutionManager OVMDumpAccount
var OVMSequencerMessageDecompressor OVMDumpAccount
var UsingOVM bool

func init() {
	var err error

	err = json.Unmarshal(RawOVMStateDump, &OVMStateDump)
	if err != nil {
		panic(fmt.Errorf("could not decode OVM state dump: %v", err))
	}

	OVMStateManager = OVMStateDump.Accounts["OVM_StateManager"]
	OVMExecutionManager = OVMStateDump.Accounts["OVM_ExecutionManager"]
	OVMSequencerMessageDecompressor = OVMStateDump.Accounts["OVM_SequencerMessageDecompressor"]

	UsingOVM = os.Getenv("USING_OVM") == "true"
}
