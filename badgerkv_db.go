// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package badgerkv

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
	"github.com/dgraph-io/badger/v4"
)

type DB interface {
	libkv.DB
	Badger() *badger.DB
}

type ChangeOptions func(opts *badger.Options)

func MinMemoryUsageOptions(opts *badger.Options) {
	opts.MemTableSize = 16 << 20
	opts.NumMemtables = 3
	opts.NumLevelZeroTables = 3
	opts.NumLevelZeroTablesStall = 8
}

func OpenPath(ctx context.Context, path string, fn ...ChangeOptions) (libkv.DB, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil
	for _, f := range fn {
		f(&opts)
	}
	db, err := badger.Open(opts)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "open badger db failed")
	}
	return NewDB(db), nil
}

func OpenMemory(ctx context.Context, fn ...ChangeOptions) (libkv.DB, error) {
	opts := badger.DefaultOptions("").WithInMemory(true)
	opts.Logger = nil
	for _, f := range fn {
		f(&opts)
	}
	db, err := badger.Open(opts)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "open badger db failed")
	}
	return NewDB(db), nil
}

func NewDB(db *badger.DB) DB {
	return &badgerdb{
		db: db,
	}
}

type badgerdb struct {
	db *badger.DB
}

func (b *badgerdb) Sync() error {
	return b.db.Sync()
}

func (b *badgerdb) Badger() *badger.DB {
	return b.db
}

func (b *badgerdb) Close() error {
	return b.db.Close()
}

func (b *badgerdb) Update(fn func(tx libkv.Tx) error) error {
	return b.db.Update(func(tx *badger.Txn) error {
		return fn(NewTx(tx))
	})
}

func (b *badgerdb) View(fn func(tx libkv.Tx) error) error {
	return b.db.View(func(tx *badger.Txn) error {
		return fn(NewTx(tx))
	})
}
