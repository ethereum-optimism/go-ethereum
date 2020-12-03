package rollup

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	// TODO(mark): deprecate these config options
	TxIngestionEnable bool
	// Number of confs before applying a L1 to L2 tx
	Eth1ConfirmationDepth uint64
	// Verifier mode
	IsVerifier bool
	// Enable the sync service
	Eth1SyncServiceEnable bool
	// Ensure that the correct layer 1 chain is being connected to
	Eth1ChainId   uint64
	Eth1NetworkId uint64
	// Gas Limit
	GasLimit uint64
	// The God Key, used to sign L1 to L2 transactions
	TxIngestionSignerKey *ecdsa.PrivateKey
	// HTTP endpoint of Layer 1 Ethereum node
	Eth1HTTPEndpoint string
	// Addresses of Layer 1 contracts
	AddressResolverAddress           common.Address
	CanonicalTransactionChainAddress common.Address
	L1ToL2TransactionQueueAddress    common.Address
	SequencerDecompressionAddress    common.Address
	L1CrossDomainMessengerAddress    common.Address
	AddressManagerOwnerAddress       common.Address
	// Deployment Height of the canonical transaction chain
	CanonicalTransactionChainDeployHeight *big.Int
	// Path to the state dump
	StateDumpPath string
}

func (c *Config) IsTxIngestionEnabled() bool {
	return c.TxIngestionEnable
}
