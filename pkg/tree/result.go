package tree

import "container/heap"

type result struct {
	Filename string
	Score    float64
}

type results []result

func (r results) Len() int           { return len(r) }
func (r results) Less(i, j int) bool { return r[i].Score < r[j].Score }
func (r results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func (r *results) Push(x interface{}) {
	*r = append(*r, x.(result))
}

func (r *results) Pop() interface{} {
	old := *r
	n := len(old)
	x := old[n-1]
	*r = old[0 : n-1]
	return x
}

func (r results) Top(limit uint) results {
	top := make(results, 0, limit)
	for heap.Init(&r); len(r) > 0 && limit > 0; limit-- {
		top = append(top, heap.Pop(&r).(result))
	}
	return top
}
