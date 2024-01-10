// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package badgerkv

import (
	"bytes"
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
	"github.com/dgraph-io/badger/v4"
	"github.com/golang/glog"
)

func NewTx(badgerTx *badger.Txn) libkv.Tx {
	return &tx{
		badgerTx: badgerTx,
	}
}

type tx struct {
	badgerTx *badger.Txn
}

func (t *tx) Bucket(ctx context.Context, name libkv.BucketName) (libkv.Bucket, error) {
	return NewBucket(t.badgerTx, name), nil
}

func (t *tx) CreateBucket(ctx context.Context, name libkv.BucketName) (libkv.Bucket, error) {
	return NewBucket(t.badgerTx, name), nil
}

func (t *tx) CreateBucketIfNotExists(ctx context.Context, name libkv.BucketName) (libkv.Bucket, error) {
	return NewBucket(t.badgerTx, name), nil
}

func (t *tx) DeleteBucket(ctx context.Context, name libkv.BucketName) error {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	it := t.badgerTx.NewIterator(opts)
	defer it.Close()
	for it.Seek(name.Bytes()); it.Valid(); it.Next() {
		key := it.Item().Key()
		if bytes.HasPrefix(key, name.Bytes()) == false {
			glog.V(3).Infof("delete all key of bucket %s completed", name)
			break
		}
		if err := t.badgerTx.Delete(key); err != nil {
			return errors.Wrapf(ctx, err, "delete bucket failed")
		}
	}
	return nil
}
