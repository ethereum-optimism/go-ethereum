package rollup

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// restart reorg height - depth to reorg at startup

type Config struct {
	// TODO(mark): deprecate these config options
	TxIngestionEnable bool
	// Maximum calldata size for a Queue Origin Sequencer Tx
	MaxCallDataSize int
	// Number of confs before applying a L1 to L2 tx
	Eth1ConfirmationDepth uint64
	// Verifier mode
	IsVerifier bool
	// Enable the sync service
	Eth1SyncServiceEnable bool
	// Ensure that the correct layer 1 chain is being connected to
	Eth1ChainId uint64
	// Gas Limit
	GasLimit uint64
	// HTTP endpoint of the data transport layer
	RollupClientHttp              string
	L1CrossDomainMessengerAddress common.Address
	AddressManagerOwnerAddress    common.Address
	// Deployment Height of the canonical transaction chain
	CanonicalTransactionChainDeployHeight *big.Int
	// Path to the state dump
	StateDumpPath string
	// Temporary setting to disable transfers
	DisableTransfers  bool
	InitialReorgDepth uint64
}

func (c *Config) IsTxIngestionEnabled() bool {
	return c.TxIngestionEnable
}
