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

var _ = Describe("Iterator", func() {
	var ctx context.Context
	var db libbadgerkv.DB
	var bucketName libkv.BucketName
	var err error

	BeforeEach(func() {
		ctx = context.Background()
		db, err = libbadgerkv.OpenMemory(ctx)
		Expect(err).To(BeNil())
		bucketName = libkv.NewBucketName("testbucket")

		err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
			bucket, err := tx.CreateBucket(ctx, bucketName)
			Expect(err).To(BeNil())

			err = bucket.Put(ctx, []byte("apple"), []byte("red"))
			Expect(err).To(BeNil())
			err = bucket.Put(ctx, []byte("banana"), []byte("yellow"))
			Expect(err).To(BeNil())
			err = bucket.Put(ctx, []byte("cherry"), []byte("red"))
			Expect(err).To(BeNil())

			return nil
		})
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		_ = db.Close()
	})

	Context("Forward iterator", func() {
		It("iterates forward through all items", func() {
			err = db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.Bucket(ctx, bucketName)
				Expect(err).To(BeNil())

				iterator := bucket.Iterator()
				defer iterator.Close()

				var keys []string
				var values []string

				iterator.Rewind()
				for iterator.Valid() {
					item := iterator.Item()
					keys = append(keys, string(item.Key()))

					err := item.Value(func(val []byte) error {
						values = append(values, string(val))
						return nil
					})
					Expect(err).To(BeNil())

					iterator.Next()
				}

				Expect(keys).To(Equal([]string{"apple", "banana", "cherry"}))
				Expect(values).To(Equal([]string{"red", "yellow", "red"}))
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("seeks to specific key", func() {
			err = db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.Bucket(ctx, bucketName)
				Expect(err).To(BeNil())

				iterator := bucket.Iterator()
				defer iterator.Close()

				iterator.Seek([]byte("banana"))
				Expect(iterator.Valid()).To(BeTrue())
				Expect(string(iterator.Item().Key())).To(Equal("banana"))

				iterator.Next()
				Expect(iterator.Valid()).To(BeTrue())
				Expect(string(iterator.Item().Key())).To(Equal("cherry"))

				return nil
			})
			Expect(err).To(BeNil())
		})
	})

	Context("Reverse iterator", func() {
		It("iterates backward through all items", func() {
			err = db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.Bucket(ctx, bucketName)
				Expect(err).To(BeNil())

				iterator := bucket.IteratorReverse()
				defer iterator.Close()

				var keys []string
				var values []string

				iterator.Rewind()
				for iterator.Valid() {
					item := iterator.Item()
					keys = append(keys, string(item.Key()))

					err := item.Value(func(val []byte) error {
						values = append(values, string(val))
						return nil
					})
					Expect(err).To(BeNil())

					iterator.Next()
				}

				Expect(keys).To(Equal([]string{"cherry", "banana", "apple"}))
				Expect(values).To(Equal([]string{"red", "yellow", "red"}))
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("seeks to specific key in reverse", func() {
			err = db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.Bucket(ctx, bucketName)
				Expect(err).To(BeNil())

				iterator := bucket.IteratorReverse()
				defer iterator.Close()

				iterator.Seek([]byte("banana"))
				Expect(iterator.Valid()).To(BeTrue())
				Expect(string(iterator.Item().Key())).To(Equal("banana"))

				iterator.Next()
				Expect(iterator.Valid()).To(BeTrue())
				Expect(string(iterator.Item().Key())).To(Equal("apple"))

				return nil
			})
			Expect(err).To(BeNil())
		})
	})
})
