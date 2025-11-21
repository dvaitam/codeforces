package main

import (
	"bufio"
	"fmt"
	"os"
)

type rangeBIT struct {
	n    int
	bit1 []int64
	bit2 []int64
}

func newRangeBIT(n int) *rangeBIT {
	return &rangeBIT{n: n, bit1: make([]int64, n+2), bit2: make([]int64, n+2)}
}

func (b *rangeBIT) add(bit []int64, idx int, delta int64) {
	for idx <= b.n {
		bit[idx] += delta
		idx += idx & -idx
	}
}

func (b *rangeBIT) prefix(bit []int64, idx int) int64 {
	res := int64(0)
	for idx > 0 {
		res += bit[idx]
		idx -= idx & -idx
	}
	return res
}

func (b *rangeBIT) rangeAdd(l, r int, delta int64) {
	if l > r {
		return
	}
	b.add(b.bit1, l, delta)
	b.add(b.bit1, r+1, -delta)
	b.add(b.bit2, l, delta*int64(l-1))
	b.add(b.bit2, r+1, -delta*int64(r))
}

func (b *rangeBIT) prefixSum(idx int) int64 {
	if idx <= 0 {
		return 0
	}
	s1 := b.prefix(b.bit1, idx)
	s2 := b.prefix(b.bit2, idx)
	return s1*int64(idx) - s2
}

func (b *rangeBIT) rangeSum(l, r int) int64 {
	if l > r {
		return 0
	}
	return b.prefixSum(r) - b.prefixSum(l-1)
}

type event struct {
	l, r   int
	slope  int64
	offset int64
}

type query struct {
	l, id int
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

		N := n - k + 1
		g := make([]int, N+1) // 1-indexed
		if N > 0 {
			d := make([]int, n+1)
			for i := 1; i <= n; i++ {
				d[i] = a[i] - i
			}

			size := 2*n + 10
			offset := n + 5
			freqVal := make([]int, size)
			freqCnt := make([]int, k+2)
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
			val := int64(g[i])
			eventsAdd[i] = append(eventsAdd[i], event{L1, L2, val, -val * int64(i-1)})
			if rEnd+1 <= N {
				eventsRemove[rEnd+1] = append(eventsRemove[rEnd+1], event{L1, L2, -val, val * int64(rEnd)})
			}
		}

		queriesByR := make([][]query, N+1)
		answers := make([]int64, q)
		for idx := 0; idx < q; idx++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			if k > n {
				answers[idx] = 0
				continue
			}
			L := l
			R := r - k + 1
			if R < L {
				answers[idx] = 0
				continue
			}
			queriesByR[R] = append(queriesByR[R], query{L, idx})
		}

		slopeBIT := newRangeBIT(N)
		offBIT := newRangeBIT(N)

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
				l := qu.l
				if l > R {
					answers[qu.id] = 0
					continue
				}
				sumSlope := slopeBIT.rangeSum(l, R)
				sumOff := offBIT.rangeSum(l, R)
				ans := int64(R)*sumSlope + sumOff
				answers[qu.id] = ans
			}
		}

		for _, v := range answers {
			fmt.Fprintln(out, v)
		}
	}
}
