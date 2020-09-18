// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package postgres_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	pgipfsethdb "github.com/ethereum/go-ethereum/ethdb/postgres"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	database         ethdb.Database
	db               *sqlx.DB
	err              error
	testHeader       = types.Header{Number: big.NewInt(1337)}
	testValue, _     = rlp.EncodeToBytes(testHeader)
	testKeccakEthKey = testHeader.Hash().Bytes()
	testMhKey, _     = pgipfsethdb.MultihashKeyFromKeccak256(testKeccakEthKey)

	testPrefixedEthKey = append(append([]byte("prefix"), pgipfsethdb.KeyDelineation...), testKeccakEthKey...)
	testPrefixedDsKey  = common.Bytes2Hex(testPrefixedEthKey)

	testSuffixedEthKey = append(append(testPrefixedEthKey, pgipfsethdb.KeyDelineation...), []byte("suffix")...)
	testSuffixedDsKey  = common.Bytes2Hex(testSuffixedEthKey)

	testHeaderEthKey = append(append(append(append(pgipfsethdb.HeaderPrefix, pgipfsethdb.KeyDelineation...),
		[]byte("number")...), pgipfsethdb.NumberDelineation...), testKeccakEthKey...)
	testHeaderDsKey = testMhKey

	testPreimageEthKey = append(append(pgipfsethdb.PreimagePrefix, pgipfsethdb.KeyDelineation...), testKeccakEthKey...)
	testPreimageDsKey  = testMhKey
)

var _ = Describe("Database", func() {
	BeforeEach(func() {
		db, err = pgipfsethdb.TestDB()
		Expect(err).ToNot(HaveOccurred())
		database = pgipfsethdb.NewDatabase(db)
	})
	AfterEach(func() {
		err = pgipfsethdb.ResetTestDB(db)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Has - Keccak keys", func() {
		It("returns false if a key-pair doesn't exist in the db", func() {
			has, err := database.Has(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).ToNot(BeTrue())
		})
		It("returns true if a key-pair exists in the db", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testMhKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testKeccakEthKey, testMhKey)
			Expect(err).ToNot(HaveOccurred())
			has, err := database.Has(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).To(BeTrue())
		})
	})

	Describe("Has - Prefixed keys", func() {
		It("returns false if a key-pair doesn't exist in the db", func() {
			has, err := database.Has(testPrefixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).ToNot(BeTrue())
		})
		It("returns true if a key-pair exists in the db", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testPrefixedDsKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testPrefixedEthKey, testPrefixedDsKey)
			Expect(err).ToNot(HaveOccurred())
			has, err := database.Has(testPrefixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).To(BeTrue())
		})
	})

	Describe("Has - Suffixed keys", func() {
		It("returns false if a key-pair doesn't exist in the db", func() {
			has, err := database.Has(testSuffixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).ToNot(BeTrue())
		})
		It("returns true if a key-pair exists in the db", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testSuffixedDsKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testSuffixedEthKey, testSuffixedDsKey)
			Expect(err).ToNot(HaveOccurred())
			has, err := database.Has(testSuffixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).To(BeTrue())
		})
	})

	Describe("Has - Header keys", func() {
		It("returns false if a key-pair doesn't exist in the db", func() {
			has, err := database.Has(testHeaderEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).ToNot(BeTrue())
		})
		It("returns true if a key-pair exists in the db", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testHeaderDsKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testHeaderEthKey, testHeaderDsKey)
			Expect(err).ToNot(HaveOccurred())
			has, err := database.Has(testHeaderEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).To(BeTrue())
		})
	})

	Describe("Has - Preimage keys", func() {
		It("returns false if a key-pair doesn't exist in the db", func() {
			has, err := database.Has(testPreimageEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).ToNot(BeTrue())
		})
		It("returns true if a key-pair exists in the db", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testPreimageDsKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testPreimageEthKey, testPreimageDsKey)
			Expect(err).ToNot(HaveOccurred())
			has, err := database.Has(testPreimageEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).To(BeTrue())
		})
	})

	Describe("Get - Keccak keys", func() {
		It("throws an err if the key-pair doesn't exist in the db", func() {
			_, err = database.Get(testKeccakEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
		It("returns the value associated with the key, if the pair exists", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testMhKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testKeccakEthKey, testMhKey)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Get - Prefixed keys", func() {
		It("throws an err if the key-pair doesn't exist in the db", func() {
			_, err = database.Get(testPrefixedEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
		It("returns the value associated with the key, if the pair exists", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testPrefixedDsKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testPrefixedEthKey, testPrefixedDsKey)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testPrefixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Get - Suffixed keys", func() {
		It("throws an err if the key-pair doesn't exist in the db", func() {
			_, err = database.Get(testSuffixedEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
		It("returns the value associated with the key, if the pair exists", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testSuffixedDsKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testSuffixedEthKey, testSuffixedDsKey)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testSuffixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Get - Header keys", func() {
		It("throws an err if the key-pair doesn't exist in the db", func() {
			_, err = database.Get(testHeaderEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
		It("returns the value associated with the key, if the pair exists", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testHeaderDsKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testHeaderEthKey, testHeaderDsKey)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testHeaderEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Get - Preimage keys", func() {
		It("throws an err if the key-pair doesn't exist in the db", func() {
			_, err = database.Get(testPreimageEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
		It("returns the value associated with the key, if the pair exists", func() {
			_, err = db.Exec("INSERT into public.blocks (key, data) VALUES ($1, $2)", testPreimageDsKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.Exec("INSERT into eth.key_preimages (eth_key, ipfs_key) VALUES ($1, $2)", testPreimageEthKey, testPreimageDsKey)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testPreimageEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Put - Keccak keys", func() {
		It("persists the key-value pair in the database", func() {
			_, err = database.Get(testKeccakEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))

			err = database.Put(testKeccakEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Put - Prefixed keys", func() {
		It("persists the key-value pair in the database", func() {
			_, err = database.Get(testPrefixedEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))

			err = database.Put(testPrefixedEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testPrefixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Put - Suffixed keys", func() {
		It("persists the key-value pair in the database", func() {
			_, err = database.Get(testSuffixedEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))

			err = database.Put(testSuffixedEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testSuffixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Put - Header keys", func() {
		It("persists the key-value pair in the database", func() {
			_, err = database.Get(testHeaderEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))

			err = database.Put(testHeaderEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testHeaderEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Put - Preimage keys", func() {
		It("persists the key-value pair in the database", func() {
			_, err = database.Get(testPreimageEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))

			err = database.Put(testPreimageEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testPreimageEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
		})
	})

	Describe("Delete - Keccak keys", func() {
		It("removes the key-value pair from the database", func() {
			err = database.Put(testKeccakEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))

			err = database.Delete(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			_, err = database.Get(testKeccakEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
	})

	Describe("Delete - Prefixed keys", func() {
		It("removes the key-value pair from the database", func() {
			err = database.Put(testPrefixedEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testPrefixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))

			err = database.Delete(testPrefixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			_, err = database.Get(testPrefixedEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
	})

	Describe("Delete - Suffixed keys", func() {
		It("removes the key-value pair from the database", func() {
			err = database.Put(testSuffixedEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testSuffixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))

			err = database.Delete(testSuffixedEthKey)
			Expect(err).ToNot(HaveOccurred())
			_, err = database.Get(testSuffixedEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
	})

	Describe("Delete - Header keys", func() {
		It("removes the key-value pair from the database", func() {
			err = database.Put(testHeaderEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testHeaderEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))

			err = database.Delete(testHeaderEthKey)
			Expect(err).ToNot(HaveOccurred())
			_, err = database.Get(testHeaderEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
	})

	Describe("Delete - Preimage keys", func() {
		It("removes the key-value pair from the database", func() {
			err = database.Put(testPreimageEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			val, err := database.Get(testPreimageEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))

			err = database.Delete(testPreimageEthKey)
			Expect(err).ToNot(HaveOccurred())
			_, err = database.Get(testPreimageEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
	})
})
