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

package postgres

import "github.com/jmoiron/sqlx"

// TestDB connect to the testing database
// it assumes the database has the IPFS public.blocks table present
// DO NOT use a production db for the test db, as it will remove all contents of the public.blocks table
func TestDB() (*sqlx.DB, error) {
	connectStr := "postgresql://localhost:5432/vulcanize_testing?sslmode=disable"
	return sqlx.Connect("postgres", connectStr)
}

// ResetTestDB drops all rows in the test db public.blocks table
func ResetTestDB(db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	if _, err := tx.Exec("TRUNCATE public.blocks CASCADE"); err != nil {
		return err
	}
	if _, err := tx.Exec("TRUNCATE eth.key_preimages CASCADE"); err != nil {
		return err
	}
	if _, err := tx.Exec("TRUNCATE eth.ancient_headers CASCADE"); err != nil {
		return err
	}
	if _, err := tx.Exec("TRUNCATE eth.ancient_hashes CASCADE"); err != nil {
		return err
	}
	if _, err := tx.Exec("TRUNCATE eth.ancient_bodies CASCADE"); err != nil {
		return err
	}
	if _, err := tx.Exec("TRUNCATE eth.ancient_receipts CASCADE"); err != nil {
		return err
	}
	_, err = tx.Exec("TRUNCATE eth.ancient_tds CASCADE")
	return err
}
