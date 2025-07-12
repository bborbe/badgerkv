// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package badgerkv_test

import (
	"context"

	libkv "github.com/bborbe/kv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libbadgerkv "github.com/bborbe/badgerkv"
)

var _ = Describe("Bucket", func() {
	var ctx context.Context
	var db libbadgerkv.DB
	var bucketName libkv.BucketName
	var err error

	BeforeEach(func() {
		ctx = context.Background()
		db, err = libbadgerkv.OpenMemory(ctx)
		Expect(err).To(BeNil())
		bucketName = libkv.NewBucketName("testbucket")
	})

	AfterEach(func() {
		_ = db.Close()
	})

	Context("Get", func() {
		BeforeEach(func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())
				return bucket.Put(ctx, []byte("key1"), []byte("value1"))
			})
			Expect(err).To(BeNil())
		})

		It("returns existing item", func() {
			err = db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.Bucket(ctx, bucketName)
				Expect(err).To(BeNil())

				item, err := bucket.Get(ctx, []byte("key1"))
				Expect(err).To(BeNil())
				Expect(item.Exists()).To(BeTrue())
				Expect(item.Key()).To(Equal([]byte("key1")))

				var value []byte
				err = item.Value(func(val []byte) error {
					value = make([]byte, len(val))
					copy(value, val)
					return nil
				})
				Expect(err).To(BeNil())
				Expect(value).To(Equal([]byte("value1")))
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("returns empty item for non-existing key", func() {
			err = db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.Bucket(ctx, bucketName)
				Expect(err).To(BeNil())

				item, err := bucket.Get(ctx, []byte("nonexistent"))
				Expect(err).To(BeNil())
				Expect(item.Key()).To(Equal([]byte("nonexistent")))

				var value []byte
				err = item.Value(func(val []byte) error {
					value = val
					return nil
				})
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
				return nil
			})
			Expect(err).To(BeNil())
		})
	})

	Context("Put and Delete", func() {
		It("stores and deletes values", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				err = bucket.Put(ctx, []byte("testkey"), []byte("testvalue"))
				Expect(err).To(BeNil())

				item, err := bucket.Get(ctx, []byte("testkey"))
				Expect(err).To(BeNil())

				var value []byte
				err = item.Value(func(val []byte) error {
					value = make([]byte, len(val))
					copy(value, val)
					return nil
				})
				Expect(err).To(BeNil())
				Expect(value).To(Equal([]byte("testvalue")))

				err = bucket.Delete(ctx, []byte("testkey"))
				Expect(err).To(BeNil())

				item, err = bucket.Get(ctx, []byte("testkey"))
				Expect(err).To(BeNil())

				err = item.Value(func(val []byte) error {
					value = val
					return nil
				})
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())

				return nil
			})
			Expect(err).To(BeNil())
		})
	})

	Context("Iterators", func() {
		BeforeEach(func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				err = bucket.Put(ctx, []byte("key1"), []byte("value1"))
				Expect(err).To(BeNil())
				err = bucket.Put(ctx, []byte("key2"), []byte("value2"))
				Expect(err).To(BeNil())
				err = bucket.Put(ctx, []byte("key3"), []byte("value3"))
				Expect(err).To(BeNil())

				return nil
			})
			Expect(err).To(BeNil())
		})

		It("creates forward iterator", func() {
			err = db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.Bucket(ctx, bucketName)
				Expect(err).To(BeNil())

				iterator := bucket.Iterator()
				Expect(iterator).NotTo(BeNil())
				defer iterator.Close()

				keys := []string{}
				iterator.Rewind()
				for iterator.Valid() {
					keys = append(keys, string(iterator.Item().Key()))
					iterator.Next()
				}

				Expect(keys).To(Equal([]string{"key1", "key2", "key3"}))
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("creates reverse iterator", func() {
			err = db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.Bucket(ctx, bucketName)
				Expect(err).To(BeNil())

				iterator := bucket.IteratorReverse()
				Expect(iterator).NotTo(BeNil())
				defer iterator.Close()

				keys := []string{}
				iterator.Rewind()
				for iterator.Valid() {
					keys = append(keys, string(iterator.Item().Key()))
					iterator.Next()
				}

				Expect(keys).To(Equal([]string{"key3", "key2", "key1"}))
				return nil
			})
			Expect(err).To(BeNil())
		})
	})
})
