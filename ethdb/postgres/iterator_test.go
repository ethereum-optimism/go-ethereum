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
	"database/sql"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	iterator            ethdb.Iterator
	testPrefix          = []byte("testPrefix")
	testEthKey1         = []byte{'\x01'}
	testEthKey2         = []byte{'\x01', '\x01'}
	testEthKey3         = []byte{'\x01', '\x02'}
	testEthKey4         = []byte{'\x01', '\x0e'}
	testEthKey5         = []byte{'\x01', '\x02', '\x01'}
	testEthKey6         = []byte{'\x01', '\x0e', '\x01'}
	prefixedTestEthKey1 = append(append(testPrefix, postgres.PrefixDelineation...), testEthKey1...)
	prefixedTestEthKey2 = append(append(testPrefix, postgres.PrefixDelineation...), testEthKey2...)
	prefixedTestEthKey3 = append(append(testPrefix, postgres.PrefixDelineation...), testEthKey3...)
	prefixedTestEthKey4 = append(append(testPrefix, postgres.PrefixDelineation...), testEthKey4...)
	prefixedTestEthKey5 = append(append(testPrefix, postgres.PrefixDelineation...), testEthKey5...)
	prefixedTestEthKey6 = append(append(testPrefix, postgres.PrefixDelineation...), testEthKey6...)
	mockValue1          = []byte{1}
	mockValue2          = []byte{2}
	mockValue3          = []byte{3}
	mockValue4          = []byte{4}
	mockValue5          = []byte{5}
	mockValue6          = []byte{6}
)

var _ = Describe("Iterator", func() {
	BeforeEach(func() {
		db, err = postgres.TestDB()
		Expect(err).ToNot(HaveOccurred())
		database = postgres.NewDatabase(db)
		// non-prefixed entries
		err = database.Put(testEthKey1, mockValue1)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(testEthKey2, mockValue2)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(testEthKey3, mockValue3)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(testEthKey4, mockValue4)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(testEthKey5, mockValue5)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(testEthKey6, mockValue6)
		Expect(err).ToNot(HaveOccurred())
		// prefixed entries
		err = database.Put(prefixedTestEthKey1, mockValue1)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(prefixedTestEthKey2, mockValue2)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(prefixedTestEthKey3, mockValue3)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(prefixedTestEthKey4, mockValue4)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(prefixedTestEthKey5, mockValue5)
		Expect(err).ToNot(HaveOccurred())
		err = database.Put(prefixedTestEthKey6, mockValue6)
		Expect(err).ToNot(HaveOccurred())
	})
	AfterEach(func() {
		err = postgres.ResetTestDB(db)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("NewIterator", func() {
		It("iterates over the entire key-set (prefixed or not)", func() {
			iterator = database.NewIterator()
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(BeNil())

			more := iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey1))
			Expect(iterator.Value()).To(Equal(mockValue1))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey3))
			Expect(iterator.Value()).To(Equal(mockValue3))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey5))
			Expect(iterator.Value()).To(Equal(mockValue5))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey4))
			Expect(iterator.Value()).To(Equal(mockValue4))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey6))
			Expect(iterator.Value()).To(Equal(mockValue6))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey1))
			Expect(iterator.Value()).To(Equal(mockValue1))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey3))
			Expect(iterator.Value()).To(Equal(mockValue3))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey5))
			Expect(iterator.Value()).To(Equal(mockValue5))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey4))
			Expect(iterator.Value()).To(Equal(mockValue4))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey6))
			Expect(iterator.Value()).To(Equal(mockValue6))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).ToNot(BeTrue())
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(Equal(sql.ErrNoRows))
		})
	})

	Describe("NewIteratorWithPrefix", func() {
		It("iterates over all db entries that have the provided prefix", func() {
			iterator = database.NewIteratorWithPrefix(testPrefix)
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(BeNil())

			more := iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey1))
			Expect(iterator.Value()).To(Equal(mockValue1))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey3))
			Expect(iterator.Value()).To(Equal(mockValue3))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey5))
			Expect(iterator.Value()).To(Equal(mockValue5))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey4))
			Expect(iterator.Value()).To(Equal(mockValue4))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey6))
			Expect(iterator.Value()).To(Equal(mockValue6))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).ToNot(BeTrue())
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(Equal(sql.ErrNoRows))
		})

		It("behaves as no prefix is provided if prefix is nil", func() {
			iterator = database.NewIteratorWithPrefix(nil)
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(BeNil())

			more := iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey1))
			Expect(iterator.Value()).To(Equal(mockValue1))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey3))
			Expect(iterator.Value()).To(Equal(mockValue3))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey5))
			Expect(iterator.Value()).To(Equal(mockValue5))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey4))
			Expect(iterator.Value()).To(Equal(mockValue4))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey6))
			Expect(iterator.Value()).To(Equal(mockValue6))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey1))
			Expect(iterator.Value()).To(Equal(mockValue1))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey3))
			Expect(iterator.Value()).To(Equal(mockValue3))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey5))
			Expect(iterator.Value()).To(Equal(mockValue5))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey4))
			Expect(iterator.Value()).To(Equal(mockValue4))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey6))
			Expect(iterator.Value()).To(Equal(mockValue6))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).ToNot(BeTrue())
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(Equal(sql.ErrNoRows))
		})

		It("considers empty but non-nil []byte a valid prefix, which precludes iteration over any other prefixed keys", func() {
			iterator = database.NewIteratorWithPrefix([]byte{})
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(BeNil())

			more := iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey1))
			Expect(iterator.Value()).To(Equal(mockValue1))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey3))
			Expect(iterator.Value()).To(Equal(mockValue3))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey5))
			Expect(iterator.Value()).To(Equal(mockValue5))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey4))
			Expect(iterator.Value()).To(Equal(mockValue4))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey6))
			Expect(iterator.Value()).To(Equal(mockValue6))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).ToNot(BeTrue())
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(Equal(sql.ErrNoRows))
		})
	})

	Describe("NewIteratorWithStart", func() {
		It("iterates over the entire key-set (prefixed or not) starting with at the provided path", func() {
			iterator = database.NewIteratorWithStart(testEthKey2)
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(Equal(testEthKey2))
			Expect(iterator.Error()).To(BeNil())

			more := iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey3))
			Expect(iterator.Value()).To(Equal(mockValue3))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey5))
			Expect(iterator.Value()).To(Equal(mockValue5))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey4))
			Expect(iterator.Value()).To(Equal(mockValue4))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey6))
			Expect(iterator.Value()).To(Equal(mockValue6))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey1))
			Expect(iterator.Value()).To(Equal(mockValue1))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey3))
			Expect(iterator.Value()).To(Equal(mockValue3))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey5))
			Expect(iterator.Value()).To(Equal(mockValue5))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey4))
			Expect(iterator.Value()).To(Equal(mockValue4))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey6))
			Expect(iterator.Value()).To(Equal(mockValue6))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).ToNot(BeTrue())
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(Equal(sql.ErrNoRows))
		})

		It("iterates over the entire key-set (prefixed or not) starting with at the provided path", func() {
			iterator = database.NewIteratorWithStart(prefixedTestEthKey3)
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey3))
			Expect(iterator.Error()).To(BeNil())

			more := iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey3))
			Expect(iterator.Value()).To(Equal(mockValue3))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey5))
			Expect(iterator.Value()).To(Equal(mockValue5))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey4))
			Expect(iterator.Value()).To(Equal(mockValue4))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(prefixedTestEthKey6))
			Expect(iterator.Value()).To(Equal(mockValue6))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).ToNot(BeTrue())
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(Equal(sql.ErrNoRows))
		})
	})

	Describe("Release", func() {
		It("releases resources associated with the Iterator", func() {
			iterator = database.NewIteratorWithStart(testEthKey2)
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(Equal(testEthKey2))
			Expect(iterator.Error()).To(BeNil())

			more := iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())

			iterator.Release()
			iterator.Release() // check that we don't panic if called multiple times

			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(BeNil())
			Expect(iterator.Error()).To(BeNil())
			Expect(func() { iterator.Next() }).To(Panic()) // check that we panic if we try to use released iterator

			// We can still create a new iterator from the same backing db
			iterator = database.NewIteratorWithStart(testEthKey2)
			Expect(iterator.Value()).To(BeNil())
			Expect(iterator.Key()).To(Equal(testEthKey2))
			Expect(iterator.Error()).To(BeNil())

			more = iterator.Next()
			Expect(more).To(BeTrue())
			Expect(iterator.Key()).To(Equal(testEthKey2))
			Expect(iterator.Value()).To(Equal(mockValue2))
			Expect(iterator.Error()).To(BeNil())
		})
	})
})
