package rawdb

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

const (
	// cleanerRecheckInterval is the frequency to check the key-value database for
	// chain progression that might permit new blocks to be removed from the consensus datastore
	cleanerRecheckInterval = time.Minute

	// cleanerBatchLimit is the maximum number of blocks to clean in one batch
	// before doing an fsync and deleting it from the key-value store.
	cleanerBatchLimit = 30000
)

type cleaner struct {
	db      ethdb.KeyValueStore
	cleaned uint64 // Number of blocks already frozen
}

func NewDBCleaner(store ethdb.KeyValueStore) *cleaner {
	hash := ReadHeadBlockHash(store)
	if hash == (common.Hash{}) {
		return &cleaner{db: store}
	}
	number := ReadHeaderNumber(store, hash)
	if number == nil || 2*params.ImmutabilityThreshold > *number {
		return &cleaner{db: store}
	}
	return &cleaner{
		db:      store,
		cleaned: *number - 2*params.ImmutabilityThreshold,
	}
}

func (c *cleaner) clean() {
	nfdb := &nofreezedb{KeyValueStore: c.db}

	for {
		// Retrieve the cleaning threshold.
		hash := ReadHeadBlockHash(nfdb)
		if hash == (common.Hash{}) {
			log.Debug("Current full block hash unavailable") // new chain, empty database
			time.Sleep(cleanerRecheckInterval)
			continue
		}
		number := ReadHeaderNumber(nfdb, hash)
		switch {
		case number == nil:
			log.Error("Current full block number unavailable", "hash", hash)
			time.Sleep(cleanerRecheckInterval)
			continue

		case *number < params.ImmutabilityThreshold:
			log.Debug("Current full block not old enough", "number", *number, "hash", hash, "delay", params.ImmutabilityThreshold)
			time.Sleep(cleanerRecheckInterval)
			continue

		case *number-params.ImmutabilityThreshold <= c.cleaned:
			log.Debug("Ancient blocks frozen already", "number", *number, "hash", hash, "frozen", c.cleaned)
			time.Sleep(cleanerRecheckInterval)
			continue
		}
		head := ReadHeader(nfdb, hash, *number)
		if head == nil {
			log.Error("Current full block unavailable", "number", *number, "hash", hash)
			time.Sleep(cleanerRecheckInterval)
			continue
		}
		// Seems we have data ready to be removed, process in usable batches
		limit := *number - params.ImmutabilityThreshold
		if limit-c.cleaned > cleanerBatchLimit {
			limit = c.cleaned + cleanerBatchLimit
		}
		var (
			start    = time.Now()
			first    = c.cleaned
			ancients = make([]common.Hash, 0, limit-first)
		)
		for c.cleaned < limit {
			// Retrieves hashes within the range
			hash := ReadCanonicalHash(nfdb, c.cleaned)
			if hash == (common.Hash{}) {
				log.Error("Canonical hash missing, can't prune", "number", c.cleaned)
				break
			}
			ancients = append(ancients, hash)
			c.cleaned++
		}
		// Prune block data from the active database
		batch := c.db.NewBatch()
		for i := 0; i < len(ancients); i++ {
			// Always keep the genesis block in active database
			if first+uint64(i) != 0 {
				DeleteBlockWithoutNumber(batch, ancients[i], first+uint64(i))
				DeleteCanonicalHash(batch, first+uint64(i))
			}
		}
		if err := batch.Write(); err != nil {
			log.Crit("Failed to prune canonical blocks", "err", err)
		}
		batch.Reset()
		// Wipe out side chain also.
		for number := first; number < c.cleaned; number++ {
			// Always keep the genesis block in active database
			if number != 0 {
				for _, hash := range ReadAllHashes(c.db, number) {
					DeleteBlock(batch, hash, number)
				}
			}
		}
		if err := batch.Write(); err != nil {
			log.Crit("Failed to prune frozen side blocks", "err", err)
		}
		// Log something friendly for the user
		context := []interface{}{
			"blocks", c.cleaned - first, "elapsed", common.PrettyDuration(time.Since(start)), "number", c.cleaned - 1,
		}
		if n := len(ancients); n > 0 {
			context = append(context, []interface{}{"hash", ancients[n-1]}...)
		}
		log.Info("Deep froze chain segment", context...)

		// Avoid database thrashing with tiny writes
		if c.cleaned-first < cleanerBatchLimit {
			time.Sleep(cleanerRecheckInterval)
		}
	}
}
