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

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgres driver
)

func NewDB(conf *Config) (*sqlx.DB, error) {
	connectString := DbConnectionString(conf)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return nil, err
	}
	if conf.MaxOpen > 0 {
		db.SetMaxOpenConns(conf.MaxOpen)
	}
	if conf.MaxIdle > 0 {
		db.SetMaxIdleConns(conf.MaxIdle)
	}
	if conf.MaxLifetime > 0 {
		db.SetConnMaxLifetime(conf.MaxLifetime)
	}
	return db, nil
}
