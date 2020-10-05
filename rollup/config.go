package rollup

import (
	"crypto/ecdsa"
	"time"
)

type Config struct {
	TxIngestionEnable       bool
	TxIngestionDBHost       string
	TxIngestionDBPort       uint
	TxIngestionDBName       string
	TxIngestionDBUser       string
	TxIngestionDBPassword   string
	TxIngestionPollInterval time.Duration
	TxIngestionSignerKey    *ecdsa.PrivateKey
}

func (c *Config) IsTxIngestionEnabled() bool {
	return c.TxIngestionEnable
}
