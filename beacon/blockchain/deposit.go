// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blockchain

import (
	"context"
	"maps"
	"slices"
	"strconv"
	"time"

	"github.com/berachain/beacon-kit/primitives/math"
)

// defaultRetryInterval processes a deposit event.
const defaultRetryInterval = 20 * time.Second

func (s *Service) depositFetcher(
	ctx context.Context,
	blockNum math.U64,
) {
	if blockNum <= s.eth1FollowDistance {
		s.logger.Info(
			"depositFetcher, nothing to fetch",
			"block num", blockNum,
			"eth1FollowDistance", s.eth1FollowDistance,
		)
		return
	}

	s.fetchAndStoreDeposits(ctx, blockNum-s.eth1FollowDistance)
}

// fetchAndStoreDeposits processes all deposits at a particular EL block height.
// TODO: This could be optimized to process a contiguous range of blocks simultaneously to minimize EL RPC calls.
func (s *Service) fetchAndStoreDeposits(
	ctx context.Context,
	blockNum math.U64,
) {
	blockNumStr := strconv.FormatUint(blockNum.Unwrap(), 10)
	deposits, err := s.depositContract.ReadDeposits(ctx, blockNum, blockNum)
	if err != nil {
		s.logger.Error("Failed to read deposits", "error", err)
		s.metrics.sink.IncrementCounter(
			"beacon_kit.execution.deposit.failed_to_get_block_logs",
			"block_num",
			blockNumStr,
		)
		s.failedBlocksMu.Lock()
		s.failedBlocks[blockNum] = struct{}{}
		s.failedBlocksMu.Unlock()
		return
	}

	if len(deposits) > 0 {
		s.logger.Info(
			"Found deposits on execution layer",
			"block", blockNum, "deposits", len(deposits),
		)
	}

	if err = s.storageBackend.DepositStore().EnqueueDeposits(ctx, deposits); err != nil {
		s.logger.Error("Failed to store deposits", "error", err)
		s.metrics.sink.IncrementCounter(
			"beacon_kit.execution.deposit.failed_to_enqueue_deposits",
			"block_num",
			blockNumStr,
		)
		s.failedBlocksMu.Lock()
		s.failedBlocks[blockNum] = struct{}{}
		s.failedBlocksMu.Unlock()
		return
	}
	s.failedBlocksMu.Lock()
	delete(s.failedBlocks, blockNum)
	s.failedBlocksMu.Unlock()
}

func (s *Service) depositCatchupFetcher(ctx context.Context) {
	ticker := time.NewTicker(defaultRetryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.failedBlocksMu.RLock()
			failedBlks := slices.Collect(maps.Keys(s.failedBlocks))
			s.failedBlocksMu.RUnlock()
			if len(failedBlks) == 0 {
				continue
			}
			s.logger.Warn(
				"Failed to get deposits from block(s), retrying...",
				"num_blocks",
				failedBlks,
			)

			// Fetch deposits for blocks that failed to be processed.
			// TODO: This can be optimized to process all the blocks queried at once by utilizing log query ranges
			// for contiguous ranges of blocks
			for _, blockNum := range failedBlks {
				s.fetchAndStoreDeposits(ctx, blockNum)
			}
		}
	}
}
