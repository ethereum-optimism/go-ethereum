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

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/postgres"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	batch             ethdb.Batch
	testHeader2       = types.Header{Number: big.NewInt(2)}
	testValue2, _     = rlp.EncodeToBytes(testHeader2)
	testKeccakEthKey2 = testHeader2.Hash().Bytes()
)

var _ = Describe("Batch", func() {
	BeforeEach(func() {
		db, err = postgres.TestDB()
		Expect(err).ToNot(HaveOccurred())
		database = postgres.NewDatabase(db)
		batch = database.NewBatch()
	})
	AfterEach(func() {
		err = postgres.ResetTestDB(db)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Put/Write", func() {
		It("adds the key-value pair to the batch", func() {
			_, err = database.Get(testKeccakEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
			_, err = database.Get(testKeccakEthKey2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))

			err = batch.Put(testKeccakEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			err = batch.Put(testKeccakEthKey2, testValue2)
			Expect(err).ToNot(HaveOccurred())
			err = batch.Write()
			Expect(err).ToNot(HaveOccurred())

			val, err := database.Get(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
			val2, err := database.Get(testKeccakEthKey2)
			Expect(err).ToNot(HaveOccurred())
			Expect(val2).To(Equal(testValue2))
		})
	})

	Describe("Delete/Reset/Write", func() {
		It("deletes the key-value pair in the batch", func() {
			err = batch.Put(testKeccakEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			err = batch.Put(testKeccakEthKey2, testValue2)
			Expect(err).ToNot(HaveOccurred())
			err = batch.Write()
			Expect(err).ToNot(HaveOccurred())

			batch.Reset()
			err = batch.Delete(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			err = batch.Delete(testKeccakEthKey2)
			Expect(err).ToNot(HaveOccurred())
			err = batch.Write()
			Expect(err).ToNot(HaveOccurred())

			_, err = database.Get(testKeccakEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
			_, err = database.Get(testKeccakEthKey2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
		})
	})

	Describe("ValueSize/Reset", func() {
		It("returns the size of data in the batch queued for write", func() {
			err = batch.Put(testKeccakEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			err = batch.Put(testKeccakEthKey2, testValue2)
			Expect(err).ToNot(HaveOccurred())
			err = batch.Write()
			Expect(err).ToNot(HaveOccurred())

			size := batch.ValueSize()
			Expect(size).To(Equal(len(testValue) + len(testValue2)))

			batch.Reset()
			size = batch.ValueSize()
			Expect(size).To(Equal(0))
		})
	})

	Describe("Replay", func() {
		It("returns the size of data in the batch queued for write", func() {
			err = batch.Put(testKeccakEthKey, testValue)
			Expect(err).ToNot(HaveOccurred())
			err = batch.Put(testKeccakEthKey2, testValue2)
			Expect(err).ToNot(HaveOccurred())

			_, err = database.Get(testKeccakEthKey)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))
			_, err = database.Get(testKeccakEthKey2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sql: no rows in result set"))

			err = batch.Replay(database)
			Expect(err).ToNot(HaveOccurred())

			val, err := database.Get(testKeccakEthKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(testValue))
			val2, err := database.Get(testKeccakEthKey2)
			Expect(err).ToNot(HaveOccurred())
			Expect(val2).To(Equal(testValue2))
		})
	})
})
