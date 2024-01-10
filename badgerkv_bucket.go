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

func NewBucket(
	badgerTx *badger.Txn,
	bucketName libkv.BucketName,
) libkv.Bucket {
	return &bucket{
		bucketName: bucketName,
		badgerTx:   badgerTx,
	}
}

type bucket struct {
	badgerTx   *badger.Txn
	bucketName libkv.BucketName
}

func (b *bucket) Iterator() libkv.Iterator {
	return b.createIterator(false)
}

func (b *bucket) IteratorReverse() libkv.Iterator {
	return b.createIterator(true)
}

func (b *bucket) createIterator(reverse bool) libkv.Iterator {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	opts.Reverse = reverse
	return NewIterator(
		b.badgerTx.NewIterator(opts),
		b.bucketName,
	)
}

func (b *bucket) Get(ctx context.Context, key []byte) (libkv.Item, error) {
	item, err := b.badgerTx.Get(BucketAddKey(b.bucketName, key))
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return libkv.NewByteItem(key, nil), nil
		}
		return nil, errors.Wrapf(ctx, err, "get failed")
	}
	return NewItem(b.bucketName, item), nil
}

func (b *bucket) Put(ctx context.Context, key []byte, value []byte) error {
	return b.badgerTx.Set(BucketAddKey(b.bucketName, key), value)
}

func (b *bucket) Delete(ctx context.Context, key []byte) error {
	return b.badgerTx.Delete(BucketAddKey(b.bucketName, key))
}
