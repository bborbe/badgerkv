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

var _ = Describe("Transaction", func() {
	var ctx context.Context
	var db libbadgerkv.DB
	var err error

	BeforeEach(func() {
		ctx = context.Background()
		db, err = libbadgerkv.OpenMemory(ctx)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		_ = db.Close()
	})

	Context("IsTransactionOpen", func() {
		It("returns false outside transaction", func() {
			Expect(libbadgerkv.IsTransactionOpen(ctx)).To(BeFalse())
		})

		It("returns true inside transaction", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				Expect(libbadgerkv.IsTransactionOpen(ctx)).To(BeTrue())
				return nil
			})
			Expect(err).To(BeNil())
		})
	})

	Context("Nested transactions", func() {
		It("prevents nested transactions", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				return db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
					return nil
				})
			})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("transaction already open"))
		})

		It("prevents nested view in update", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				return db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
					return nil
				})
			})
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("transaction already open"))
		})
	})

	Context("Bucket operations", func() {
		var bucketName libkv.BucketName

		BeforeEach(func() {
			bucketName = libkv.NewBucketName("testbucket")
		})

		It("creates bucket", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())
				Expect(bucket).NotTo(BeNil())
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("fails to create existing bucket", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				_, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				_, err = tx.CreateBucket(ctx, bucketName)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("already exists"))
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("creates bucket if not exists", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket1, err := tx.CreateBucketIfNotExists(ctx, bucketName)
				Expect(err).To(BeNil())
				Expect(bucket1).NotTo(BeNil())

				bucket2, err := tx.CreateBucketIfNotExists(ctx, bucketName)
				Expect(err).To(BeNil())
				Expect(bucket2).NotTo(BeNil())
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("gets existing bucket", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				_, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				bucket, err := tx.Bucket(ctx, bucketName)
				Expect(err).To(BeNil())
				Expect(bucket).NotTo(BeNil())
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("fails to get non-existing bucket", func() {
			err = db.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
				_, err := tx.Bucket(ctx, bucketName)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("not found"))
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("deletes bucket", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				err = bucket.Put(ctx, []byte("key1"), []byte("value1"))
				Expect(err).To(BeNil())

				err = tx.DeleteBucket(ctx, bucketName)
				Expect(err).To(BeNil())

				_, err = tx.Bucket(ctx, bucketName)
				Expect(err).NotTo(BeNil())
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("fails to delete non-existing bucket", func() {
			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				err := tx.DeleteBucket(ctx, bucketName)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("not found"))
				return nil
			})
			Expect(err).To(BeNil())
		})

		It("lists bucket names", func() {
			bucket1Name := libkv.NewBucketName("bucket1")
			bucket2Name := libkv.NewBucketName("bucket2")

			err = db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				_, err := tx.CreateBucket(ctx, bucket1Name)
				Expect(err).To(BeNil())

				_, err = tx.CreateBucket(ctx, bucket2Name)
				Expect(err).To(BeNil())

				bucketNames, err := tx.ListBucketNames(ctx)
				Expect(err).To(BeNil())
				Expect(len(bucketNames)).To(Equal(2))

				names := make([]string, len(bucketNames))
				for i, name := range bucketNames {
					names[i] = name.String()
				}
				Expect(names).To(ContainElements("bucket1", "bucket2"))
				return nil
			})
			Expect(err).To(BeNil())
		})
	})
})
