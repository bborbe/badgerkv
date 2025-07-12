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

var _ = Describe("Item", func() {
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

	Context("Exists", func() {
		It("always returns true for badger items", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				err = bucket.Put(ctx, []byte("testkey"), []byte("testvalue"))
				Expect(err).To(BeNil())

				item, err := bucket.Get(ctx, []byte("testkey"))
				Expect(err).To(BeNil())
				Expect(item.Exists()).To(BeTrue())

				return nil
			})
			Expect(err).To(BeNil())
		})
	})

	Context("Key", func() {
		It("returns the original key without bucket prefix", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				testKey := []byte("mykey")
				err = bucket.Put(ctx, testKey, []byte("value"))
				Expect(err).To(BeNil())

				item, err := bucket.Get(ctx, testKey)
				Expect(err).To(BeNil())
				Expect(item.Key()).To(Equal(testKey))

				return nil
			})
			Expect(err).To(BeNil())
		})
	})

	Context("Value", func() {
		It("provides access to stored value", func() {
			testValue := []byte("stored value")
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				err = bucket.Put(ctx, []byte("key"), testValue)
				Expect(err).To(BeNil())

				item, err := bucket.Get(ctx, []byte("key"))
				Expect(err).To(BeNil())

				var retrievedValue []byte
				err = item.Value(func(val []byte) error {
					retrievedValue = make([]byte, len(val))
					copy(retrievedValue, val)
					return nil
				})
				Expect(err).To(BeNil())
				Expect(retrievedValue).To(Equal(testValue))

				return nil
			})
			Expect(err).To(BeNil())
		})

		It("handles empty values", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				err = bucket.Put(ctx, []byte("emptykey"), []byte{})
				Expect(err).To(BeNil())

				item, err := bucket.Get(ctx, []byte("emptykey"))
				Expect(err).To(BeNil())

				var retrievedValue []byte
				err = item.Value(func(val []byte) error {
					retrievedValue = val
					return nil
				})
				Expect(err).To(BeNil())
				Expect(retrievedValue).To(Equal([]byte{}))

				return nil
			})
			Expect(err).To(BeNil())
		})
	})
})
