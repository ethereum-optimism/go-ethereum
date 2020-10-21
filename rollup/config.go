package rollup

import (
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	// TODO(mark): deprecate these config options
	TxIngestionEnable       bool
	TxIngestionDBHost       string
	TxIngestionDBPort       uint
	TxIngestionDBName       string
	TxIngestionDBUser       string
	TxIngestionDBPassword   string
	TxIngestionPollInterval time.Duration

	// Ensure that the correct layer 1 chain is being connected to
	Eth1ChainID   big.Int
	Eth1NetworkID big.Int
	// The God Key, used to sign L1 to L2 transactions
	TxIngestionSignerKey *ecdsa.PrivateKey
	// HTTP endpoint of Layer 1 Ethereum node
	httpEndpoint string
	// Addresses of Layer 1 contracts
	CanonicalTransactionChainAddress common.Address
	L1ToL2TransactionQueueAddress    common.Address

	// Deployment Height of the canonical transaction chain
	CanonicalTransactionChainDeployHeight *big.Int
}

func (c *Config) IsTxIngestionEnabled() bool {
	return c.TxIngestionEnable
}
