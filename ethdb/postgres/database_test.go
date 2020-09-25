// Copyright 2020 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package postgres_test

import (
	"math/big"

	"github.com/jmoiron/sqlx"

	"github.com/ethereum/go-ethereum/core/rawdb"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/postgres"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	database           ethdb.Database
	db                 *sqlx.DB
	err                error
	testHeader         = types.Header{Number: big.NewInt(1337)}
	testValue, _       = rlp.EncodeToBytes(testHeader)
	testKeccakEthKey   = testHeader.Hash().Bytes()
	testPrefixedEthKey = append(append(testPrefix, postgres.PrefixDelineation...), testKeccakEthKey...)
	testSuffixedEthKey = append(testPrefixedEthKey, []byte("suffix")...)
	testHeaderEthKey   = append(append(append(append(rawdb.HeaderPrefix, postgres.PrefixDelineation...),
		[]byte("number")...), rawdb.NumberDelineation...), testKeccakEthKey...)
	testPreimageEthKey = append(append(rawdb.PreimagePrefix, postgres.PrefixDelineation...), testKeccakEthKey...)
)

var _ = Describe("Database", func() {
	BeforeEach(func() {
		db, err = postgres.TestDB()
		Expect(err).ToNot(HaveOccurred())
		database = postgres.NewDatabase(db)
	})
	AfterEach(func() {
		err = postgres.ResetTestDB(db)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Has - Keccak keys", func() {
		It("returns false if a key-pair doesn't exist in the db", func() {
			has, err := database.Has(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(has).ToNot(BeTrue())
		})
		It("returns true if a key-pair exists in the db", func() {
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data) VALUES ($1, $2)", testKeccakEthKey, testValue)
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
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data, prefix) VALUES ($1, $2, $3)", testPrefixedEthKey, testValue, testPrefix)
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
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data) VALUES ($1, $2)", testSuffixedEthKey, testValue)
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
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data, prefix) VALUES ($1, $2, $3)", testHeaderEthKey, testValue, rawdb.HeaderPrefix)
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
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data, prefix) VALUES ($1, $2, $3)", testPreimageEthKey, testValue, rawdb.PreimagePrefix)
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
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data) VALUES ($1, $2)", testKeccakEthKey, testValue)
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
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data, prefix) VALUES ($1, $2, $3)", testPrefixedEthKey, testValue, testPrefix)
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
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data) VALUES ($1, $2)", testSuffixedEthKey, testValue)
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
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data, prefix) VALUES ($1, $2, $3)", testHeaderEthKey, testValue, rawdb.HeaderPrefix)
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
			_, err = db.Exec("INSERT into eth.kvstore (eth_key, eth_data, prefix) VALUES ($1, $2, $3)", testPreimageEthKey, testValue, rawdb.PreimagePrefix)
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
