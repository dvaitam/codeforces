package main

import (
	"bufio"
	"fmt"
	"os"
)

type BIT struct {
	n   int
	bit []int64
}

func newBIT(n int) *BIT {
	return &BIT{n: n, bit: make([]int64, n+2)}
}

func (b *BIT) add(idx int, delta int64) {
	for idx <= b.n {
		b.bit[idx] += delta
		idx += idx & -idx
	}
}

func (b *BIT) sum(idx int) int64 {
	res := int64(0)
	for idx > 0 {
		res += b.bit[idx]
		idx -= idx & -idx
	}
	return res
}

func (b *BIT) rangeAdd(l, r int, delta int64) {
	if l > r {
		return
	}
	b.add(l, delta)
	b.add(r+1, -delta)
}

type event struct {
	l, r   int
	slope  int64
	offset int64
}

type query struct {
	l, r, id int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k, q int
		fmt.Fscan(in, &n, &k, &q)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if k == 0 {
			for i := 0; i < q; i++ {
				fmt.Fscan(in, new(int), new(int))
				fmt.Fprintln(out, 0)
			}
			continue
		}
		N := n - k + 1
		d := make([]int, n+1)
		for i := 1; i <= n; i++ {
			d[i] = a[i] - i
		}
		g := make([]int, N+1)
		if N > 0 {
			size := 2*n + 10
			offset := n + 5
			freqVal := make([]int, size)
			freqCnt := make([]int, k+3)
			best := 0
			inc := func(val int) {
				idx := val + offset
				cur := freqVal[idx]
				if cur > 0 {
					freqCnt[cur]--
				}
				freqVal[idx] = cur + 1
				freqCnt[cur+1]++
				if cur+1 > best {
					best = cur + 1
				}
			}
			dec := func(val int) {
				idx := val + offset
				cur := freqVal[idx]
				freqCnt[cur]--
				freqVal[idx] = cur - 1
				if freqVal[idx] > 0 {
					freqCnt[cur-1]++
				}
				for best > 0 && freqCnt[best] == 0 {
					best--
				}
			}
			for i := 1; i <= k; i++ {
				inc(d[i])
			}
			if best > k {
				best = k
			}
			g[1] = k - best
			for s := 2; s <= N; s++ {
				dec(d[s-1])
				inc(d[s+k-1])
				if best > k {
					best = k
				}
				g[s] = k - best
			}
		}
		prevLess := make([]int, N+1)
		stack := make([]int, 0, N)
		for i := 1; i <= N; i++ {
			for len(stack) > 0 && g[stack[len(stack)-1]] >= g[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				prevLess[i] = 0
			} else {
				prevLess[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}
		nextLess := make([]int, N+1)
		stack = stack[:0]
		for i := N; i >= 1; i-- {
			for len(stack) > 0 && g[stack[len(stack)-1]] > g[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				nextLess[i] = N + 1
			} else {
				nextLess[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}
		eventsAdd := make([][]event, N+2)
		eventsRemove := make([][]event, N+2)
		for i := 1; i <= N; i++ {
			L1 := prevLess[i] + 1
			L2 := i
			rEnd := nextLess[i] - 1
			slope := int64(g[i])
			offStart := -int64(g[i]) * int64(i-1)
			eventsAdd[i] = append(eventsAdd[i], event{L1, L2, slope, offStart})
			if rEnd+1 <= N {
				offEnd := int64(g[i]) * int64(rEnd)
				eventsRemove[rEnd+1] = append(eventsRemove[rEnd+1], event{L1, L2, -slope, offEnd})
			}
		}
		queriesByR := make([][]query, N+1)
		answers := make([]int64, q)
		for idx := 0; idx < q; idx++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			R := r - k + 1
			if R < l {
				answers[idx] = 0
				continue
			}
			if R > N {
				R = N
			}
			queriesByR[R] = append(queriesByR[R], query{l, R, idx})
		}
		slopeBIT := newBIT(N)
		offBIT := newBIT(N)
		for R := 1; R <= N; R++ {
			for _, ev := range eventsAdd[R] {
				slopeBIT.rangeAdd(ev.l, ev.r, ev.slope)
				offBIT.rangeAdd(ev.l, ev.r, ev.offset)
			}
			for _, ev := range eventsRemove[R] {
				slopeBIT.rangeAdd(ev.l, ev.r, ev.slope)
				offBIT.rangeAdd(ev.l, ev.r, ev.offset)
			}
			for _, qu := range queriesByR[R] {
				if qu.l > N {
					answers[qu.id] = 0
					continue
				}
				s := slopeBIT.sum(qu.l)
				c := offBIT.sum(qu.l)
				answers[qu.id] = s*int64(R) + c
			}
		}
		for _, ans := range answers {
			fmt.Fprintln(out, ans)
		}
	}
}
