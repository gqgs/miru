package tree

import (
	"container/heap"
	"sort"
)

var (
	maxHeapSorter = func(r1, r2 result) bool {
		return r1.Score > r2.Score
	}
	increasingSorter = func(r1, r2 result) bool {
		return r1.Score < r2.Score
	}
)

type (
	result struct {
		Filename string
		Score    float64
	}

	results []result

	sorter struct {
		r      results
		sortBy func(r1, r2 result) bool
	}
)

func (s sorter) Len() int           { return len(s.r) }
func (s sorter) Less(i, j int) bool { return s.sortBy(s.r[i], s.r[j]) }
func (s sorter) Swap(i, j int)      { s.r[i], s.r[j] = s.r[j], s.r[i] }

func (s *sorter) Push(x interface{}) {
	s.r = append(s.r, x.(result))
}

func (s *sorter) Pop() interface{} {
	old := s.r
	n := len(old)
	x := old[n-1]
	s.r = old[0 : n-1]
	return x
}

func (r results) Top(limit uint) results {
	if len(r) <= int(limit) {
		s := sorter{
			r:      r,
			sortBy: increasingSorter,
		}
		sort.Sort(s)
		return s.r
	}
	top := sorter{
		r:      r[:limit],
		sortBy: maxHeapSorter,
	}
	heap.Init(&top)
	for _, res := range r[limit:] {
		if res.Score < top.r[0].Score {
			heap.Pop(&top)
			heap.Push(&top, res)
		}
	}
	top.sortBy = increasingSorter
	sort.Sort(top)
	return top.r
}
