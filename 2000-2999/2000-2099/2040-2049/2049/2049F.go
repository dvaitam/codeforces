package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// We only care about k = 2^t - 1 where 2^t <= n, because a good subarray
// of length <= n cannot have mex larger than n.

type entry struct {
	size int
	root int
	ver  int
}

type maxHeap []entry

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].size > h[j].size }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(entry))
}
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	v := old[n-1]
	*h = old[:n-1]
	return v
}

func processThreshold(k int, finalVal []int, updates [][2]int, globalAns []int) {
	n := len(finalVal)
	q := len(updates)
	need := k + 1

	// working copy of values while we walk updates in reverse
	val := make([]int, n)
	copy(val, finalVal)

	parent := make([]int, n)
	size := make([]int, n)     // number of active positions in component
	distinct := make([]int, n) // number of distinct values in component (only for roots)
	version := make([]int, n)
	full := make([]bool, n) // component has all values 0..k
	active := make([]bool, n)
	mp := make([]map[int]int, n) // only for roots of active components

	find := func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}

	h := maxHeap{}

	pushIfFull := func(r int) {
		if full[r] {
			heap.Push(&h, entry{size: size[r], root: r, ver: version[r]})
		}
	}

	union := func(a, b int) {
		ra, rb := find(a), find(b)
		if ra == rb {
			return
		}
		version[ra]++
		version[rb]++
		if len(mp[ra]) < len(mp[rb]) {
			ra, rb = rb, ra
		}
		parent[rb] = ra
		size[ra] += size[rb]
		for v, c := range mp[rb] {
			if mp[ra] == nil {
				mp[ra] = make(map[int]int)
			}
			prev := mp[ra][v]
			if prev == 0 {
				distinct[ra]++
			}
			mp[ra][v] = prev + c
		}
		mp[rb] = nil
		distinct[rb] = 0
		size[rb] = 0
		full[rb] = false

		full[ra] = distinct[ra] == need
		version[ra]++
		pushIfFull(ra)
	}

	// initialize components from final array
	for i := 0; i < n; i++ {
		parent[i] = i
		version[i] = 1
		if val[i] <= k {
			active[i] = true
			size[i] = 1
			m := make(map[int]int, 1)
			m[val[i]] = 1
			mp[i] = m
			distinct[i] = 1
			full[i] = distinct[i] == need
			pushIfFull(i)
		}
	}

	for i := 0; i+1 < n; i++ {
		if active[i] && active[i+1] {
			union(i, i+1)
		}
	}

	getBest := func() int {
		for len(h) > 0 {
			top := h[0]
			r := find(top.root)
			if r != top.root || version[r] != top.ver || !full[r] {
				heap.Pop(&h)
				continue
			}
			return top.size
		}
		return 0
	}

	best := getBest()
	if best > globalAns[q] {
		globalAns[q] = best
	}

	for idx := q - 1; idx >= 0; idx-- {
		pos := updates[idx][0]
		delta := updates[idx][1]
		oldVal := val[pos]
		newVal := oldVal - delta

		if oldVal <= k {
			if newVal <= k {
				r := find(pos)
				version[r]++
				// remove oldVal
				cntOld := mp[r][oldVal]
				if cntOld == 1 {
					delete(mp[r], oldVal)
					distinct[r]--
				} else {
					mp[r][oldVal] = cntOld - 1
				}
				// add newVal
				cntNew := mp[r][newVal]
				if cntNew == 0 {
					distinct[r]++
				}
				mp[r][newVal] = cntNew + 1

				full[r] = distinct[r] == need
				version[r]++
				pushIfFull(r)
			}
			// oldVal <= k and newVal > k cannot happen while decrementing
		} else {
			if newVal <= k {
				// activation
				active[pos] = true
				parent[pos] = pos
				size[pos] = 1
				mp[pos] = map[int]int{newVal: 1}
				distinct[pos] = 1
				version[pos]++
				full[pos] = distinct[pos] == need
				pushIfFull(pos)

				if pos > 0 && active[pos-1] {
					union(pos, pos-1)
				}
				if pos+1 < n && active[pos+1] {
					union(pos, pos+1)
				}
			}
		}

		val[pos] = newVal
		best = getBest()
		if best > globalAns[idx] {
			globalAns[idx] = best
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		updates := make([][2]int, q)
		finalVal := make([]int, n)
		copy(finalVal, a)
		for i := 0; i < q; i++ {
			var pos, x int
			fmt.Fscan(in, &pos, &x)
			pos--
			updates[i] = [2]int{pos, x}
			finalVal[pos] += x
		}

		globalAns := make([]int, q+1)

		for p := 1; p <= n; p <<= 1 {
			k := p - 1
			processThreshold(k, finalVal, updates, globalAns)
		}

		for i := 1; i <= q; i++ {
			fmt.Fprintln(out, globalAns[i])
		}
	}
}
