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

type DB struct {
	*sqlx.DB
	Node   *NodeInfo
	NodeID int64
}

func NewDB(conf *Config) (*DB, error) {
	connectString := DbConnectionString(conf)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return &DB{}, err
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
	pg := &DB{DB: db, Node: conf.NodeInfo}
	if err := pg.CreateNode(conf.NodeInfo); err != nil {
		return &DB{}, err
	}
	return pg, nil
}

func (db *DB) CreateNode(node *NodeInfo) error {
	var nodeID int64
	err := db.QueryRow(
		`INSERT INTO nodes (genesis_block, network_id, node_id, client_name, chain_id)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT (genesis_block, network_id, node_id, chain_id)
                  DO UPDATE
                    SET genesis_block = $1,
                        network_id = $2,
                        node_id = $3,
                        client_name = $4,
						chain_id = $5
                RETURNING id`,
		node.GenesisBlock, node.NetworkID, node.ID, node.ClientName, node.ChainID).Scan(&nodeID)
	if err != nil {
		return err
	}
	db.NodeID = nodeID
	return nil
}
