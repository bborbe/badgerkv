// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package badgerkv_test

import (
	"github.com/bborbe/badgerkv"
	libkv "github.com/bborbe/kv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Key", func() {
	var bytes []byte
	Context("BucketToPrefix", func() {
		JustBeforeEach(func() {
			bytes = badgerkv.BucketToPrefix(libkv.NewBucketName("mybucket"))
		})
		It("returns correct prefix", func() {
			Expect(bytes).To(Equal([]byte("mybucket_")))
		})
	})
	Context("BucketAddKey", func() {
		JustBeforeEach(func() {
			bytes = badgerkv.BucketAddKey(libkv.NewBucketName("mybucket"), []byte("1337"))
		})
		It("returns correct key", func() {
			Expect(bytes).To(Equal([]byte("mybucket_1337")))
		})
	})
	Context("BucketRemoveKey", func() {
		JustBeforeEach(func() {
			bytes = badgerkv.BucketRemoveKey(libkv.NewBucketName("mybucket"), []byte("mybucket_1337"))
		})
		It("returns correct key", func() {
			Expect(bytes).To(Equal([]byte("1337")))
		})
	})
})
