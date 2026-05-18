// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package badgerkv_test

import (
	"context"

	libkv "github.com/bborbe/kv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/badgerkv"
)

var _ = Describe("BadgerKV Stats", func() {
	var ctx context.Context
	var db libkv.DB

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		db, err = badgerkv.OpenMemory(ctx)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		_ = db.Close()
	})

	Context("Stats (fast)", func() {
		It("reports backend name and Detailed=false", func() {
			stats, err := db.Stats(ctx)
			Expect(err).To(BeNil())
			Expect(stats).NotTo(BeNil())
			Expect(stats.Backend).To(Equal("badger"))
			Expect(stats.Detailed).To(BeFalse())
		})

		It("returns empty buckets list for fresh db", func() {
			stats, err := db.Stats(ctx)
			Expect(err).To(BeNil())
			Expect(stats.Buckets).To(BeEmpty())
		})

		It("returns bucket names but leaves KeyCount at zero", func() {
			bucketName := libkv.NewBucketName("test-bucket")
			err := db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists(ctx, bucketName)
				Expect(err).To(BeNil())
				Expect(bucket.Put(ctx, []byte("k1"), []byte("v1"))).To(Succeed())
				Expect(bucket.Put(ctx, []byte("k2"), []byte("v2"))).To(Succeed())
				return nil
			})
			Expect(err).To(BeNil())

			stats, err := db.Stats(ctx)
			Expect(err).To(BeNil())
			Expect(stats.Buckets).To(HaveLen(1))
			Expect(stats.Buckets[0].Name).To(Equal(bucketName))
			Expect(stats.Buckets[0].KeyCount).To(Equal(int64(0)))
		})
	})

	Context("StatsDetailed", func() {
		It("reports backend name and Detailed=true", func() {
			stats, err := db.StatsDetailed(ctx)
			Expect(err).To(BeNil())
			Expect(stats.Backend).To(Equal("badger"))
			Expect(stats.Detailed).To(BeTrue())
		})

		It("counts keys per bucket", func() {
			bucketName := libkv.NewBucketName("test-bucket")
			err := db.Update(ctx, func(ctx context.Context, tx libkv.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists(ctx, bucketName)
				Expect(err).To(BeNil())
				Expect(bucket.Put(ctx, []byte("k1"), []byte("v1"))).To(Succeed())
				Expect(bucket.Put(ctx, []byte("k2"), []byte("v2"))).To(Succeed())
				Expect(bucket.Put(ctx, []byte("k3"), []byte("v3"))).To(Succeed())
				return nil
			})
			Expect(err).To(BeNil())

			stats, err := db.StatsDetailed(ctx)
			Expect(err).To(BeNil())
			Expect(stats.Buckets).To(HaveLen(1))
			Expect(stats.Buckets[0].Name).To(Equal(bucketName))
			Expect(stats.Buckets[0].KeyCount).To(Equal(int64(3)))
		})
	})
})
