package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const negInf = -1000000000

// Event represents a kick with transformed coordinates.
type Event struct {
	E    int64 // a + v*t
	Didx int   // compressed index of a - v*t (reversed)
	t    int   // time (for sorting)
}

type BIT struct {
	n        int
	tree     []int
	modified []int
}

func NewBIT(n int) *BIT {
	b := &BIT{n: n, tree: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		b.tree[i] = negInf
	}
	return b
}

func (b *BIT) update(pos, val int) {
	for i := pos; i <= b.n; i += i & -i {
		if b.tree[i] < val {
			b.tree[i] = val
			b.modified = append(b.modified, i)
		}
	}
}

func (b *BIT) query(pos int) int {
	res := negInf
	for i := pos; i > 0; i -= i & -i {
		if b.tree[i] > res {
			res = b.tree[i]
		}
	}
	return res
}

func (b *BIT) clear() {
	for _, i := range b.modified {
		b.tree[i] = negInf
	}
	b.modified = b.modified[:0]
}

// cdq computes DP via divide and conquer over events[l:r].
func cdq(events []Event, dp []int, bit *BIT, l, r int) {
	if r-l <= 1 {
		return
	}
	mid := (l + r) >> 1
	cdq(events, dp, bit, l, mid)

	left := make([]int, mid-l)
	right := make([]int, r-mid)
	for i := l; i < mid; i++ {
		left[i-l] = i
	}
	for i := mid; i < r; i++ {
		right[i-mid] = i
	}
	sort.Slice(left, func(i, j int) bool {
		return events[left[i]].E < events[left[j]].E
	})
	sort.Slice(right, func(i, j int) bool {
		return events[right[i]].E < events[right[j]].E
	})

	p := 0
	for _, idx := range right {
		for p < len(left) && events[left[p]].E <= events[idx].E {
			bit.update(events[left[p]].Didx, dp[left[p]])
			p++
		}
		best := bit.query(events[idx].Didx)
		if best+1 > dp[idx] {
			dp[idx] = best + 1
		}
	}
	bit.clear()

	cdq(events, dp, bit, mid, r)
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var v int64
	fmt.Fscan(in, &n, &v)
	t := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &t[i])
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	events := make([]Event, n)
	Ds := make([]int64, n)
	for i := 0; i < n; i++ {
		E := a[i] + v*int64(t[i])
		D := a[i] - v*int64(t[i])
		events[i] = Event{E: E, t: t[i]}
		Ds[i] = D
	}

	sort.Slice(Ds, func(i, j int) bool { return Ds[i] < Ds[j] })
	uniq := Ds[:0]
	for i, v := range Ds {
		if i == 0 || v != Ds[i-1] {
			uniq = append(uniq, v)
		}
	}
	m := len(uniq)
	mp := make(map[int64]int, m)
	for i, v := range uniq {
		mp[v] = m - i
	}
	for i := 0; i < n; i++ {
		D := a[i] - v*int64(t[i])
		events[i].Didx = mp[D]
	}

	// sort events by time (already increasing but ensure)
	sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })

	dp := make([]int, n)
	for i := 0; i < n; i++ {
		if abs64(a[i]) <= v*int64(t[i]) {
			dp[i] = 1
		} else {
			dp[i] = negInf
		}
	}

	bit := NewBIT(m)
	cdq(events, dp, bit, 0, n)

	ans := 0
	for i := 0; i < n; i++ {
		if dp[i] > ans {
			ans = dp[i]
		}
	}
	fmt.Println(ans)
}
