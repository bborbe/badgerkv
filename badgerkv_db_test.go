// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package badgerkv_test

import (
	"context"
	"os"

	libbadgerkv "github.com/bborbe/badgerkv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DB", func() {
	var ctx context.Context
	var err error
	var db libbadgerkv.DB
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("OpenPath", func() {
		JustBeforeEach(func() {
			db, err = libbadgerkv.OpenPath(ctx, os.TempDir())
		})
		AfterEach(func() {
			_ = db.Close()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns db", func() {
			Expect(db).NotTo(BeNil())
		})
	})
	Context("OpenMemory", func() {
		JustBeforeEach(func() {
			db, err = libbadgerkv.OpenMemory(ctx)
		})
		AfterEach(func() {
			_ = db.Close()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns db", func() {
			Expect(db).NotTo(BeNil())
		})
	})
})
