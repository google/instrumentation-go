// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package stats defines the stats collection API and its native Go
// implementation.
package stats

import "time"

// aggregatorSlidingTime indicates that the aggregation occurs over a sliding
// window of time: i.e. last n seconds, minutes, hours...
type aggregatorSlidingCount struct {
	// desiredCount is the actual sample size desired to be aggregated. The
	// subBucketCount is the number of sample  to store in each
	// sub-aggregation. The entries is the set of buckets to keep in memory in
	// order to compute an approximation of the collected data without storing
	// every instance.
	desiredCount   uint64
	subBucketCount uint64
	entries        []*subBucketEntry
	idx            int
}

func (a *aggregatorSlidingCount) isAggregator() bool {
	return true
}

func (a *aggregatorSlidingCount) addSample(v interface{}, now time.Time) {
	e := a.entries[a.idx]
	if e.count == a.subBucketCount {
		a.idx = (a.idx + 1) % len(a.entries)
		e = a.entries[a.idx]
		e.av.clear()
	}
	e.av.addSample(v)
}

func (a *aggregatorSlidingCount) retrieveCollected(now time.Time) AggregationValue {
	e := a.entries[a.idx]
	remaining := float64(a.subBucketCount-e.count) / float64(a.subBucketCount)
	oldestIdx := (a.idx + 1) % len(a.entries)

	e = a.entries[oldestIdx]
	ret := e.av.multiplyByFraction(remaining)

	for j := 1; j < len(a.entries); j++ {
		oldestIdx = (oldestIdx + 1) % len(a.entries)
		e = a.entries[oldestIdx]
		ret.addToIt(e.av)
	}
	return ret
}

type subBucketEntry struct {
	count uint64
	av    AggregationValue
}
