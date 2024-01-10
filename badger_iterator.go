// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package badgerkv

import (
	libkv "github.com/bborbe/kv"
	"github.com/dgraph-io/badger/v4"
)

func NewIterator(
	badgerIterator *badger.Iterator,
	bucketName libkv.BucketName,
) libkv.Iterator {
	return &iterator{
		badgerIterator: badgerIterator,
		bucketName:     bucketName,
	}

}

type iterator struct {
	badgerIterator *badger.Iterator
	bucketName     libkv.BucketName
}

func (i iterator) Close() {
	i.badgerIterator.Close()
}

func (i iterator) Item() libkv.Item {
	return NewItem(i.bucketName, i.badgerIterator.Item())
}

func (i iterator) Next() {
	i.badgerIterator.Next()
}

func (i iterator) Valid() bool {
	return i.badgerIterator.ValidForPrefix(i.bucketName)
}

func (i iterator) Rewind() {
	i.badgerIterator.Rewind()
}

func (i iterator) Seek(key []byte) {
	i.badgerIterator.Seek(BucketAddKey(i.bucketName, key))
}
