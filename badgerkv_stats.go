// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package badgerkv

import (
	"context"

	"github.com/bborbe/errors"
	libkv "github.com/bborbe/kv"
)

// Stats returns per-bucket key counts plus the total LSM + value-log size
// reported by Badger. Per-bucket key counting is O(n) because Badger has no
// native bucket abstraction — callers should not poll this hot.
// Per-bucket SizeB is not reported (Badger does not expose per-prefix size).
func (b *badgerdb) Stats(ctx context.Context) (libkv.Stats, error) {
	lsm, vlog := b.db.Size()
	s := libkv.Stats{
		Backend: "badger",
		SizeB:   lsm + vlog,
	}
	err := b.View(ctx, func(ctx context.Context, tx libkv.Tx) error {
		names, err := tx.ListBucketNames(ctx)
		if err != nil {
			return errors.Wrapf(ctx, err, "list bucket names failed")
		}
		for _, name := range names {
			bucket, err := tx.Bucket(ctx, name)
			if err != nil {
				return errors.Wrapf(ctx, err, "get bucket %s failed", name)
			}
			count, err := libkv.Count(ctx, bucket)
			if err != nil {
				return errors.Wrapf(ctx, err, "count bucket %s failed", name)
			}
			s.Buckets = append(s.Buckets, libkv.BucketStats{
				Name:     name,
				KeyCount: count,
			})
		}
		return nil
	})
	if err != nil {
		return libkv.Stats{}, errors.Wrapf(ctx, err, "stats failed")
	}
	return s, nil
}
