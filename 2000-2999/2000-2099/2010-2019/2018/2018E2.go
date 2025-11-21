package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type segment struct {
	l int
	r int
}

type pair struct {
	idx int
	val int
}

type bit struct {
	n    int
	tree []int
}

func newBIT(n int) *bit {
	return &bit{
		n:    n,
		tree: make([]int, n+2),
	}
}

func (b *bit) add(idx, delta int) {
	for idx <= b.n {
		b.tree[idx] += delta
		idx += idx & -idx
	}
}

func (b *bit) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += b.tree[idx]
		idx -= idx & -idx
	}
	return res
}

func (b *bit) lowerBound(target int) int {
	if target <= 0 {
		return 1
	}
	idx := 0
	bitMask := 1
	for (bitMask << 1) <= b.n {
		bitMask <<= 1
	}
	sum := 0
	for bitMask > 0 {
		next := idx + bitMask
		if next <= b.n && sum+b.tree[next] < target {
			sum += b.tree[next]
			idx = next
		}
		bitMask >>= 1
	}
	return idx + 1
}

func countGroups(b *bit, k, n int, total int, lefts []int, logs []int, st [][]int) int {
	if total == 0 || k > n {
		return 0
	}
	pos := k - 1
	cnt := 0
	startIdx := 0
	for {
		pos = startIdx + k - 1
		if pos >= n {
			break
		}
		prefix := 0
		if pos > 0 {
			prefix = b.sum(pos)
		}
		if prefix == total {
			break
		}
		idx := b.lowerBound(prefix + 1)
		if idx > n {
			break
		}
		endIdx := idx - 1
		if endIdx < pos {
			endIdx = pos
		}
		start := endIdx - k + 1
		maxR := queryMax(start, endIdx, logs, st)
		cnt++
		startIdx = upperBound(lefts, maxR)
	}
	return cnt
}

func queryMax(l, r int, logs []int, st [][]int) int {
	if l > r {
		return 0
	}
	length := r - l + 1
	k := logs[length]
	v1 := st[k][l]
	v2 := st[k][r-(1<<k)+1]
	if v1 > v2 {
		return v1
	}
	return v2
}

func upperBound(arr []int, target int) int {
	lo, hi := 0, len(arr)
	for lo < hi {
		mid := (lo + hi) / 2
		if arr[mid] > target {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		segs := make([]segment, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &segs[i].l)
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &segs[i].r)
		}

		sort.Slice(segs, func(i, j int) bool {
			if segs[i].l == segs[j].l {
				return segs[i].r < segs[j].r
			}
			return segs[i].l < segs[j].l
		})

		lenArr := make([]int, n)
		lefts := make([]int, n)
		rights := make([]int, n)
		for i := 0; i < n; i++ {
			lefts[i] = segs[i].l
			rights[i] = segs[i].r
		}
		deque := make([]pair, 0)
		j := 0
		for i := 0; i < n; i++ {
			curR := rights[i]
			for len(deque) > 0 && deque[len(deque)-1].val >= curR {
				deque = deque[:len(deque)-1]
			}
			deque = append(deque, pair{i, curR})
			for {
				for len(deque) > 0 && deque[0].idx < j {
					deque = deque[1:]
				}
				if len(deque) == 0 || deque[0].val >= lefts[i] {
					break
				}
				j++
			}
			lenArr[i] = i - j + 1
		}

		logs := make([]int, n+1)
		for i := 2; i <= n; i++ {
			logs[i] = logs[i/2] + 1
		}
		K := logs[n] + 1
		st := make([][]int, K)
		st[0] = make([]int, n)
		copy(st[0], rights)
		for k := 1; k < K; k++ {
			length := 1 << k
			st[k] = make([]int, n-length+1)
			for i := 0; i+length-1 < n; i++ {
				v1 := st[k-1][i]
				v2 := st[k-1][i+(length>>1)]
				if v1 > v2 {
					st[k][i] = v1
				} else {
					st[k][i] = v2
				}
			}
		}

		buckets := make([][]int, n+1)
		for i := 0; i < n; i++ {
			L := lenArr[i]
			if L > n {
				L = n
			}
			if L > 0 {
				buckets[L] = append(buckets[L], i)
			}
		}

		b := newBIT(n)
		active := 0
		best := 0
		for k := n; k >= 2; k-- {
			for _, idx := range buckets[k] {
				b.add(idx+1, 1)
				active++
			}
			if active == 0 {
				continue
			}
			groups := countGroups(b, k, n, active, lefts, logs, st)
			value := groups * k
			if value > best {
				best = value
			}
		}

		bestK1 := maxNonOverlapping(segs)
		if bestK1 > best {
			best = bestK1
		}
		fmt.Fprintln(out, best)
	}
}

func maxNonOverlapping(segs []segment) int {
	type interval struct {
		l int
		r int
	}
	arr := make([]interval, len(segs))
	for i, s := range segs {
		arr[i] = interval{s.l, s.r}
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].r == arr[j].r {
			return arr[i].l < arr[j].l
		}
		return arr[i].r < arr[j].r
	})
	count := 0
	lastEnd := -1
	for _, seg := range arr {
		if seg.l > lastEnd {
			count++
			lastEnd = seg.r
		}
	}
	return count
}
