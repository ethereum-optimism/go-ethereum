package rollup

import "time"

type Config struct {
	TxIngestionEnable       bool
	TxIngestionDBHost       string
	TxIngestionDBPort       uint32
	TxIngestionDBName       string
	TxIngestionDBUser       string
	TxIngestionDBPassword   string
	TxIngestionPollInterval time.Duration
}

func (c *Config) IsTxIngestionEnabled() bool {
	return c.TxIngestionEnable
}
